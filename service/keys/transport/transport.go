package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	endpoint "github.com/cloudtrust/vault-bridge/service/keys/endpoint"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// key used for signing the JWT token
var keyJWT = "secret"

//DecodeWriteKeyRequest processes the write key request
func DecodeWriteKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the path of the key
	var path = strings.TrimPrefix(r.URL.Path, "/key/")

	var request endpoint.WritekeyRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	request.PathKey = path
	return request, nil
}

//DecodeReadKeyRequest processes the read key request
func DecodeReadKeyRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the path of the key
	var path = strings.TrimPrefix(r.URL.Path, "/key/")

	var request = endpoint.ReadkeyRequest{
		PathKey: path,
	}
	return request, nil
}

//DecodeCreateKeyRequest processes the create key request
func DecodeCreateKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the name of the key that we will create
	var keyName = strings.TrimPrefix(r.URL.Path, "/createkey/")

	var request endpoint.CreatekeyRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	request.KeyName = keyName
	return request, nil
}

//DecodeExportKeyRequest processes the export key request
func DecodeExportKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the path of the key we will export
	var path = strings.TrimPrefix(r.URL.Path, "/exportkey/")

	var request = endpoint.ExportkeyRequest{
		KeyPath: path,
	}
	return request, nil
}

//DecodeEncryptRequest processes the encrypt request
func DecodeEncryptRequest(_ context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the name of the key used for encryption
	var keyName = strings.TrimPrefix(r.URL.Path, "/encrypt/")

	var request endpoint.EncryptRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	request.KeyName = keyName

	return request, nil
}

//DecodeDecryptRequest processes the decrypt request
func DecodeDecryptRequest(_ context.Context, r *http.Request) (interface{}, error) {

	//verify the signature of the Jason Web Token
	var err = verifyJWT(r)

	if err != nil {
		return nil, errors.New("Verifying the signature of the JWT: " + err.Error())
	}

	//extract the name of the key used for decryption
	var keyName = strings.TrimPrefix(r.URL.Path, "/decrypt/")

	var request endpoint.DecryptRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	request.KeyName = keyName
	return request, nil
}

//EncodeResponse encodes the passed response object to the HTTP response writer
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

//verify the signature of the JWT
func verifyJWT(r *http.Request) error {

	var jwtToken string
	jwtToken = r.Header.Get("Authorization")

	var _, err = jwt.Parse(jwtToken[7:], func(token *jwt.Token) (interface{}, error) {
		return []byte(keyJWT), nil
	})

	if err != nil {
		return errors.New("Verifying the signature of the JWT: " + err.Error())
	}
	return nil
}
