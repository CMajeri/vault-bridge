package modules

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	vault "github.com/hashicorp/vault/api"
)

var ciphertextTest = "vault:v1:xC8gAgTcnP1qFTR5GuUaZJqdcsyv6k1VmxFJjbGHxsE="
var keyTest = "abcdfg"
var plaintextTest = "YWJjZA=="
var keysTest = map[string]interface{}{"1": "DXaRjwH4eup1XViCOyWyEW3aNGlfuO6Xjxeve/BEIiw="}
var respSecretTest = &vault.Secret{
	RequestID:     "91b4063d-19da-10d9-1879-8361af1656ee",
	LeaseID:       "",
	LeaseDuration: 2764800,
	Renewable:     false,
	Data: map[string]interface{}{"key": keyTest,
		"ciphertext": ciphertextTest,
		"plaintext":  plaintextTest,
		"keys":       keysTest,
	},
	Warnings: nil,
	Auth:     nil,
	WrapInfo: nil,
}
var jwtToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc"

type mockClient struct {
	fail bool
}

func (m *mockClient) Write(path string, data map[string]interface{}, token string) (*vault.Secret, error) {
	if m.fail == true {
		errFail := errors.New("Writing failed")
		return nil, errFail
	}

	return respSecretTest, nil
}

func (m *mockClient) Read(path string, token string) (*vault.Secret, error) {
	if m.fail == true {
		errFail := errors.New("Reading failed")
		return nil, errFail
	}

	return respSecretTest, nil
}

func (m *mockClient) CreatePolicy(path string, role string, policyName string) error {
	return nil
}

func (m *mockClient) CreateToken(policyName string) (string, error) {
	return "", nil
}

func TestBasicService_WriteKey(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var pathKey = "tenants/rolex/f1/key1"
	var keyValue = "123456"

	var ctx = context.WithValue(context.Background(), httptransport.ContextKeyRequestAuthorization, jwtToken)
	var errorWrite = bs.WriteKey(ctx, pathKey, keyValue)
	if errorWrite != nil {
		t.Errorf("I have %s and I should get nil", errorWrite)
	}
}

func TestBasicService_ReadKey(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var pathKey = "tenants/rolex/f1/key1"
	var ctx = context.WithValue(context.Background(), httptransport.ContextKeyRequestAuthorization, jwtToken)

	var key, errorRead = bs.ReadKey(ctx, pathKey)
	if errorRead != nil {
		t.Errorf("I have %s and I should get nil", errorRead)
	}

	if strings.Compare(key, keyTest) != 0 {
		t.Errorf("I have %s and I should get %s", key, keyTest)
	}
}

func TestBasicService_CreateKey(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var keyName = "key10"
	var params = map[string]interface{}{
		"type":       "aes256-gcm96",
		"derived":    false,
		"exportable": true,
	}

	var errorCreate = bs.CreateKey(context.Background(), keyName, params)
	if errorCreate != nil {
		t.Errorf("I have %s and I should get nil", errorCreate)
	}
}

func TestBasicService_ExportKey(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var keyPath = "encryption-key/key100/1"

	var keys, errorExport = bs.ExportKey(context.Background(), keyPath)
	if errorExport != nil {
		t.Errorf("I have %s and I should get nil", errorExport)
	}

	var eq = reflect.DeepEqual(keys, keysTest)
	if !eq {
		t.Errorf("I have %s and I should get %s", keys, keysTest)
	}
}

func TestBasicService_Encrypt(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var keyName = "key10"
	var params = map[string]interface{}{
		"plaintext":   "YWJjZA==",
		"key_version": 1,
	}

	var ciphertext, errorEncrypt = bs.Encrypt(context.Background(), keyName, params)
	if errorEncrypt != nil {
		t.Errorf("I have %s and I should get nil", errorEncrypt)
	}

	if strings.Compare(ciphertext, ciphertextTest) != 0 {
		t.Errorf("I have %s and I should get %s", ciphertext, ciphertextTest)
	}
}

func TestBasicService_Decrypt(t *testing.T) {
	client := mockClient{
		false,
	}
	var bs = NewBasicService(&client)
	var params = map[string]interface{}{
		"ciphertext": "vault:v1:dKG1C2bFFdLMJQkau6v3lDmhHfLtTMBB9cd0XBy6Id8=",
	}

	var plaintext, errorDecrypt = bs.Decrypt(context.Background(), ciphertextTest, params)
	if errorDecrypt != nil {
		t.Errorf("I have %s and I should get nil", errorDecrypt)
	}

	if strings.Compare(plaintext, plaintextTest) != 0 {
		t.Errorf("I have %s and I should get %s", plaintext, plaintextTest)
	}
}
