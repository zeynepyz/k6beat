package k6metricset

import (
	"encoding/json"
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
	ID         string        `json:"id"`
	Attributes common.MapStr `json:"attributes"`
}

type Data struct {
	Metrics []Metric `json:"data"`
}

func Mapping() (common.MapStr, error) {
	res, err := http.Get("http://localhost:6565/v1/metrics")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

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

	for _, metric := range data.Metrics {
		sample := metric.Attributes["sample"].(Sample)
		metricFields := common.MapStr{}

		if sample.Rate != 0 {
			metricFields["rate"] = sample.Rate
		}

		if sample.Value != 0 {
			metricFields["value"] = sample.Value
		}

		if sample.Avg != 0 {
			metricFields["avg"] = sample.Avg
		}

		if sample.Max != 0 {
			metricFields["max"] = sample.Max
		}

		if sample.Med != 0 {
			metricFields["med"] = sample.Med
		}

		if sample.Min != 0 {
			metricFields["min"] = sample.Min
		}

		if sample.P90 != 0 {
			metricFields["p(90)"] = sample.P90
		}

		if sample.P95 != 0 {
			metricFields["p(95)"] = sample.P95
		}

		event["data"].(common.MapStr)["metrics"].(common.MapStr)[metric.ID] = metricFields
	}

	print(event)

	return event, nil
}

func GetReport(report mb.ReporterV2) error {
	event, err := Mapping()
	if err != nil {
		return errors.Wrap(err, "failure")
	}

	report.Event(mb.Event{MetricSetFields: event})
	print(report)

	return nil
}
