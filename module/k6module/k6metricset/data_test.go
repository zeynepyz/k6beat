package k6metricset

import (
	"net/http"
	"net/http/httptest"
	"time"
)

var testServer *httptest.Server

func initTestServer() {
	// Test sunucusu oluşturma
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// İstenen duruma göre test verileriyle yanıt verme
		if r.URL.Path == "/v1/metrics" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"data": {
					"metrics": [
						{
							"id": "vus",
							"attributes": {
								"sample": {
									"value": 10
								}
							}
						}
					]
				}
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func closeTestServer() {
	testServer.Close()
}

func SetHTTPClientForTesting() {
	initTestServer()
	httpClient := &http.Client{
		Transport: testServer.Client().Transport,
		Timeout:   5 * time.Second,
	}
	SetHTTPClient(httpClient)
}

/*
func TestMapping(t *testing.T) {
	// Test sunucusunu başlat
	SetHTTPClientForTesting()
	defer closeTestServer()

	// Mapping fonksiyonunu çağırma
	event, err := Mapping()

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	t.Logf("Event: %+v", event)

	// Test işlemleri devam eder...
}
*/
