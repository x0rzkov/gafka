package influxdb

import (
	"net/url"
	"testing"
	"time"

	"github.com/funkygao/assert"
	"github.com/funkygao/gafka/ctx"
	"github.com/funkygao/go-metrics"
	"github.com/influxdata/influxdb/client"
)

func createRunner() *runner {
	ctx.LoadFromHome()
	cf, _ := NewConfig("uri", "db", "user", "pass", time.Second)
	return New(metrics.DefaultRegistry, cf).(*runner)

}

func TestDump(t *testing.T) {
	g := metrics.NewRegisteredGaugeFloat64("gauge", nil)
	g.Update(1)

	r := createRunner()
	pts := make([]client.Point, 0, 1<<8)
	r.export(&pts)
	t.Logf("%+v", pts)
}

func TestExtractTags(t *testing.T) {
	r := createRunner()
	realName, tags := r.extractTagsFromMetricsName("pub.qps")
	t.Logf("%s %+v", realName, tags)
	assert.Equal(t, 1, len(tags)) // at least "host" attribute exists
	assert.Equal(t, "pub.qps", realName)

	realName, tags = r.extractTagsFromMetricsName("appid=5&topic=a.b.c&ver=v1#pub.qps")
	assert.Equal(t, "pub.qps", realName)
	for k, v := range tags {
		t.Logf("%s: %s", k, v)
	}

}

// 318 ns/op	     336 B/op	       2 allocs/op
func BenchmarkExtractTagsWithoutTags(b *testing.B) {
	r := createRunner()
	for i := 0; i < b.N; i++ {
		_, _ = r.extractTagsFromMetricsName("pub.qps")
	}
}

// 1854 ns/op	     784 B/op	       7 allocs/op
func BenchmarkExtractTagsWithTags(b *testing.B) {
	r := createRunner()
	for i := 0; i < b.N; i++ {
		_, _ = r.extractTagsFromMetricsName("appid=5&topic=a.b.c&ver=v1#pub.qps")
	}
}

func BenchmarkRawUrlParseQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url.ParseQuery("appid=5&topic=a.b.c&ver=v1#pub.qps")
	}
}
