package main

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func main() {
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("source.accessToken: %v\n", source.accessToken)

	tracks, err := source.ScanFolder("/Music/Video Games/Darren Korb/")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		track.Print()
	}

	if len(tracks) > 0 {
		err = tracks[len(tracks)-1].Play()
		if err != nil {
			fmt.Println(err)
		}
	}

	// fileData, _ := source.DownloadFile("F12027F22382A4D!505343")

	// play(fileData)
}
