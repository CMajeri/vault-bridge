package modules

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

//Middleware type
type Middleware func(ServiceVault) ServiceVault

type tracingMiddleware struct {
	tracer stdopentracing.Tracer
	next   ServiceVault
}

func (mw tracingMiddleware) WriteKey(ctx context.Context, pathKey string, keyValue string) (err error) {
	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("writekey_module")
		span.LogFields(log.String("operation", "writekey"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("writekey_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "writekey"),
			log.String("microservice_level", "module"))

	}

	err = mw.next.WriteKey(ctx, pathKey, keyValue)
	return
}

func (mw tracingMiddleware) ReadKey(ctx context.Context, pathKey string) (keyValue string, err error) {

	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("readkey_module")
		span.LogFields(log.String("operation", "readkey"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("readkey_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "readkey"),
			log.String("microservice_level", "module"))

	}
	//defer span.Finish()
	keyValue, err = mw.next.ReadKey(ctx, pathKey)
	return
}

func (mw tracingMiddleware) CreateKey(ctx context.Context, keyName string, params map[string]interface{}) (err error) {

	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("createkey_module")
		span.LogFields(log.String("operation", "createkey"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("createkey_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "createkey"),
			log.String("microservice_level", "module"))

	}

	err = mw.next.CreateKey(ctx, keyName, params)
	return
}

func (mw tracingMiddleware) ExportKey(ctx context.Context, keyPath string) (keys map[string]interface{}, err error) {

	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("exportkey_module")
		span.LogFields(log.String("operation", "exportkey"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("exportkey_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "exportkey"),
			log.String("microservice_level", "module"))

	}

	keys, err = mw.next.ExportKey(ctx, keyPath)
	return
}

func (mw tracingMiddleware) Encrypt(ctx context.Context, keyName string, params map[string]interface{}) (ciphertext string, err error) {

	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("encrypt_module")
		span.LogFields(log.String("operation", "encrypt"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("encrypt_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "encrypt"),
			log.String("microservice_level", "module"))

	}

	ciphertext, err = mw.next.Encrypt(ctx, keyName, params)
	return
}

func (mw tracingMiddleware) Decrypt(ctx context.Context, keyName string, params map[string]interface{}) (plaintext string, err error) {

	var span = stdopentracing.SpanFromContext(ctx)

	if span == nil {
		// create a new root span
		span = mw.tracer.StartSpan("decrypt_module")
		span.LogFields(log.String("operation", "decrypt"),
			log.String("microservice_level", "module"))

	} else { // create child span
		var cspan = stdopentracing.StartSpan("decrypt_module", stdopentracing.ChildOf(span.Context()))
		defer cspan.Finish()

		cspan.LogFields(log.String("operation", "decrypt"),
			log.String("microservice_level", "module"))

	}

	plaintext, err = mw.next.Decrypt(ctx, keyName, params)
	return
}

//MakeServiceTracingMiddleware wraps ServiceVault with a tracer
func MakeServiceTracingMiddleware(tracer stdopentracing.Tracer) Middleware {
	return func(next ServiceVault) ServiceVault {
		return &tracingMiddleware{
			tracer: tracer,
			next:   next,
		}
	}
}
