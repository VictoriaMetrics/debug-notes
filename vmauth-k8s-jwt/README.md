# vmauth + Kubernetes ServiceAccount JWT (OIDC discovery)

Goal: authenticate requests to vmauth with Kubernetes ServiceAccount tokens, using the
kube-apiserver as OIDC issuer (`/.well-known/openid-configuration` + `/openid/v1/jwks`).

Expected: with the plain config, vmauth cannot fetch discovery because
- the apiserver cert is signed by the cluster CA, and vmauth's OIDC HTTP client has no `tls_ca_file` option;
- discovery endpoints are not readable anonymously by default, and vmauth's OIDC HTTP client sends no credentials.

## 1. Cluster

On arm64 Macs make sure the node image is arm64 (an amd64 image under emulation fails kubeadm
with `seccomp is not supported`). kind >= v0.33 is needed for `kindest/node:v1.37.0`.

```bash
docker pull --platform linux/arm64 kindest/node:v1.37.0
DOCKER_DEFAULT_PLATFORM= kind create cluster --name vmauthissue2 --image kindest/node:v1.37.0
alias k='kubectl --context kind-vmauthissue2'

k get --raw /.well-known/openid-configuration | jq   # issuer: https://kubernetes.default.svc.cluster.local
k get --raw /openid/v1/jwks | jq
```

## 2. Deploy vmsingle + vmauth (OIDC issuer = kube-apiserver)

```bash
k apply -f vmauth-oidc.yaml
k rollout status deploy/vmsingle
k rollout status deploy/vmauth
k logs deploy/vmauth | grep -i oidc
```

The manifest sets `SSL_CERT_DIR=/etc/ssl/certs:/var/run/secrets/kubernetes.io/serviceaccount`
(Go on Linux adds every PEM from these dirs to the system roots), so TLS to the apiserver verifies.
Expect:

```
failed to refresh OIDC provider metadata ... unexpected status code 403 when fetching openid config
```

Without `SSL_CERT_DIR` the error is `tls: failed to verify certificate: x509: certificate signed by unknown authority`.
Avoid `SSL_CERT_FILE=<cluster ca.crt>`: it replaces the system bundle, breaking public IdPs/backends.

## 3. Send a request with a ServiceAccount token

```bash
k create serviceaccount tester
TOKEN=$(k create token tester --audience=vmauth --duration=1h)

k port-forward svc/vmauth 8427 &
curl -i -H "Authorization: Bearer $TOKEN" 'http://localhost:8427/api/v1/query?query=up'
```

Expect `401`: no keys were fetched, so vmauth skips the JWT user.

## 4. Workaround (no vmauth changes)

TLS is already handled by `SSL_CERT_DIR` in the manifest; what remains is the 403.

```bash
# allow anonymous read of the discovery endpoints
k apply -f oidc-anon-rbac.yaml
```

Then repeat step 3: the request is proxied to vmsingle (`200`).
A token without `--audience=vmauth` gets `401` (`no user match jwt token`), as does a request without a token.

SA tokens have no `vm_access` claim, so the config also needs `default_vm_access_claim: {}` (tenant `0:0`,
no filters); without it vmauth verifies the token and still returns `401`
(`jwt token for user "jwt" has no vm_access claim and default_vm_access_claim is not configured`).
`-logInvalidAuthTokens` (set in the manifest) is what surfaces these reasons in the vmauth log.
vmauth reloads `-auth.config` only on SIGHUP, so after changing the ConfigMap run `k rollout restart deploy/vmauth`.

## 5. In-cluster client with a projected token (per-cluster audience)

`curl-client.yaml` mounts a projected SA token with `audience: vmauth-prod-eu-1` and calls vmauth every 5s
(re-reading the token, since kubelet rotates it). `vmauth-oidc.yaml` has a matching user
(`aud: vmauth-prod-eu-1`, `sub: system:serviceaccount:default:curl-client`).

```bash
k apply -f curl-client.yaml
k logs -f deploy/curl-client   # ... http_code=200
```

A unique `aud` per cluster is how to tell clusters apart when they share the default issuer and SA name:
the token from another cluster can claim the same `aud`, but fails signature verification against this
user's keys.

## Cleanup

```bash
kind delete cluster --name vmauthissue2
```

## Notes

- Code: `app/vmauth/oidc.go` — `oidcConfig` has only `issuer`; `oidcHTTPClient` is a global client
  with no auth and default TLS roots.
- Possible fix: add `bearer_token_file` / `tls_ca_file` (via `lib/promauth`) to `jwt.oidc`, so vmauth
  can use its own mounted SA token + CA. The default `system:service-account-issuer-discovery`
  binding already allows all ServiceAccounts, so no anonymous RBAC would be needed.
- Docs mismatch: `docs/victoriametrics/vmauth.md:274` says `vm_access` is optional since v1.147, but
  `app/vmauth/main.go:244` (commit 32f8c4e434) requires `vm_access` or `default_vm_access_claim`.
- `jwks_uri` points to the node IP (`https://172.19.0.2:6443/...`), not the issuer host — keep in mind
  when deciding where credentials are sent.
