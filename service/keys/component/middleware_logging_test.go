package components

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/influxdata/influxdb/pkg/testing/assert"
	"github.com/pkg/errors"
)

var keyTest = "dummy key"
var pathKeyTest = "dummy path"
var plainTest = "dummy plaintext"
var cipherTest = "dummy ciphertext"
var paramsTest = map[string]interface{}{
	"plaintext":  plainTest,
	"ciphertext": cipherTest,
}

func TestLoggingMiddleware_ReadKey(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var key string
	var err error

	key, err = mw.ReadKey(context.Background(), pathKeyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	if strings.Compare(key, keyTest) != 0 {
		t.Errorf("I have %s and I should get %s", key, keyTest)
	}
	assert.Equal(t, called, true)
}

func TestLoggingMiddleware_WriteKey(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var err error

	err = mw.WriteKey(context.Background(), pathKeyTest, keyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestLoggingMiddleware_CreateKey(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var err error

	err = mw.CreateKey(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestLoggingMiddleware_ExportKey(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var err error

	_, err = mw.ExportKey(context.Background(), pathKeyTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestLoggingMiddleware_Encrypt(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var err error

	_, err = mw.Encrypt(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestLoggingMiddleware_Decrypt(t *testing.T) {
	var called = false
	var mw = loggingMiddleware{&mockLogging{&called}, &mockServiceVault{false}}
	var err error

	_, err = mw.Decrypt(context.Background(), keyTest, paramsTest)

	if err != nil {
		t.Errorf("I have %s and I should get nil", err)
	}
	assert.Equal(t, called, true)
}

func TestMakeServiceLoggingMiddleware(t *testing.T) {
	var called = false
	var mServiceVault = &mockServiceVault{false}
	var mLog = &mockLogging{&called}
	var newService = MakeServiceLoggingMiddleware(mLog)(mServiceVault)

	newService.WriteKey(context.Background(), pathKeyTest, keyTest)

	assert.Equal(t, called, true)

}

type mockLogging struct {
	Called *bool
}

func (m *mockLogging) Log(keyvals ...interface{}) error {
	*(m.Called) = true
	fmt.Println(keyvals)
	return nil
}

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
		return keyTest, errors.New("Read key failed")
	}
	return keyTest, nil

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
