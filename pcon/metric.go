package main

import (
	"expvar"
	"net/http"
	"runtime"
	"time"

	metric "github.com/weareplanet/ifcv5-main/utils/metrics"
)

func initMetric(addr string, name string) {

	if len(addr) == 0 {
		return
	}

	goroutine := metric.NewGauge("2m5s", "15m30s", "1h1m")
	expvar.Publish("go:numgoroutine", goroutine)

	alloc := metric.NewGauge("2m5s", "15m30s", "1h1m")
	expvar.Publish("go:alloc:mb", alloc)

	inuseheap := metric.NewGauge("2m5s", "15m30s", "1h1m")
	expvar.Publish("go:inuseheap:mb", inuseheap)

	inusestack := metric.NewGauge("2m5s", "15m30s", "1h1m")
	expvar.Publish("go:inusestack:mb", inusestack)

	expvar.Publish("noise:send:kb", metric.NewCounter("2m1s", "15m30s", "1h1m"))
	expvar.Publish("noise:read:kb", metric.NewCounter("2m1s", "15m30s", "1h1m"))

	latency := metric.NewGauge("2m1s", "15m30s", "1h1m")
	expvar.Publish("noise:latency:ms", latency)

	go func() {
		for range time.Tick(5000 * time.Millisecond) {
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)
			goroutine.Add(float64(runtime.NumGoroutine()))
			alloc.Add(float64(m.Alloc) / 1000000)
			inuseheap.Add(float64(m.HeapInuse) / 1000000)
			inusestack.Add(float64(m.StackInuse) / 1000000)
		}
	}()

	http.Handle("/debug/metrics", metric.Handler(metric.Exposed))

	go func() {

		defer func() {
			metric.Enabled = false
		}()

		metric.Enabled = true
		metric.Title = name
		http.ListenAndServe(addr, nil)

	}()
}
