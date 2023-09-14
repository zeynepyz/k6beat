package k6metricset

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/beats/v7/libbeat/common"
)

type Sample struct {
	Value int     `json:"value,omitempty"`
	Count int     `json:"count,omitempty"`
	Rate  float64 `json:"rate,omitempty"`
	Avg   float64 `json:"avg,omitempty"`
	Max   float64 `json:"max,omitempty"`
	Med   float64 `json:"med,omitempty"`
	Min   float64 `json:"min,omitempty"`
	P90   float64 `json:"p(90),omitempty"`
	P95   float64 `json:"p(95),omitempty"`
}

type Metric struct {
	ID         string `json:"id"`
	Attributes struct {
		Sample Sample `json:"sample"`
	} `json:"attributes"`
}

type Data struct {
	Metrics []Metric `json:"data"`
}

func eventMapping(response []byte) (common.MapStr, error) {

	var data Data
	var err error

	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshall fail: %v", err)
	}

	event := common.MapStr{
		"data": common.MapStr{
			"metrics": common.MapStr{},
		},
	}

	wantedMetricIDs := []string{"vus", "vus_max", "http_reqs", "http_req_duration", "http_req_connecting", "http_req_receiving",
		"http_req_sending", "http_req_tls_handshaking", "http_req_waiting"}

	contains := func(items []string, item string) bool {
		for _, i := range items {
			if i == item {
				return true
			}
		}
		return false
	}

	for _, metric := range data.Metrics {

		if contains(wantedMetricIDs, metric.ID) {
			sample := metric.Attributes.Sample
			metricFields := common.MapStr{}
			if sample.Rate != 0 {
				metricFields["rate"] = sample.Rate
				metricFields["count"] = sample.Count
			} else {
				if sample.Value != 0 {
					metricFields["value"] = sample.Value
				} else {
					metricFields["avg"] = sample.Avg
					metricFields["max"] = sample.Max
					metricFields["med"] = sample.Med
					metricFields["p(90)"] = sample.P90
					metricFields["p(95)"] = sample.P95

				}

			}
			event["data"].(common.MapStr)["metrics"].(common.MapStr)[metric.ID] = metricFields
		}

	}

	return event, nil
}
