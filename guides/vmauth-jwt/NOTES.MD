Get jwt token from keycloak:
```
curl -s -X POST "http://localhost:8080/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "client_id=grafana" \
    -d "client_secret=yyy" \
    -d "grant_type=password" \
    -d "username=admin" \
    -d "password=admin" | jq -r '.access_token'
```

Test vmgateway auth works curl:
```
curl http://localhost:8431 -H "Authorization: Bearer xxx" -i 
```

Test datasource auth curl:
```
curl http://localhost:8431/api/v1/status/buildinfo -H "Authorization: Bearer xxx" -i
```