package endpoints

import (
	"context"

	service "github.com/cloudtrust/vault-bridge/service/keys/component"

	"github.com/go-kit/kit/endpoint"
)

//MakeWriteKeyEndpoint calls the ServiceVault in order to process the write key request
func MakeWriteKeyEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(WritekeyRequest)

		var err = svc.WriteKey(ctx, req.PathKey, req.KeyValue)
		if err != nil {
			return WritekeyResponse{err.Error()}, nil
		}
		return WritekeyResponse{""}, nil
	}
}

//MakeReadKeyEndpoint calls the ServiceVault in order to process the read key request
func MakeReadKeyEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(ReadkeyRequest)
		var respReadKey string
		var err error

		respReadKey, err = svc.ReadKey(ctx, req.PathKey)
		if err != nil {
			return ReadkeyResponse{"", err.Error()}, nil
		}
		return ReadkeyResponse{respReadKey, ""}, nil
	}
}

//MakeCreateKeyEndpoint calls the ServiceVault in order to process the create key request
func MakeCreateKeyEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(CreatekeyRequest)
		var err error

		err = svc.CreateKey(ctx, req.KeyName, req.Parameters)
		if err != nil {
			return CreatekeyResponse{err.Error()}, nil
		}
		return CreatekeyResponse{""}, nil
	}
}

//MakeExportKeyEndpoint calls the ServiceVault in order to process the export key request
func MakeExportKeyEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(ExportkeyRequest)
		var keys map[string]interface{}
		var err error

		keys, err = svc.ExportKey(ctx, req.KeyPath)
		if err != nil {
			return ExportkeyResponse{nil, err.Error()}, nil
		}
		return ExportkeyResponse{keys, ""}, nil
	}
}

//MakeEncryptEndpoint calls the ServiceVault in order to process the encrypt request
func MakeEncryptEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(EncryptRequest)
		var ciphertext string
		var err error

		ciphertext, err = svc.Encrypt(ctx, req.KeyName, req.Parameters)
		if err != nil {
			return EncryptResponse{"", err.Error()}, nil
		}
		return EncryptResponse{ciphertext, ""}, nil
	}
}

//MakeDecryptEndpoint calls the ServiceVault in order to process the decrypt request
func MakeDecryptEndpoint(svc service.ServiceVault) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var req = request.(DecryptRequest)
		var plaintext string
		var err error

		plaintext, err = svc.Decrypt(ctx, req.KeyName, req.Parameters)
		if err != nil {
			return DecryptResponse{"", err.Error()}, nil
		}
		return DecryptResponse{plaintext, ""}, nil
	}
}

//WritekeyRequest format
type WritekeyRequest struct {
	PathKey  string `json:"pathk"`
	KeyValue string `json:"key"`
}

//WritekeyResponse format
type WritekeyResponse struct {
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

//ReadkeyRequest format
type ReadkeyRequest struct {
	PathKey string `json:"pathk"`
}

//ReadkeyResponse format
type ReadkeyResponse struct {
	KeyValue string `json:"key"`
	Err      string `json:"err,omitempty"`
}

//CreatekeyRequest format
type CreatekeyRequest struct {
	KeyName    string                 `json:"keyname"`
	Parameters map[string]interface{} `json:"params"`
}

//CreatekeyResponse format
type CreatekeyResponse struct {
	Err string `json:"err,omitempty"`
}

//ExportkeyRequest format
type ExportkeyRequest struct {
	KeyPath string `json:"keypath"`
}

//ExportkeyResponse format
type ExportkeyResponse struct {
	KeyValue map[string]interface{} `json:"keys"`
	Err      string                 `json:"err,omitempty"`
}

//EncryptRequest format
type EncryptRequest struct {
	Parameters map[string]interface{} `json:"params"`
	KeyName    string                 `json:"keyname"`
}

//EncryptResponse format
type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
	Err        string `json:"err,omitempty"`
}

//DecryptRequest format
type DecryptRequest struct {
	Parameters map[string]interface{} `json:"params"`
	KeyName    string                 `json:"keyname"`
}

//DecryptResponse format
type DecryptResponse struct {
	Plaintext string `json:"plaintext"`
	Err       string `json:"err,omitempty"`
}
