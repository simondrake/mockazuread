# Introduction

The purpose of this service is to mock the endpoints needed to run local Service Principal authentication for Azure, specifically so it can be used with [Azurite](https://learn.microsoft.com/en-us/azure/storage/common/storage-use-azurite) **without** needing use a real Service Principal.

If [Azurite issue 2373](https://github.com/Azure/Azurite/issues/2373) is ever implemented, it will hopefully make this redundant.

# Prerequisites

* [mkcert](https://github.com/FiloSottile/mkcert)

# Notes and Caveats

* This service has only been tested with the Go SDK

# Set-up

* Because the Azure SDK expects a valid TLS certificate, generate a new certificate with `mkcert` with the required SANs. **Note:** substitute the SANS below with the ones relevant to your set-up.

```bash
mkcert 127.0.0.1 localhost azurite
```

This will create a certificate and key in the current directory (e.g.  `127.0.0.1+2.pem` and `127.0.0.1+2-key.pem`)

* To ensure the certificates are trusted, run `mkcert -install`
* Set the following environment variables:
  * `MOCKAZURE_CONNECTION_CERTDIRECTORY` to the directory where the certificates were created
  * `MOCKAZURE_CONNECTION_CERTNAME` to the name of the certificate file
  * `MOCKAZURE_CONNECTION_KEYNAME` to the name of the key file
* Run `go run .`

```bash
MOCKAZURE_CONNECTION_CERTDIRECTORY=/tmp/mockazurecerts MOCKAZURE_CONNECTION_CERTNAME=cert.pem MOCKAZURE_CONNECTION_KEYNAME=key.pem go run .
```

* Run the Azurite docker-compose file, ensuring that the volume is set to the same directory mentioned above (**Note:** The provided docker compose file assumes a certificate filename of `cert.pem` and a key filename of `key.pem`)

```yaml
    volumes:
      - /tmp/mockazurecerts:/workspace
```

* You should then be able to run the example by running `cd examples/simple; go run main.go`

# AZ CLI

If you need to use the [az CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli), for example to query the data that has been stored, you can use the `REQUESTS_CA_BUNDLE` environment variable pointed to the mkcert Root CA file. For example:

```bash
$ REQUESTS_CA_BUNDLE=~/.local/share/mkcert/rootCA.pem az storage blob download --container-name test-data --name "path/to/data" --connection-string "DefaultEndpointsProtocol=https;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=https://minikube:30700/devstoreaccount1;"
```
