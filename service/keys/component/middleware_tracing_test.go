package components

import (
	"context"
	"testing"

	"github.com/influxdata/influxdb/pkg/testing/assert"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func TestMakeServiceTracingMiddleware(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}

	var newService = MakeServiceTracingMiddleware(&mockTracer)(&mockService)
	newService.WriteKey(context.Background(), pathKeyTest, keyTest)

	assert.Equal(t, called, true)
}

func TestTracingMiddleware_WriteKey(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	err = mw.WriteKey(context.Background(), pathKeyTest, keyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestTracingMiddleware_ReadKey(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	_, err = mw.ReadKey(context.Background(), pathKeyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)

}

func TestTracingMiddleware_CreateKey(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	err = mw.CreateKey(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)

}

func TestTracingMiddleware_ExportKey(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	_, err = mw.ExportKey(context.Background(), pathKeyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)

}

func TestTracingMiddleware_Decrypt(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	_, err = mw.Decrypt(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)

}

func TestTracingMiddleware_Encrypt(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	_, err = mw.Encrypt(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

type mockTracer struct {
	Called *bool
}

func (m *mockTracer) StartSpan(operationName string, opts ...stdopentracing.StartSpanOption) stdopentracing.Span {
	*(m.Called) = true
	return &mockSpan{}
}

func (m *mockTracer) Inject(sm stdopentracing.SpanContext, format interface{}, carrier interface{}) error {
	return nil
}

func (m *mockTracer) Extract(format interface{}, carrier interface{}) (stdopentracing.SpanContext, error) {
	return nil, nil
}

type mockSpan struct{}

func (m *mockSpan) Finish()                                                        {}
func (m *mockSpan) FinishWithOptions(opts stdopentracing.FinishOptions)            {}
func (m *mockSpan) Context() stdopentracing.SpanContext                            { return nil }
func (m *mockSpan) SetOperationName(operationName string) stdopentracing.Span      { return nil }
func (m *mockSpan) SetTag(key string, value interface{}) stdopentracing.Span       { return nil }
func (m *mockSpan) LogFields(fields ...log.Field)                                  {}
func (m *mockSpan) LogKV(alternatingKeyValues ...interface{})                      {}
func (m *mockSpan) SetBaggageItem(restrictedKey, value string) stdopentracing.Span { return nil }
func (m *mockSpan) BaggageItem(restrictedKey string) string                        { return "" }
func (m *mockSpan) Tracer() stdopentracing.Tracer                                  { return nil }
func (m *mockSpan) LogEvent(event string)                                          {}
func (m *mockSpan) LogEventWithPayload(event string, payload interface{})          {}
func (m *mockSpan) Log(data stdopentracing.LogData)                                {}
