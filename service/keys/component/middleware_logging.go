package components

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
)

//Middleware type
type Middleware func(ServiceVault) ServiceVault

type loggingMiddleware struct {
	logger log.Logger
	next   ServiceVault
}

func (mw loggingMiddleware) WriteKey(ctx context.Context, pathKey string, keyValue string) (err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "writekey",
			"pathkey", pathKey,
			"keyvalue", keyValue,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.WriteKey(ctx, pathKey, keyValue)
	return
}

func (mw loggingMiddleware) ReadKey(ctx context.Context, pathKey string) (keyValue string, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "readkey",
			"pathkey", pathKey,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	keyValue, err = mw.next.ReadKey(ctx, pathKey)
	return
}

func (mw loggingMiddleware) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) (err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "createkey",
			"keyname", keyName,
			"params", paramsInfo,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.CreateKey(ctx, keyName, params)
	return
}

func (mw loggingMiddleware) ExportKey(ctx context.Context, keyPath string) (keys map[string]interface{}, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "exportkey",
			"keypath", keyPath,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	keys, err = mw.next.ExportKey(ctx, keyPath)
	return
}

func (mw loggingMiddleware) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (ciphertext string, err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "encrypt",
			"keyname", keyName,
			"params", paramsInfo,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	ciphertext, err = mw.next.Encrypt(ctx, keyName, params)
	return
}

func (mw loggingMiddleware) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (plaintext string, err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "decrypt",
			"keyname", keyName,
			"params", paramsInfo,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	plaintext, err = mw.next.Decrypt(ctx, keyName, params)
	return
}

//MakeServiceLoggingMiddleware wraps ServiceVault with logging tools
func MakeServiceLoggingMiddleware(log log.Logger) Middleware {
	return func(next ServiceVault) ServiceVault {
		return &loggingMiddleware{
			logger: log,
			next:   next,
		}
	}
}
