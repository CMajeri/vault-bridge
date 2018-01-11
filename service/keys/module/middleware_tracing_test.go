package modules

import (
	"context"
	"testing"

	"github.com/influxdata/influxdb/pkg/testing/assert"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

var keyTestName = "dummy key"
var keyTestValue = "dummy value"
var pathKeyTest = "dummy path"
var plainTest = "dummy plaintext"
var cipherTest = "dummy ciphertext"
var paramsTest = map[string]interface{}{
	"plaintext":  plainTest,
	"ciphertext": cipherTest,
}

func TestMakeServiceTracingMiddleware(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var newService = MakeServiceTracingMiddleware(&mockTracer)(&mockService)

	newService.WriteKey(context.Background(), pathKeyTest, keyTestValue)
	assert.Equal(t, called, true)
}

func TestTracingMiddleware_WriteKey(t *testing.T) {
	var called = false
	var mockService = mockServiceVault{}
	var mockTracer = mockTracer{&called}
	var err error
	var mw = tracingMiddleware{&mockTracer, &mockService}

	err = mw.WriteKey(context.Background(), pathKeyTest, keyTestValue)

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

	err = mw.CreateKey(context.Background(), keyTestName, paramsTest)

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

	_, err = mw.Decrypt(context.Background(), keyTestName, paramsTest)

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

	_, err = mw.Encrypt(context.Background(), keyTestName, paramsTest)

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

type mockServiceVault struct {
	Fail bool
}

func (m *mockServiceVault) WriteKey(ctx context.Context, pathKey string, keyValue string) error {
	if m.Fail == true {
		return errors.New("Write key failed")

	}
	return nil
}

func (m *mockServiceVault) ReadKey(ctx context.Context, pathKey string) (string, error) {
	if m.Fail == true {
		return keyTestValue, errors.New("Read key failed")

	}
	return keyTestValue, nil

}

func (m *mockServiceVault) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) error {
	if m.Fail == true {
		return errors.New("Create key failed")
	}
	return nil
}

func (m *mockServiceVault) ExportKey(ctx context.Context, keyPath string) (map[string]interface{}, error) {
	if m.Fail == true {
		return map[string]interface{}{}, errors.New("Export key failed")
	}
	return map[string]interface{}{}, nil
}

func (m *mockServiceVault) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	if m.Fail == true {
		return cipherTest, errors.New("Encrypt failed")
	}
	return cipherTest, nil
}

func (m *mockServiceVault) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	if m.Fail == true {
		return plainTest, errors.New("Decrypt failed")
	}
	return plainTest, nil
}
