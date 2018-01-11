package transport

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

//HTTPToContext middleware at the transport level
func HTTPToContext(tracer stdopentracing.Tracer, operationName string, logger log.Logger) httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		// Try to join to a trace propagated in `req`.
		var span stdopentracing.Span
		var wireContext, err = tracer.Extract(
			stdopentracing.TextMap,
			stdopentracing.HTTPHeadersCarrier(req.Header),
		)
		if err != nil && err != stdopentracing.ErrSpanContextNotFound {
			logger.Log("error", err)
		}

		span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
		span.LogFields(otlog.String("operation", operationName),
			otlog.String("microservice_level", "transport"))
		//defer span.Finish()
		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL.String())
		return stdopentracing.ContextWithSpan(ctx, span)
	}
}
