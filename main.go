package main

import (
	"context"
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azcache "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/joho/godotenv"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

type drive struct {
	client_id     string
	client_secret string
	tenant_id     string
	object_id     string
	redirect_url  string
}

type config struct {
	drive drive
}

func main() {
	godotenv.Load()
	cfg := config{
		drive: drive{
			client_id:     os.Getenv("DRIVE_CLIENT_ID"),
			client_secret: os.Getenv("DRIVE_CLIENT_SECRET"),
			tenant_id:     os.Getenv("DRIVE_TENANT_ID"),
			object_id:     os.Getenv("DRIVE_OBJECT_ID"),
			redirect_url:  os.Getenv("DRIVE_REDIRECT_URL"),
		},
	}

	record, err := retrieveRecord()
	if err != nil {
		//
	}
	c, err := azcache.New(nil)
	if err != nil {
		// Persistent Cache impossible
	}

	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		Cache:                c,
	})
	// cred, _ := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
	// 	TenantID:    cfg.drive.tenant_id,
	// 	ClientID:    cfg.drive.client_id,
	// 	RedirectURL: cfg.drive.redirect_url,

	// 	Cache:       c,
	// })

	client, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"User.Read", "Files.Read"})

	result, err := client.Me().Drive().Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting the drive: %v\n", err)
		printOdataError(err)
		return
	}
	fmt.Printf("Found Drive : %v\n", *result.GetId())
	driveId := *result.GetId()

	root, _ := client.Drives().ByDriveId(driveId).Root().Get(context.Background(), nil)
	rootId := *root.GetId()

	folders, _ := client.Drives().ByDriveId(driveId).Items().ByDriveItemId(rootId).Children().Get(context.Background(), nil)

	writer :=
		folders.Serialize(s)
	// for item := range result.GetRoot().GetChildren() {
	// 	fmt.Printf("item: %v\n", item)
	// }

	// fmt.Println(client.Me().Drives())
	// children.
	// children.GetId()
	// https://learn.microsoft.com/en-us/graph/api/driveitem-get?view=graph-rest-1.0&tabs=go
	// https://learn.microsoft.com/en-us/graph/api/resources/onedrive?view=graph-rest-1.0

}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
