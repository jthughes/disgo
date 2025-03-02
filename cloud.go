package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azcache "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/joho/godotenv"
)

// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity@v1.8.2

// this example shows file storage but any form of byte storage would work
func retrieveRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}
	b, err := os.ReadFile("./entra.record.json")
	if err == nil {
		err = json.Unmarshal(b, &record)
	}
	return record, err
}

func storeRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	if err == nil {
		err = os.WriteFile("./entra.record.json", b, 0600)
	}
	return err
}

func SetupOneDrive(tokenOptions policy.TokenRequestOptions) (*azidentity.InteractiveBrowserCredential, error) {
	godotenv.Load()

	record, err := retrieveRecord()
	if err != nil {
		//
		fmt.Println("unable to retrieve record")
		// return
	}
	c, err := azcache.New(nil)
	if err != nil {
		// Persistent Cache impossible
		fmt.Println("persistent cache impossible")
		return nil, err
	}

	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		ClientID:             os.Getenv("DRIVE_CLIENT_ID"),
		TenantID:             os.Getenv("DRIVE_TENANT_ID"),
		Cache:                c,
	})
	if err != nil {
		// handle errorstore
		fmt.Println("unable to get credential")
		return nil, err
	}
	fmt.Println("credential acquired")

	if record == (azidentity.AuthenticationRecord{}) {
		// No stored record; call Authenticate to acquire one.
		record, err = cred.Authenticate(context.TODO(), &tokenOptions)
		if err != nil {
			fmt.Println("unable to authenticate credential")
			return nil, err
		}
		fmt.Println("credential authenticated")
		err = storeRecord(record)
		if err != nil {
			fmt.Println("unable to store record")
			return nil, err
		}
		fmt.Println("record stored")
	}
	return cred, nil
}

func GetFile(cred *azidentity.InteractiveBrowserCredential, tokenOptions policy.TokenRequestOptions, path string) (*http.Response, error) {
	accessToken, err := cred.GetToken(context.TODO(), tokenOptions)
	if err != nil {
		fmt.Printf("unable to get access token: %v", err)
		return nil, err
	}

	itemId := "F12027F22382A4D!505343"
	endpoint := fmt.Sprintf("drive/items/%s/content", itemId)
	baseurl := "https://graph.microsoft.com/v1.0/me"
	url := fmt.Sprintf("%s/%s", baseurl, endpoint)

	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v\n", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken.Token))

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Response: %d\n", response.StatusCode)

	fmt.Printf("Content Length: %d\n", response.ContentLength)
	return response, nil
}

func GetFolder(cred *azidentity.InteractiveBrowserCredential, tokenOptions policy.TokenRequestOptions, path string) {
	accessToken, err := cred.GetToken(context.TODO(), tokenOptions)
	if err != nil {
		fmt.Printf("unable to get access token: %v", err)
		return
	}

	filePath := "/Music/Video Games/Darren Korb/Songs of Supergiant Games/"
	endpoint := fmt.Sprintf("drive/root:%s:/children", filePath)
	baseurl := "https://graph.microsoft.com/v1.0/me"
	url := fmt.Sprintf("%s/%s", baseurl, endpoint)

	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %v\n", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken.Token))

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Response: %d\n", response.StatusCode)

	fmt.Printf("Content Length: %d\n", response.ContentLength)

	// fmt.Printf("Response:\n")
	// for k, v := range response.Header {
	// 	fmt.Printf("%v: %v\n", k, v)
	// }

	respDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RESPONSE:\n%s", string(respDump))
}
