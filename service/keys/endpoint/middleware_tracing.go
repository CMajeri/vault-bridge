package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	stdopentracing "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

//MakeEndpointTracingMiddleware wraps Endpoint with a tracer
func MakeEndpointTracingMiddleware(tracer stdopentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var span = stdopentracing.SpanFromContext(ctx)
			if span == nil {
				//create a new root span.
				span = tracer.StartSpan(operationName)

				span.LogFields(log.String("operation", operationName),
					log.String("microservice_level", "endpoint"))

				otext.SpanKindRPCServer.Set(span)
				ctx = stdopentracing.ContextWithSpan(ctx, span)
				return next(ctx, request)
			}
			cspan := stdopentracing.StartSpan(operationName, stdopentracing.ChildOf(span.Context()))
			//defer cspan.Finish()
			defer span.Finish()
			cspan.LogFields(log.String("operation", operationName),
				log.String("microservice_level", "endpoint"))

			otext.SpanKindRPCServer.Set(cspan)
			ctx = stdopentracing.ContextWithSpan(ctx, cspan)
			return next(ctx, request)
		}
	}
}
