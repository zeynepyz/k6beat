package k6metricset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/pkg/errors"
)

type AutoGenerated struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			Sample struct {
				Avg   float64 `json:"avg"`
				Max   float64 `json:"max"`
				Med   float64 `json:"med"`
				Min   float64 `json:"min"`
				P90   float64 `json:"p(90)"`
				P95   float64 `json:"p(95)"`
				Count int     `json:"count"`
				Value int     `json:"value"`
				Rate  float64 `json:"rate"`
			} `json:"sample"`
		} `json:"attributes"`
	} `json:"data"`
}

func Mapping() ([]common.MapStr, error) {
	var err error
	// http://localhost:6565/v1/metrics'e http isteği gönderdiğimde yukarıdaki oluşturduğum structla uyumlu çalışmıyordu. ben de bakacağım metriklere
	// teker teker bakma kararı aldım. url adında çekmek istediğim metrikleri içeren bir array oluşturdum ve for döngüsüyle hepsine teker teker bakmasını sağladım.
	url := [9]string{"vus", "vus_max", "http_reqs", "http_req_duration", "http_req_connecting", "http_req_receiving",
		"http_req_sending", "http_req_tls_handshaking", "http_req_waiting"}

	var events []common.MapStr

	for i := range url {
		res, err := http.Get("http://localhost:6565/v1/metrics/" + url[i])

		if err != nil {
			print("Error connecting K6:", err)
			logp.Err("Error connecting K6: ", err)
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			logp.Err("Returned wrong status code: HTTP %s ", res.Status)
			print("Returned wrong status code:", err)
			return nil, fmt.Errorf("HTTP %s", res.Status)
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			logp.Err("Error reading stats: %v", err)
			print("Error reading stats:", err)
			return nil, fmt.Errorf("HTTP%s", res.Status)
		}

		var data AutoGenerated
		json.Unmarshal(bodyBytes, &data)

		if err != nil {
			logp.Err("Error unmarshal: ", err)
			print("Error unmarsall:", err)
			return nil, err
		}

		event := common.MapStr{
			"data": common.MapStr{
				"id": data.Data.ID,
				"attributes": common.MapStr{
					"sample": common.MapStr{},
				},
			},
		}

		if data.Data.Attributes.Sample.Rate != 0 {
			event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["rate"] = data.Data.Attributes.Sample.Rate
			event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["count"] = data.Data.Attributes.Sample.Count
		} else {
			if data.Data.Attributes.Sample.Value != 0 {
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["value"] = data.Data.Attributes.Sample.Value
			} else {
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["avg"] = data.Data.Attributes.Sample.Avg
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["max"] = data.Data.Attributes.Sample.Max
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["med"] = data.Data.Attributes.Sample.Med
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["p(90)"] = data.Data.Attributes.Sample.P90
				event["data"].(common.MapStr)["attributes"].(common.MapStr)["sample"].(common.MapStr)["p(95)"] = data.Data.Attributes.Sample.P95

			}

		}

		events = append(events, event)

	}

	fmt.Println(events)

	return events, err

}

func GetReport(report mb.ReporterV2) error {

	events, err := Mapping()
	if err != nil {
		print("FAILURE", err)
		return errors.Wrap(err, "failure")
	}

	fmt.Println("Report:", report)

	for _, event := range events {
		report.Event(mb.Event{MetricSetFields: event})
	}

	return nil
}
