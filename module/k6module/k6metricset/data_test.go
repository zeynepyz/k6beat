package k6metricset

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/elastic/beats/v7/libbeat/common"
	mbtest "github.com/elastic/beats/v7/metricbeat/mb/testing"
	"github.com/stretchr/testify/assert"
)

func TestEventMapping(t *testing.T) {
	content, err := ioutil.ReadFile("./_meta/test/data_for_test.json")
	assert.NoError(t, err)

	event, _ := eventMapping(content)
	if err != nil {
		t.Fatal(err)
	}

	expected := common.MapStr{
		"data": common.MapStr{
			"metrics": common.MapStr{
				"vus": common.MapStr{
					"value": 1,
				},
				"http_req_duration": common.MapStr{
					"avg":   386.793313,
					"max":   430.822131,
					"med":   386.793313,
					"p(90)": 422.01636740000004,
					"p(95)": 426.4192492,
				},
				"http_req_tls_handshaking": common.MapStr{
					"avg":   61.844635,
					"max":   62.025681,
					"med":   61.844635,
					"p(90)": 61.9894718,
					"p(95)": 62.0075764,
				},
				"vus_max": common.MapStr{
					"value": 1,
				},
				"http_req_receiving": common.MapStr{
					"avg":   7.3451695,
					"max":   13.807923,
					"med":   7.345169500000001,
					"p(90)": 12.515372300000001,
					"p(95)": 13.16164765,
				},
				"http_req_sending": common.MapStr{
					"avg":   1.2124685,
					"max":   2.309611,
					"med":   1.2124685,
					"p(90)": 2.0901824999999996,
					"p(95)": 2.1998967499999997,
				},
				"http_req_connecting": common.MapStr{
					"avg":   19.264117499999998,
					"max":   19.343144,
					"med":   19.264117499999998,
					"p(90)": 19.3273387,
					"p(95)": 19.33524135,
				},
				"http_reqs": common.MapStr{
					"count": 2,
					"rate":  0.07031716558239108,
				},
				"http_req_waiting": common.MapStr{
					"avg":   378.235675,
					"max":   416.898882,
					"med":   378.235675,
					"p(90)": 409.16624060000004,
					"p(95)": 413.0325613,
				},
			},
		},
	}

	assert.Equal(t, expected, event)

}

func TestEventMapping_InvalidJSON(t *testing.T) {
	// Geçersiz JSON verisi oluşturun
	invalidJSON := []byte(`{"data": { "metrics": [1, 2, 3] } }`)

	// eventMapping'i çağırın ve bir hata bekleyin
	_, err := eventMapping(invalidJSON)

	// Hatanın beklenen bir hata türü olup olmadığını kontrol edin
	if err == nil {
		t.Errorf("Hata bekleniyor, ancak hata alınmadı.")
	} else if err.Error() != "JSON unmarshall fail: json: cannot unmarshal object into Go struct field Data.data of type []k6metricset.Metric" {
		t.Errorf("Beklenen hata metni alınmadı. Beklenen: 'json: cannot unmarshal object into Go struct field Data.data of type []k6metricset.Metric', Alınan: %v", err)
	}

}

func TestFetchEventContent(t *testing.T) {
	absPath, err := filepath.Abs("../_meta/testdata/")
	assert.NoError(t, err)

	response, err := ioutil.ReadFile(absPath + "/k6metrics")
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Write([]byte(response))
	}))
	defer server.Close()

	config := map[string]interface{}{
		"module":      "k6module",
		"metricsets":  []string{"k6metricset"},
		"k6metricset": []string{server.URL},
	}

	f := mbtest.NewReportingMetricSetV2Error(t, config)
	events, errs := mbtest.ReportingFetchV2Error(f)
	if len(errs) > 0 {
		t.Fatalf("Expected 0 error, had %d. %v\n", len(errs), errs)
	}
	assert.NotEmpty(t, events)

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), events[0])
}
