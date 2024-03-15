package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	tenant := "mytenantid"
	clientID := "myrandomid"
	clientSecret := "myrandomsecret"
	container := "test-data"
	adEndpoint := "https://127.0.0.1:8080"
	azuriteEndpoint := "https://127.0.0.1:10000/devstoreaccount1"

	cred, err := azidentity.NewClientSecretCredential(tenant, clientID, clientSecret, &azidentity.ClientSecretCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloud.Configuration{
				ActiveDirectoryAuthorityHost: adEndpoint,
			},
		},
		DisableInstanceDiscovery: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	client, err := azblob.NewClient(azuriteEndpoint, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateContainer(context.Background(), container, nil)
	if err != nil {
		// This is not the proper way of handling such thing - for demonstration purposes only
		if strings.Contains(err.Error(), "ContainerAlreadyExists") {
			fmt.Println("container already exists -- continuing")
		} else {
			log.Fatal(err)
		}
	}

	res, err := client.UploadStream(context.Background(),
		container,
		"test/path/to/upload/to/",
		bytes.NewReader([]byte("Hello World")),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
