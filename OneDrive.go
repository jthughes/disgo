package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azcache "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/joho/godotenv"
)

type OneDriveSource struct {
	cred         *azidentity.InteractiveBrowserCredential
	tokenOptions policy.TokenRequestOptions
	accessToken  azcore.AccessToken
}

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

func (s OneDriveSource) ScanFolder(path string) ([]Track, error) {
	// Build request
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
		return nil, err
	}

	fmt.Printf("Response: %d\n", resp.StatusCode)

	// Process request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result OneDriveResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "FLAC") {
		result.Print()
	}

	tracks := []Track{}
	for _, item := range result.Value {
		if item.Audio != (AudioMetadata{}) {
			// Add Track
			tracks = append(tracks, Track{
				FileName: item.Name,
				Metadata: item.Audio,
				MimeType: item.File.MimeType,
				Data: OneDriveFileSource{
					id:     item.ID,
					source: &s,
				},
			})
		} else if item.Folder != (OneDriveFolder{}) {
			// Recursively add tracks from nested folder
			nestedTracks, err := s.ScanFolder(fmt.Sprintf("%s%s/", path, item.Name))
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Adding %d tracks from %s\n", len(nestedTracks), item.Name)
			tracks = append(tracks, nestedTracks...)
		}

	}
	return tracks, nil
}

func (s OneDriveSource) GetFolder(path string) ([]Track, error) {
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
		return nil, err
	}

	fmt.Printf("Response: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result OneDriveResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	result.Print()
	// result.Print()
	tracks := []Track{}
	for _, item := range result.Value {
		if item.Audio == (AudioMetadata{}) {
			continue
		}
		tracks = append(tracks, Track{
			FileName: item.Name,
			Metadata: item.Audio,
			MimeType: item.File.MimeType,
			Data: OneDriveFileSource{
				id:     item.ID,
				source: &s,
			},
		})
	}
	return tracks, nil
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

type OneDriveFolder struct {
	ChildCount int `json:"childCount"`
	View       struct {
		ViewType  string `json:"viewType"`
		SortBy    string `json:"sortBy"`
		SortOrder string `json:"sortOrder"`
	} `json:"view"`
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
		Audio AudioMetadata `json:"audio,omitempty"`
		File  struct {
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
		Folder OneDriveFolder `json:"folder"`
	} `json:"value"`
}

func (r OneDriveResponse) Print() {
	str, _ := json.MarshalIndent(r, "", "\t")
	fmt.Printf("%s\n", str)
}
