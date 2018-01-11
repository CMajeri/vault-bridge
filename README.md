[![Build Status](https://travis-ci.org/cloustrust/vault-bridge.svg?branch=master)](https://travis-ci.org/cloustrust/vault-bridge)
[![Coverage Status](https://coveralls.io/repos/github/cloustrust/vault-bridge/badge.svg?branch=master)](https://coveralls.io/github/cloustrust/vault-bridge?branch=master)
 [![Go Report Card](https://goreportcard.com/badge/github.com/cloustrust/vault-bridge)](https://goreportcard.com/report/github.com/cloustrust/vault-bridge)


# Vault bridge
Vault bridge is a microservice that serves as a bridge between an application and Vault. The application can send the following requests to the vault bridge: 
- **write a key** in Vault (in key/value secret backend)
- **read a key** (from the key/value secret backend)
- **create a key** (in the transit backend) and **export** it
- **encrypt** a plaintext and **decrypt** a ciphertext (with keys from the transit backend)


## Launch

```bash
./scripts/build.sh --env <value>
``` 

```bash
./bin/vaultBridge
``` 

When you launch the vault bridge, the parameters are read from ``` conf/<value>/vault_bridge.yaml ```.

For loading a different configuratiuon file launch the sevice with the following command: 

```bash
./bin/vaultBridge --config-file "path/to/the/file.yaml".
``` 

## Configuration

In order to run the service, you need to configure: 
- **http address**: the vault bridge listens to the *http address* given in the configuration file.     

- **Vault**: the Vault service needs to be set up. In order to communicate with Vault, the vault bridge needs to have a *token* and the *address* of Vault. 

- **InfluxDB**: the metrics will be send to an Influx time series DB. The service needs the *url*, *username*, *password*, *db name* of the DB and the *table names* where the metrics information is stored. Also it needs the parameters of the *influxBatchPointsConfig*.  

- **Sentry**: the errors and crashes are sent to a Sentry error tracking system. We need to provide the *sentry DSN*.

- **Jaeger**: the tracing will be sent to Jaeger. The vault bridge needs to provide the *Jaeger configuration* parameters.


## Usage 

The vault bridge listens for http requests at different prefix paths, depending on the requests: ```/key```, ```/createkey```, ```/exportkey```, ```/encrypt```, ```/decrypt```. 

Each http request must contain a JWT, the path and, if needed, data. 

The parameters used in the http requests follow the Vault API. 

## Usage - examples

For these examples we work localhost and the vault bridge listens at port 8080.

*read a key*

```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" localhost:8080/key/tenants/rolex/f1/key1
```
With this request, the microservice reads the key that is stored at ```tenants/rolex/f1/key1```.

*write a key*
```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"key":"abc"}' localhost:8080/key/tenants/rolex/f1/key1
```
With this request, the microservice writes the key/value ```{"key":"abc"}``` at ```tenants/rolex/f1/key1```.


*create a key*
```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"type": "aes256-gcm96", "derived": false, "exportable": true}}' localhost:8080/createkey/key1
```
With this request, the microservice asks Vault to create an ```exportable``` ```aes256-gcm96``` key named ```key1``` in the transit backend.

*export a key*
```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" localhost:8080/exportkey/encryption-key/key100/1
```
With this request, the microservice fetches from Vault the key named ```key100``` from the transit backend.

*encrypt*
```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"plaintext": "YWJjZA==", "key_version": 1}}'  localhost:8080/encrypt/key100
```

With this request, the microservice asks Vault to encrypt the plaintext (base64 encoded) ```YWJjZA==``` with the key named ```key100``` from the transit backend.

*decrypt*
```bash
curl -H "authorization:  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnQiOiJyb2xleCIsImZpbCI6ImYxIn0.qt8lC6BOTVVx1RiEShpdgF43v1TAvTPGVdtL2rdixcc" -d'{"params": {"ciphertext": "vault:v1:dKG1C2bFFdLMJQkau6v3lDmhHfLtTMBB9cd0XBy6Id8="}}'  localhost:8080/decrypt/key100
```    
With this request, the microservice asks Vault to decrypt the ciphertext (base64 encoded) ```vault:v1:dKG1C2bFFdLMJQkau6v3lDmhHfLtTMBB9cd0XBy6Id8=``` with the key named ```key100``` from the transit backend.














