package components

import (
	"context"
	"fmt"

	sentry "github.com/getsentry/raven-go"
	"github.com/go-kit/kit/log"
)

type sentryClient interface {
	CaptureErrorAndWait(err error, tags map[string]string, interfaces ...sentry.Interface) string
}

type errorMiddleware struct {
	log    log.Logger
	client sentryClient
	next   ServiceVault
}

func (mw errorMiddleware) WriteKey(ctx context.Context, pathKey string, keyValue string) (err error) {

	err = mw.next.WriteKey(ctx, pathKey, keyValue)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "writekey", "pathkey": pathKey, "keyvalue": keyValue})
	}
	return
}

func (mw errorMiddleware) ReadKey(ctx context.Context, pathKey string) (keyValue string, err error) {

	keyValue, err = mw.next.ReadKey(ctx, pathKey)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "readkey", "pathkey": pathKey})
	}
	return
}

func (mw errorMiddleware) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) (err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	err = mw.next.CreateKey(ctx, keyName, params)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "createkey", "keyname": keyName, "parameters": paramsInfo})
	}
	return
}

func (mw errorMiddleware) ExportKey(ctx context.Context, keyPath string) (keys map[string]interface{}, err error) {

	keys, err = mw.next.ExportKey(ctx, keyPath)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "exportkey", "keypath": keyPath})
	}
	return
}

func (mw errorMiddleware) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (ciphertext string, err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	ciphertext, err = mw.next.Encrypt(ctx, keyName, params)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "encrypt", "parameters": paramsInfo})
	}
	return
}

func (mw errorMiddleware) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (plaintext string, err error) {

	var paramsInfo string
	for k, v := range params {
		paramsInfo = paramsInfo + fmt.Sprint(k, " -> ", v, " , ")
	}

	plaintext, err = mw.next.Decrypt(ctx, keyName, params)
	if err != nil {
		mw.log.Log("error_to_Sentry", err)
		mw.client.CaptureErrorAndWait(err, map[string]string{"method": "decrypt", "parameters": paramsInfo})
	}
	return
}

//MakeServiceErrorMiddleware wraps ServiceVault with Sentry
func MakeServiceErrorMiddleware(log log.Logger, client sentryClient) Middleware {
	return func(next ServiceVault) ServiceVault {
		return &errorMiddleware{
			log:    log,
			client: client,
			next:   next,
		}
	}
}
