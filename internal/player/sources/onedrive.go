package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/jthughes/disgo/internal/player"
	"github.com/jthughes/disgo/internal/repl"
)

type OneDriveSource struct {
	tokenOptions policy.TokenRequestOptions
	userAccount  public.Account
	accessToken  string
}

func (s OneDriveSource) String() string {
	return "onedrive"
}

type TokenCache struct {
	file string
}

func (t *TokenCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	data, err := os.ReadFile(t.file)
	if err != nil {
		log.Println(err)
	}
	return cache.Unmarshal(data)
}

func (t *TokenCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) error {
	data, err := cache.Marshal()
	if err != nil {
		log.Println(err)
	}
	return os.WriteFile(t.file, data, 0600)
}

// Creates a OneDriveSource (implements Source) by authenticating with OneDrive
func InitOneDriveSource(cfg repl.Config) (OneDriveSource, error) {

	cache := TokenCache{
		file: fmt.Sprintf("%s/.onedrive.cache", cfg.ConfigPath),
	}
	client, err := public.New(
		"01cf73e8-6601-4df1-8282-03ccd68a7075", // Client ID
		public.WithCache(&cache),
		public.WithAuthority("https://login.microsoftonline.com/common"), // Tennant ID
	)
	if err != nil {
		return OneDriveSource{}, err
	}

	accounts, err := client.Accounts(context.TODO())
	if err != nil {
		return OneDriveSource{}, err
	}
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}
	var result public.AuthResult
	if len(accounts) > 0 {
		// There may be more accounts; here we assume the first one is wanted.
		// AcquireTokenSilent returns a non-nil error when it can't provide a token.
		result, err = client.AcquireTokenSilent(context.TODO(), tokenOptions.Scopes, public.WithSilentAccount(accounts[0]))
	}
	if err != nil || len(accounts) == 0 {
		// cache miss, authenticate a user with another AcquireToken* method
		result, err = client.AcquireTokenInteractive(context.TODO(), tokenOptions.Scopes)
		if err != nil {
			// TODO: handle error
		}
	}

	s := OneDriveSource{}
	s.tokenOptions = tokenOptions
	s.accessToken = result.AccessToken
	s.userAccount = result.Account
	return s, nil
}

// Makes an authenticated request against a provided endpoint
func (s OneDriveSource) Request(request string, endpoint string) (*http.Response, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(request, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	return resp, err
}

func (s OneDriveSource) ScanFolder(path string) ([]player.Track, error) {
	// Build request
	baseurl := "https://graph.microsoft.com/v1.0/me"
	var endpoint string
	if path != "/" && path != "" {
		endpoint = fmt.Sprintf("drive/root:%s:/children", path)
	} else {
		endpoint = "drive/root/children"
	}
	url := fmt.Sprintf("%s/%s", baseurl, endpoint)
	// fmt.Println("Getting: " + url)

	resp, err := s.Request("GET", url)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Response: %d\n", resp.StatusCode)

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
		// result.Print()
	}

	tracks := []player.Track{}
	for _, item := range result.Value {
		if item.Audio != (player.AudioMetadata{}) {
			// Add Track
			tracks = append(tracks, player.Track{
				FileName: item.Name,
				Metadata: item.Audio,
				MimeType: item.File.MimeType,
				Data:     player.NewFile(item.ID, s.String(), s),
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

// Given a OneDrive fileId string, gets the io.ReadCloser data for that file
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
		Audio player.AudioMetadata `json:"audio,omitempty"`
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
