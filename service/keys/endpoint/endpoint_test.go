package endpoints

import (
	"context"
	"errors"
	"testing"
)

func TestMakeWriteKeyEndpoint(t *testing.T) {
	var service = &mockService{
		false,
	}

	var makeWriteKeyEndpoint = MakeWriteKeyEndpoint(service)
	var req = WritekeyRequest{
		KeyValue: "dummy_key",
		PathKey:  "dummy_path",
	}
	var resp, _ = makeWriteKeyEndpoint(context.Background(), req)
	var r = resp.(WritekeyResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}

}

func TestMakeReadKeyEndpoint(t *testing.T) {
	service := &mockService{
		false,
	}
	var makeReadKeyEndpoint = MakeReadKeyEndpoint(service)
	var req = ReadkeyRequest{
		PathKey: "dummy_path",
	}
	var resp, _ = makeReadKeyEndpoint(context.Background(), req)
	var r = resp.(ReadkeyResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}
}

func TestMakeCreateKeyEndpoint(t *testing.T) {
	service := &mockService{
		false,
	}
	var makeCreateKeyEndpoint = MakeCreateKeyEndpoint(service)
	var req = CreatekeyRequest{
		KeyName:    "dummy_name",
		Parameters: map[string]interface{}{"dummy": "dummy"},
	}
	var resp, _ = makeCreateKeyEndpoint(context.Background(), req)
	var r = resp.(CreatekeyResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}
}

func TestMakeExportKeyEndpoint(t *testing.T) {
	service := &mockService{
		false,
	}
	var makeExportKeyEndpoint = MakeExportKeyEndpoint(service)
	var req = ExportkeyRequest{
		KeyPath: "dummy_path",
	}
	var resp, _ = makeExportKeyEndpoint(context.Background(), req)
	var r = resp.(ExportkeyResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}
}

func TestMakeEncryptEndpoint(t *testing.T) {
	service := &mockService{
		false,
	}
	var makeEncryptEndpoint = MakeEncryptEndpoint(service)
	var req = EncryptRequest{
		KeyName:    "dummy_name",
		Parameters: map[string]interface{}{"dummy": "dummy"},
	}
	var resp, _ = makeEncryptEndpoint(context.Background(), req)
	var r = resp.(EncryptResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}
}

func TestMakeDecryptEndpoint(t *testing.T) {
	service := &mockService{
		false,
	}
	var makeDecryptEndpoint = MakeDecryptEndpoint(service)
	var req = DecryptRequest{
		KeyName:    "dummy_name",
		Parameters: map[string]interface{}{"dummy": "dummy"},
	}
	var resp, _ = makeDecryptEndpoint(context.Background(), req)
	var r = resp.(DecryptResponse)
	if r.Err != "" {
		t.Errorf("I have %s and I should get nil", r.Err)
	}
}

type mockService struct {
	fail bool
}

func (ms *mockService) WriteKey(ctx context.Context, pathKey string, keyValue string) error {
	if ms.fail == true {
		return errors.New("Writing failed")
	}
	return nil
}

func (ms *mockService) ReadKey(ctx context.Context, pathKey string) (string, error) {
	if ms.fail == true {
		error := errors.New("Reading failed")
		return "", error
	}
	return "dummy_key", nil
}
func (ms *mockService) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) error {
	if ms.fail == true {
		return errors.New("Creating key failed")
	}
	return nil
}
func (ms *mockService) ExportKey(ctx context.Context, keyPath string) (map[string]interface{}, error) {
	if ms.fail == true {
		error := errors.New("Exporting key failed")
		return map[string]interface{}{"dummy": "dummy"}, error
	}
	return map[string]interface{}{"dummy": "dummy"}, nil
}
func (ms *mockService) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	if ms.fail == true {
		error := errors.New("Encryption failed")
		return "", error
	}
	return "dummy_ciphertext", nil
}
func (ms *mockService) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	if ms.fail == true {
		error := errors.New("Decryption failed")
		return "", error
	}
	return "dummy_plaintext", nil
}
