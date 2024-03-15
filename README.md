# Introduction

The purpose of this service is to mock the endpoints needed to run local Service Principal authentication for Azure, specifically so it can be used with [Azurite](https://learn.microsoft.com/en-us/azure/storage/common/storage-use-azurite) **without** needing use a real Service Principal.

If [Azurite issue 2373](https://github.com/Azure/Azurite/issues/2373) is ever implemented, it will hopefully make this redundant.

# Notes and Caveats

* This service has only been tested with the Go SDK
* The certificates were generated using [mkcert](https://github.com/FiloSottile/mkcert). You may need to install and run `mkcert -install` for it to work.

# Example

* Run the provided `docker-compose.yaml` file - `docker-compose up -d`
* Then `cd examples/simple; go run main.go`
