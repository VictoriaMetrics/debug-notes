package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	export("gh8807.json")
}

type Metric struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

// /api/v1/export
func export(file string) {
	bb, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var m Metric
	if err := json.Unmarshal(bb, &m); err != nil {
		log.Fatal(err)
	}

	var total, duplicates, stepSum int
	for i := range m.Timestamps {
		t := m.Timestamps[i]
		v := m.Values[i]

		var prevT int64
		var prevV float64
		if i > 0 {
			prevT = m.Timestamps[i-1]
			prevV = m.Values[i-1]
		}
		var step int64
		if prevT > 0 {
			step = t - prevT
			stepSum += int(step)
		}

		fmt.Printf("%v %.0f (step: %d, value change: %.0f)\n", toTime(t), v, step, v-prevV)

		if t == prevT {
			duplicates++
			fmt.Printf("\t>> duplicate at %v", time.Unix(int64(t)/1000, 0))
		}
		total++
	}

	log.Println("total samples", total)
	log.Println("duplicates", duplicates)
	log.Println("avg step", stepSum/(total-duplicates)/1000, "sec")

	minT := m.Timestamps[0]
	maxT := m.Timestamps[len(m.Timestamps)-1]
	log.Println("min date", toTime(minT))
	log.Println("min date", toTime(maxT))
}

func toTime(v int64) time.Time {
	return time.Unix(int64(v)/1000, 0)
}
