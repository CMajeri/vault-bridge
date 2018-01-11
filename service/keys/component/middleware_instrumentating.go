package components

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	influx "github.com/go-kit/kit/metrics/influx"
)

//InfluxCounter is the interface that implements the same functions as metrics.Counter
type InfluxCounter interface {
	With(labelValues ...string) metrics.Counter
	Add(delta float64)
}

//InfluxHistogram is the interface that implements the same functions as metrics.Histogram
type InfluxHistogram interface {
	With(labelValues ...string) metrics.Histogram
	Observe(value float64)
}

//InfluxMetrics uses a counter and a histogram
type InfluxMetrics interface {
	NewCounter(name string) InfluxCounter
	NewHistogram(name string) InfluxHistogram
}

type instrumentatingMiddleware struct {
	influxCounter   InfluxCounter
	influxHistogram InfluxHistogram
	next            ServiceVault
}

//InfluxClient struct
type InfluxClient struct {
	C *influx.Influx
}

//NewCounter instantiates an InfluxCounter for an InfluxClient
func (ic *InfluxClient) NewCounter(name string) InfluxCounter {
	return ic.C.NewCounter(name)
}

//NewHistogram instantiates an InfluxHistogram for an InfluxClient
func (ic *InfluxClient) NewHistogram(name string) InfluxHistogram {
	return ic.C.NewHistogram(name)
}

func (mw instrumentatingMiddleware) WriteKey(ctx context.Context, pathKey string, keyValue string) (err error) {
	defer func(begin time.Time) {
		var lvs = []string{"method", "writekey", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.WriteKey(ctx, pathKey, keyValue)
	return
}

func (mw instrumentatingMiddleware) ReadKey(ctx context.Context, pathKey string) (keyValue string, err error) {
	func(begin time.Time) {
		var lvs = []string{"method", "readkey", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	keyValue, err = mw.next.ReadKey(ctx, pathKey)
	return
}

func (mw instrumentatingMiddleware) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) (err error) {
	defer func(begin time.Time) {
		var lvs = []string{"method", "createkey", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.CreateKey(ctx, keyName, params)
	return
}

func (mw instrumentatingMiddleware) ExportKey(ctx context.Context, keyPath string) (keys map[string]interface{}, err error) {
	defer func(begin time.Time) {
		var lvs = []string{"method", "exportkey", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	keys, err = mw.next.ExportKey(ctx, keyPath)
	return
}

func (mw instrumentatingMiddleware) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (ciphertext string, err error) {
	defer func(begin time.Time) {
		var lvs = []string{"method", "encrypt", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ciphertext, err = mw.next.Encrypt(ctx, keyName, params)
	return
}

func (mw instrumentatingMiddleware) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (plaintext string, err error) {
	defer func(begin time.Time) {
		var lvs = []string{"method", "decrypt", "error", fmt.Sprint(err != nil)}
		mw.influxCounter.With(lvs...).Add(1)
		mw.influxHistogram.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	plaintext, err = mw.next.Decrypt(ctx, keyName, params)
	return
}

//MakeServiceInstrumentingMiddleware wraps ServiceVault with instrumenting tools
func MakeServiceInstrumentingMiddleware(influxCounter InfluxCounter, influxHistogram InfluxHistogram) Middleware {
	return func(next ServiceVault) ServiceVault {
		return &instrumentatingMiddleware{
			influxCounter:   influxCounter,
			influxHistogram: influxHistogram,
			next:            next,
		}
	}
}
