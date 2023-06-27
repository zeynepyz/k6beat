package k6metricset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/pkg/errors"
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

func Mapping() (common.MapStr, error) {
	res, err := http.Get("http://localhost:6565/v1/metrics")
	if err != nil {
		return nil, err
	}

	//defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Returned wrong status code: HTTP " + res.Status)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	event := common.MapStr{
		"data": common.MapStr{
			"metrics": common.MapStr{},
		},
	}

	wantedMetricIDs := []string{"vus", "vus_max", "http_reqs", "http_req_duration", "http_req_connecting", "http_req_receiving",
		"http_req_sending", "http_req_tls_handshaking", "http_req_waiting"}

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
					metricFields["P(95)"] = sample.P95

				}

			}
			event["data"].(common.MapStr)["metrics"].(common.MapStr)[metric.ID] = metricFields
		}

	}

	fmt.Printf("%+v\n", event)

	return event, nil
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func GetReport(report mb.ReporterV2) error {
	event, err := Mapping()
	if err != nil {
		return errors.Wrap(err, "failure")
	}

	report.Event(mb.Event{MetricSetFields: event})
	return nil
}
