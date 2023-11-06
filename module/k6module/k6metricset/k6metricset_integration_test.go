package k6metricset

import (
	"testing"

	"github.com/elastic/beats/v7/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/v7/metricbeat/mb/testing"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	service := compose.EnsureUp(t, "k6module")

	f := mbtest.NewReportingMetricSetV2Error(t, getConfig(service.Host()))
	events, errs := mbtest.ReportingFetchV2Error(f)
	if len(errs) > 0 {
		t.Fatalf("Expected 0 error, had %d. %v\n", len(errs), errs)
	}
	assert.NotEmpty(t, events)

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), events[0])
}

func TestData(t *testing.T) {
	service := compose.EnsureUp(t, "k6module")

	f := mbtest.NewReportingMetricSetV2Error(t, getConfig(service.Host()))
	if err := mbtest.WriteEventsReporterV2Error(f, t, ""); err != nil {
		t.Fatal("write", err)
	}
}

func getConfig(host string) map[string]interface{} {
	return map[string]interface{}{
		"module":     "k6module",
		"metricsets": []string{"k6metricset"},
		"hosts":      []string{host},
	}
}
