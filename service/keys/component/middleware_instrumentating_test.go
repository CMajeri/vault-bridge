package components

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kit/kit/metrics"
	"github.com/influxdata/influxdb/pkg/testing/assert"
)

type mockCounter struct {
	CalledW *bool
	CalledA *bool
}

func (m *mockCounter) With(labelValues ...string) metrics.Counter {
	*(m.CalledW) = true
	fmt.Println(labelValues)
	return m
}

func (m *mockCounter) Add(delta float64) {
	*(m.CalledA) = true
	fmt.Println(delta)
}

type mockHistogram struct {
	CalledW *bool
	CalledO *bool
}

func (m *mockHistogram) With(labelValues ...string) metrics.Histogram {
	*(m.CalledW) = true
	fmt.Println(labelValues)
	return m
}

func (m *mockHistogram) Observe(value float64) {
	*(m.CalledO) = true
	fmt.Println(value)
}

type mockInfluxMetrics struct {
}

func (m *mockInfluxMetrics) NewCounter(name string) InfluxCounter {
	fmt.Println(name)
	var called = false
	var counter = &mockCounter{&called, &called}
	return counter
}

func (m *mockInfluxMetrics) NewHistogram(name string) InfluxHistogram {
	fmt.Println(name)
	var called = false
	var hist = &mockHistogram{&called, &called}
	return hist
}

func TestInstrumentatingMiddleware_ReadKey(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	_, err = m.ReadKey(context.Background(), pathKeyTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)
}

func TestInstrumentatingMiddleware_WriteKey(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	err = m.WriteKey(context.Background(), pathKeyTest, keyTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)

}

func TestInstrumentatingMiddleware_CreateKey(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	err = m.CreateKey(context.Background(), keyTest, paramsTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)
}

func TestInstrumentatingMiddleware_ExportKey(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	_, err = m.ExportKey(context.Background(), pathKeyTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)
}

func TestInstrumentatingMiddleware_Encrypt(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	_, err = m.Encrypt(context.Background(), keyTest, paramsTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)
}

func TestInstrumentatingMiddleware_Decrypt(t *testing.T) {
	var fail = false
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false
	var m = instrumentatingMiddleware{&mockCounter{&calledCountW, &calledCountA},
		&mockHistogram{&calledHistW, &calledHistO}, &mockServiceVault{fail}}
	var err error

	_, err = m.Decrypt(context.Background(), keyTest, paramsTest)
	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)
}

func TestMakeServiceInstrumentingMiddleware(t *testing.T) {
	var calledHistW = false
	var calledHistO = false
	var calledCountW = false
	var calledCountA = false

	var mCounter = mockCounter{&calledCountW, &calledCountA}
	var mHist = mockHistogram{&calledHistW, &calledHistO}
	var mService = mockServiceVault{false}
	var ns = MakeServiceInstrumentingMiddleware(&mCounter, &mHist)(&mService)

	ns.WriteKey(context.Background(), pathKeyTest, keyTest)
	assert.Equal(t, calledCountW, true)
	assert.Equal(t, calledCountA, true)
	assert.Equal(t, calledHistW, true)
	assert.Equal(t, calledHistO, true)

}
