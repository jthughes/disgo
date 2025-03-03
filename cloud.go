package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

type OneDriveSource struct {
	cred         *azidentity.InteractiveBrowserCredential
	tokenOptions policy.TokenRequestOptions
	accessToken  azcore.AccessToken
}

func NewOneDriveSource(tokenOptions policy.TokenRequestOptions) (OneDriveSource, error) {
	s := OneDriveSource{}
	s.tokenOptions = tokenOptions
	godotenv.Load()
	record, err := retrieveRecord()
	if err != nil {
		fmt.Println("unable to retrieve record")
	}
	c, err := azcache.New(nil)
	if err != nil {
		return OneDriveSource{}, fmt.Errorf("persistent cache impossible: %w", err)
	}

	s.cred, err = azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		ClientID:             os.Getenv("DRIVE_CLIENT_ID"),
		TenantID:             os.Getenv("DRIVE_TENANT_ID"),
		Cache:                c,
	})
	if err != nil {
		return OneDriveSource{}, fmt.Errorf("unable to get credential: %w", err)
	}

	if record == (azidentity.AuthenticationRecord{}) {
		fmt.Println("prompting user for authentication...")
		// No stored record; call Authenticate to acquire one.
		record, err = s.cred.Authenticate(context.TODO(), &tokenOptions)
		if err != nil {
			return OneDriveSource{}, fmt.Errorf("unable to authenticate credential: %w", err)
		}
		fmt.Println("credential authenticated")
		err = storeRecord(record)
		if err != nil {

			return OneDriveSource{}, fmt.Errorf("unable to store record: %w", err)
		}
		fmt.Println("record stored")
	}
	s.accessToken, err = s.cred.GetToken(context.TODO(), s.tokenOptions)
	if err != nil {
		return OneDriveSource{}, fmt.Errorf("unable to get access token: %w", err)
	}
	return s, nil
}

func (s OneDriveSource) Request(request string, endpoint string) (*http.Response, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(request, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.accessToken.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	return resp, err
}

func (s OneDriveSource) GetFolder(path string) {
	baseurl := "https://graph.microsoft.com/v1.0/me"
	var endpoint string
	if path != "/" && path != "" {
		endpoint = fmt.Sprintf("drive/root:%s:/children", path)
	} else {
		endpoint = "drive/root/children"
	}
	url := fmt.Sprintf("%s/%s", baseurl, endpoint)
	fmt.Println("Getting: " + url)

	resp, err := s.Request("GET", url)
	if err != nil {

	}

	fmt.Printf("Response: %d\n", resp.StatusCode)

	fmt.Printf("Content Length: %d\n", resp.ContentLength)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result OneDriveResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	result.Print()
}

func (s OneDriveSource) DownloadFile(id string) (io.ReadCloser, error) {
	if id == "" {
		// error
	}
	baseurl := "https://graph.microsoft.com/v1.0/me"
	endpoint := fmt.Sprintf("drive/items/%s/content", id)
	url := fmt.Sprintf("%s/%s", baseurl, endpoint)

	resp, err := s.Request("GET", url)
	if err != nil {
		// error
	}

	return resp.Body, nil
}

type OneDriveResponse struct {
	OdataCount int `json:"@odata.count"`
	Value      []struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		ID                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		Name                 string    `json:"name"`
		Size                 int       `json:"size"`
		ParentReference      struct {
			DriveID   string `json:"driveId"`
			DriveType string `json:"driveType"`
			ID        string `json:"id"`
			Name      string `json:"name"`
			Path      string `json:"path"`
		} `json:"parentReference"`
		Audio struct {
			Album             string `json:"album"`
			AlbumArtist       string `json:"albumArtist"`
			Artist            string `json:"artist"`
			Bitrate           int    `json:"bitrate"`
			Duration          int    `json:"duration"`
			Genre             string `json:"genre"`
			HasDrm            bool   `json:"hasDrm"`
			IsVariableBitrate bool   `json:"isVariableBitrate"`
			Title             string `json:"title"`
			Track             int    `json:"track"`
			Year              int    `json:"year"`
		} `json:"audio,omitempty"`
		File struct {
			MimeType string `json:"mimeType"`
			Hashes   struct {
				QuickXorHash string `json:"quickXorHash"`
				Sha1Hash     string `json:"sha1Hash"`
				Sha256Hash   string `json:"sha256Hash"`
			} `json:"hashes"`
		} `json:"file"`
		FileSystemInfo struct {
			CreatedDateTime      time.Time `json:"createdDateTime"`
			LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		} `json:"fileSystemInfo"`
		Folder struct {
			ChildCount int `json:"childCount"`
			View       struct {
				ViewType  string `json:"viewType"`
				SortBy    string `json:"sortBy"`
				SortOrder string `json:"sortOrder"`
			} `json:"view"`
		} `json:"folder"`
	} `json:"value"`
}

func (r OneDriveResponse) Print() {
	str, _ := json.MarshalIndent(r, "", "\t")
	fmt.Printf("%s\n", str)
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

	filePath := "/Music/Video Games/Darren Korb/" + "Songs of Supergiant Games/"
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

	// respDump, err := httputil.DumpResponse(response, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("RESPONSE:\n%s", string(respDump))

	// body, err := io.ReadAll(response.Body) // response body is []byte
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // fmt.Printf("%s\n", body)
	// var bufferOut bytes.Buffer
	// err = json.Indent(&bufferOut, body, "", "\t")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(bufferOut.String())

	// Print files
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// var result FolderResponse
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// s, _ := json.MarshalIndent(result, "", "\t")
	// // fmt.Printf("%+v\n", result)
	// fmt.Printf("%s\n", s)

	// Print all?
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result OneDriveResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, _ := json.MarshalIndent(result, "", "\t")
	// fmt.Printf("%+v\n", result)
	fmt.Printf("%s\n", s)
}
