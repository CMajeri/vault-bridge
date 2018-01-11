package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	service "github.com/cloustrust/vault-bridge/service/keys/component"
	endpoint "github.com/cloustrust/vault-bridge/service/keys/endpoint"
	module "github.com/cloustrust/vault-bridge/service/keys/module"
	transport "github.com/cloustrust/vault-bridge/service/keys/transport"
	vaultClient "github.com/cloustrust/vault-client/client"

	sentry "github.com/getsentry/raven-go"
	"github.com/go-kit/kit/log"
	influx "github.com/go-kit/kit/metrics/influx"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	influxdb "github.com/influxdata/influxdb/client/v2"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var (
	// Version of the component.
	Version = "1.0.0"
	// Environment is filled by the compiler.
	Environment = "unknown"
	// GitCommit is filled by the compiler.
	GitCommit = "unknown"
)

func main() {

	//Logger
	var logger = log.NewLogfmtLogger(os.Stdout)
	{
		logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
		defer logger.Log("msg", "Goodbye")
	}

	//Configurations
	config := config(log.With(logger, "component", "config_loader"))
	var (
		influxHTTPConfig = influxdb.HTTPConfig{
			Addr:     config["influx-url"].(string),
			Username: config["influx-username"].(string),
			Password: config["influx-password"].(string),
		}
		influxBatchPointsConfig = influxdb.BatchPointsConfig{
			Precision:        config["influx-precision"].(string),
			Database:         config["influx-database"].(string),
			RetentionPolicy:  config["influx-retention-policy"].(string),
			WriteConsistency: config["influx-write-consistency"].(string),
		}
		influxDBName = fmt.Sprintf(config["influx-database"].(string))
		tableCount   = config["influx-counter-table-name"].(string)
		tableHist    = config["influx-histogram-table-name"].(string)

		httpAddr  = fmt.Sprintf(config["component-http-address"].(string))
		sentryDSN = fmt.Sprintf(config["sentry-dsn"].(string))

		vaultToken = fmt.Sprintf(config["vault-token"].(string))
		vaultURL   = fmt.Sprintf(config["vault-url"].(string))

		jaegerConfiguration = jaegerConfig.Configuration{
			Sampler: &jaegerConfig.SamplerConfig{
				Type:              config["jaeger-sampler-type"].(string),
				Param:             float64(config["jaeger-sampler-param"].(int)),
				SamplingServerURL: config["jaeger-sampler-url"].(string),
			},
			Reporter: &jaegerConfig.ReporterConfig{
				LogSpans:            config["jaeger-reporter-logspan"].(bool),
				BufferFlushInterval: time.Duration(config["jaeger-reporter-flushinterval-ms"].(int)) * time.Millisecond,
			},
		}
		jaegerName = config["jaeger-service-name"].(string)
	)

	// Vault client
	var vClient vaultClient.Client
	{
		var err error
		vClient, err = vaultClient.NewClient(vaultToken, vaultURL)
		if err != nil {
			panic(err)
		}
	}

	//Tracer
	var tracer stdopentracing.Tracer
	var closer io.Closer
	{
		var err error
		tracer, closer, err = jaegerConfiguration.New(
			jaegerName,
		)
		if err != nil {
			logger.Log("error", err)
		}
	}
	stdopentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	//Module
	var moduleVault module.ServiceVault
	moduleVault = module.NewBasicService(vClient)
	moduleVault = module.MakeServiceTracingMiddleware(tracer)(moduleVault)

	//Service
	var svc service.ServiceVault
	svc = service.NewBasicService(moduleVault)
	svc = service.MakeServiceLoggingMiddleware(logger)(svc)
	svc = service.MakeServiceTracingMiddleware(tracer)(svc)

	//Sentry client
	var sentryClient *sentry.Client
	{
		var logger = log.With(logger, "service", "Sentry", "config", sentryDSN)
		var err error
		sentryClient, err = sentry.New(sentryDSN)
		if err != nil {
			logger.Log("error", err)
			return
		}
	}
	svc = service.MakeServiceErrorMiddleware(logger, sentryClient)(svc)

	//InfluxDB
	var influxClient influxdb.Client
	{
		var logger = log.With(logger, "service", "Influx")
		{
			var err error
			influxClient, err = influxdb.NewHTTPClient(influxHTTPConfig)
			if err != nil {
				logger.Log("error", err)
				return
			}
		}
	}

	var clientInfluxKit *influx.Influx
	var influxCounter service.InfluxCounter
	var influxHistogram service.InfluxHistogram
	var client service.InfluxClient
	{
		var logger = log.With(logger, "service", "Influx", "database_name", influxDBName)
		clientInfluxKit = influx.New(map[string]string{}, influxBatchPointsConfig, logger)
		client = service.InfluxClient{C: clientInfluxKit}
		influxCounter = client.NewCounter(tableCount)
		influxHistogram = client.NewHistogram(tableHist)
	}
	svc = service.MakeServiceInstrumentingMiddleware(influxCounter, influxHistogram)(svc)
	//svc = service.MakeServiceInstrumentingMiddleware(influxCounter, influxHistogram, &client)(svc)

	//Influx Handling
	go func() {
		var tic = time.NewTicker(1 * time.Millisecond)
		clientInfluxKit.WriteLoop(tic.C, influxClient)
	}()

	//Http handlers
	var writekeyHandler = httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "writekey_endpoint")(endpoint.MakeWriteKeyEndpoint(svc)),
		transport.DecodeWriteKeyRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext, transport.HTTPToContext(tracer, "writekey", logger)), // add JWT in the context
	)

	readkeyHandler := httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "readkey_endpoint")(endpoint.MakeReadKeyEndpoint(svc)),
		transport.DecodeReadKeyRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext, transport.HTTPToContext(tracer, "readkey", logger)),
		// add JWT in the context

	)

	createkeyHandler := httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "createkey_endpoint")(endpoint.MakeCreateKeyEndpoint(svc)),
		transport.DecodeCreateKeyRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(transport.HTTPToContext(tracer, "createkey", logger)),
		//httptransport.ServerBefore(httptransport.PopulateRequestContext), // add JWT in the context
	)

	exportkeyHandler := httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "exportkey_endpoint")(endpoint.MakeExportKeyEndpoint(svc)),
		transport.DecodeExportKeyRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(transport.HTTPToContext(tracer, "exportkey", logger)),
		//httptransport.ServerBefore(httptransport.PopulateRequestContext), // add JWT in the context
	)

	encryptHandler := httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "encrypt_endpoint")(endpoint.MakeEncryptEndpoint(svc)),
		transport.DecodeEncryptRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(transport.HTTPToContext(tracer, "encrypt", logger)),
		//httptransport.ServerBefore(httptransport.PopulateRequestContext), // add JWT in the context
	)

	decryptHandler := httptransport.NewServer(
		endpoint.MakeEndpointTracingMiddleware(tracer, "decrypt_endpoint")(endpoint.MakeDecryptEndpoint(svc)),
		transport.DecodeDecryptRequest,
		transport.EncodeResponse,
		httptransport.ServerBefore(transport.HTTPToContext(tracer, "decrypt", logger)),
		//httptransport.ServerBefore(httptransport.PopulateRequestContext), // add JWT in the context
	)

	var router = mux.NewRouter()
	router.PathPrefix("/key").Methods("GET").Handler(readkeyHandler)
	router.PathPrefix("/key").Methods("POST").Handler(writekeyHandler)
	router.PathPrefix("/createkey").Handler(createkeyHandler)
	router.PathPrefix("/exportkey").Handler(exportkeyHandler)
	router.PathPrefix("/encrypt").Handler(encryptHandler)
	router.PathPrefix("/decrypt").Handler(decryptHandler)

	http.ListenAndServe(httpAddr, router)

}

func config(logger log.Logger) map[string]interface{} {

	logger.Log("msg", "Loading configuration & command args")
	var configFile = fmt.Sprintf("conf/%s/vault_bridge.yaml", Environment)

	/*
		Component default
	*/
	viper.SetDefault("config-file", configFile)
	viper.SetDefault("component-name", "vault-bridge")
	viper.SetDefault("component-http-address", "127.0.0.1:8080")

	/*
		Vault default
	*/
	viper.SetDefault("vault-token", "6aeec159-592b-a807-ac54-3060f204116c") // we need to provide a Vault token with enough rights to run all the operations
	viper.SetDefault("vault-url", "http://127.0.0.1:8200")

	/*
		Influx DB client default
	*/
	viper.SetDefault("influx-url", "http://localhost:8086")
	viper.SetDefault("influx-username", "admin")
	viper.SetDefault("influx-password", "admin")
	viper.SetDefault("influx-database", "metrics")
	viper.SetDefault("influx-precision", "ms")
	viper.SetDefault("influx-retention-policy", "")
	viper.SetDefault("influx-write-consistency", "")
	viper.SetDefault("influx-counter-table-name", "vault_counter")
	viper.SetDefault("influx-histogram-table-name", "vault_histogram")

	/*
		Sentry client default
	*/
	viper.SetDefault("sentry-dsn", "https://1b65cbc7295f451e8b8e2711cd5ec3af:8ca839ccff4a49fb8b90e8ad7fa5f3b2@sentry.io/253288")

	/*
		Jaeger default
	*/
	viper.SetDefault("jaeger-service-name", "vault_bridge")
	viper.SetDefault("jaeger-sampler-type", "const")
	viper.SetDefault("jaeger-sampler-param", 1)
	viper.SetDefault("jaeger-sampler-url", "http://127.0.0.1:5775/")
	viper.SetDefault("jaeger-reporter-logspan", false)
	viper.SetDefault("jaeger-reporter-flushinterval-ms", 1000)

	/*
		First level of overhide
	*/
	pflag.String("config-file", viper.GetString("config-file"), "The configuration file path can be relative or absolute.")
	viper.BindPFlag("config-file", pflag.Lookup("config-file"))
	pflag.Parse()

	/*
		Load & log Config
	*/
	viper.SetConfigFile(viper.GetString("config-file"))
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log("msg", err)
	}

	var config = viper.AllSettings()
	for k, v := range config {
		logger.Log(k, v)
	}

	return config
}

//These are just some examples on how the requests to the bridge can be done
//JWT are generated with https://jwt.io/#debugger

//read a key
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" localhost:8080/key/tenants/rolex/f1/key1

//write a key
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"key":"abc"}' localhost:8080/key/tenants/rolex/f1/key1

//create a key
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"type": "aes256-gcm96", "derived": false, "exportable": true}}' localhost:8080/createkey/key1

//export a key
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" localhost:8080/exportkey/encryption-key/key100/1

//encrypt
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"plaintext": "YWJjZA==", "key_version": 1}}'  localhost:8080/encrypt/key100

//decrypt
//curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"ciphertext": "vault:v1:dKG1C2bFFdLMJQkau6v3lDmhHfLtTMBB9cd0XBy6Id8="}}'  localhost:8080/decrypt/key100
