package components

import (
	"context"
	"fmt"
	"testing"

	sentry "github.com/getsentry/raven-go"
	"github.com/influxdata/influxdb/pkg/testing/assert"
)

func TestErrorMiddleware_ReadKey(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_, _ = mw.ReadKey(context.Background(), pathKeyTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)

}

func TestErrorMiddleware_WriteKey(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_ = mw.WriteKey(context.Background(), pathKeyTest, keyTest)

	//If err == nil then I should not have called Sentry and would not log anything
	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)
}

func TestErrorMiddleware_CreateKey(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_ = mw.CreateKey(context.Background(), keyTest, paramsTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)

}

func TestErrorMiddleware_ExportKey(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_, _ = mw.ExportKey(context.Background(), pathKeyTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)

}

func TestErrorMiddleware_Encrypt(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_, _ = mw.Encrypt(context.Background(), keyTest, paramsTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)

}

func TestErrorMiddleware_Decrypt(t *testing.T) {
	var calledL = false
	var calledS = false

	var mw = errorMiddleware{&mockLogging{&calledL}, &mockSentry{&calledS}, &mockServiceVault{true}}
	_, _ = mw.Decrypt(context.Background(), keyTest, paramsTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)
}

func TestMakeServiceErrorMiddleware(t *testing.T) {
	var calledL = false
	var calledS = false
	var mServiceVault = &mockServiceVault{true}
	var mSentryClient = &mockSentry{&calledS}
	var mLog = &mockLogging{&calledL}

	var newService = MakeServiceErrorMiddleware(mLog, mSentryClient)(mServiceVault)
	newService.WriteKey(context.Background(), pathKeyTest, keyTest)

	assert.Equal(t, calledS, true)
	assert.Equal(t, calledL, true)
}

type mockSentry struct {
	Called *bool
}

func (m *mockSentry) CaptureErrorAndWait(err error, tags map[string]string, interfaces ...sentry.Interface) string {
	*(m.Called) = true
	fmt.Println(err)
	return ""
}
