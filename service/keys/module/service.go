package modules

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	clientVault "github.com/cloustrust/vault-client/client"

	jwt "github.com/dgrijalva/jwt-go"
	httptransport "github.com/go-kit/kit/transport/http"
)

//ServiceVault interface
type ServiceVault interface {
	WriteKey(ctx context.Context, pathKey string, keyValue string) error
	ReadKey(ctx context.Context, pathKey string) (string, error)
	CreateKey(ctx context.Context, keyName string, params map[string]interface{}) error
	ExportKey(ctx context.Context, keyPath string) (map[string]interface{}, error)
	Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error)
	Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error)
}

type basicService struct {
	client clientVault.Client
}

//NewBasicService instantiates ServiceVault
func NewBasicService(client clientVault.Client) ServiceVault {
	return &basicService{
		client: client,
	}
}

func (bs *basicService) WriteKey(ctx context.Context, pathKey string, keyValue string) error {

	var err error
	var claims Claims
	//extract the JWT from the context and store it in a struct Claims
	claims, err = extractJWT(ctx)

	if err != nil {
		return err
	}

	var pathPolicy = "tenants/" + claims.Tenant + "/" + claims.Fil + "/*"
	var policyName = "writekey_" + claims.Tenant + "_" + claims.Fil
	var token string
	err = bs.client.CreatePolicy(pathPolicy, "writekey", policyName)
	if err != nil {
		return err
	}

	var errToken, errWrite error
	token, errToken = bs.client.CreateToken(policyName)
	if errToken != nil {
		return errToken
	}

	_, errWrite = bs.client.Write(pathKey, map[string]interface{}{"key": keyValue}, token)

	if errWrite != nil {
		return errWrite
	}
	return nil

}

func (bs *basicService) ReadKey(ctx context.Context, pathKey string) (string, error) {

	var err error
	var claims Claims
	//extract the JWT from the context and store it in the struct Claims
	claims, err = extractJWT(ctx)

	if err != nil {
		return "", err
	}

	var pathPolicy = "tenants/" + claims.Tenant + "/" + claims.Fil + "/*"
	var policyName = "readkey_" + claims.Tenant + "_" + claims.Fil
	err = bs.client.CreatePolicy(pathPolicy, "readkey", policyName)
	if err != nil {
		return "", err
	}

	var token string
	var errToken error
	token, errToken = bs.client.CreateToken(policyName)
	if errToken != nil {
		return "", errToken
	}

	var secret, errRead = bs.client.Read(pathKey, token)

	if errRead != nil {
		return "", errRead
	}

	if secret == nil {
		return "", errors.New("Searched key could not be found")
	}

	// we have the convention that the key value is at secret.Data["key"]
	var keyValue string
	var ok bool
	if keyValue, ok = secret.Data["key"].(string); ok {
		return keyValue, nil
	}

	return "", errors.New("Read key: the key read from Vault is not of type string or it does not exist")

}

func (bs *basicService) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) error {
	// transit backend must be mounted before we can perform this operation
	//create a key in the transit backend

	var policyName = "createkey"
	var pathPolicy = "transit/keys/*"
	var err error

	err = bs.client.CreatePolicy(pathPolicy, "createkey", policyName)
	if err != nil {
		return err
	}

	var token string
	var errToken, errKey error
	token, errToken = bs.client.CreateToken(policyName)
	if errToken != nil {
		return errToken
	}

	_, errKey = bs.client.Write("transit/keys/"+keyName, params, token)
	if errKey != nil {
		return errKey
	}
	return nil
}

func (bs *basicService) ExportKey(ctx context.Context, keyPath string) (map[string]interface{}, error) {
	// transit backend must be mounted before we can perform this operation

	var policyName = "exportkey"
	var pathPolicy = "transit/export/encryption-key/*"
	var err error

	err = bs.client.CreatePolicy(pathPolicy, "exportkey", policyName)
	if err != nil {
		return nil, err
	}

	var token string
	var errToken error
	token, errToken = bs.client.CreateToken("exportkey")
	if errToken != nil {
		return nil, errToken
	}

	var secret, errRead = bs.client.Read("transit/export/"+keyPath, token)
	if errRead != nil {
		return nil, errRead
	}

	if secret == nil {
		return nil, errors.New("Encryption key could not be found")
	}

	var keys map[string]interface{}
	var ok bool
	if keys, ok = secret.Data["keys"].(map[string]interface{}); ok {
		return keys, nil
	}
	return nil, errors.New("Export key: the keys data structure is not of the correct type \n Expected map[string]interface{}")

}

func (bs *basicService) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	//encrypt in the transit backend
	//transit backend must be mounted before we can perform this operation

	var policyName = "encrypt"
	var pathPolicy = "transit/encrypt/*"
	var err error

	err = bs.client.CreatePolicy(pathPolicy, "encrypt", policyName)
	if err != nil {
		return "", err
	}

	var token string
	var errToken error
	token, errToken = bs.client.CreateToken(policyName)
	if errToken != nil {
		return "", errToken
	}

	var ciphertextData, errEncrypt = bs.client.Write("transit/encrypt/"+keyName, params, token)
	if errEncrypt != nil {
		return "", errEncrypt
	}

	var ciphertext string
	var ok bool
	if ciphertext, ok = ciphertextData.Data["ciphertext"].(string); ok {
		return ciphertext, nil
	}
	return "", errors.New("Encryption failed")

}

func (bs *basicService) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	//decrypt in the transit backend
	// transit backend must be mounted before we can perform this operation

	var policyName = "decrypt"
	var pathPolicy = "transit/decrypt/*"
	var err error

	err = bs.client.CreatePolicy(pathPolicy, "decrypt", policyName)
	if err != nil {
		return "", err
	}

	var token string
	var errToken error
	token, errToken = bs.client.CreateToken(policyName)
	if errToken != nil {
		return "", errToken
	}

	var plaintextData, errDecrypt = bs.client.Write("transit/decrypt/"+keyName, params, token)
	if errDecrypt != nil {
		return "", errDecrypt
	}

	var plaintext string
	var ok bool
	if plaintext, ok = plaintextData.Data["plaintext"].(string); ok {
		return plaintext, nil
	}
	return "", errors.New("Decryption failed")

}

func decodeJWT(token string) ([]byte, error) {
	var parts []string
	parts = strings.Split(token, ".")

	var claims, err = jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func extractJWT(ctx context.Context) (Claims, error) {

	var token string
	token = ctx.Value(httptransport.ContextKeyRequestAuthorization).(string)

	var err error
	var decodedJWT []byte
	decodedJWT, err = decodeJWT(token)

	var claims Claims
	err = json.Unmarshal(decodedJWT, &claims)
	if err != nil {
		return claims, err
	}
	return claims, nil
}

//Claims format
type Claims struct {
	Tenant string `json:"tenant"`
	Fil    string `json:"fil"`
}
