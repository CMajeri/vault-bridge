package components

import (
	"context"
	module "vault-bridge/service/keys/module"
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
	module module.ServiceVault
}

//NewBasicService instatiates ServiceVault
func NewBasicService(module module.ServiceVault) ServiceVault {
	return &basicService{
		module: module,
	}
}

func (bs *basicService) WriteKey(ctx context.Context, pathKey string, keyValue string) error {
	return bs.module.WriteKey(ctx, pathKey, keyValue)
}

func (bs *basicService) ReadKey(ctx context.Context, pathKey string) (string, error) {
	return bs.module.ReadKey(ctx, pathKey)
}

func (bs *basicService) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) error {
	return bs.module.CreateKey(ctx, keyName, params)
}

func (bs *basicService) ExportKey(ctx context.Context, keyPath string) (map[string]interface{}, error) {
	return bs.module.ExportKey(ctx, keyPath)
}

func (bs *basicService) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	return bs.module.Encrypt(ctx, keyName, params)
}

func (bs *basicService) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (string, error) {
	return bs.module.Decrypt(ctx, keyName, params)
}
