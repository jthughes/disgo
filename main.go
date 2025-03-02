package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func main() {
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}
	cred, err := SetupOneDrive(tokenOptions)
	if err != nil {
		return
	}

	GetFolder(cred, tokenOptions, "")
	// response, err := GetFile(cred, tokenOptions, "")
	// if err != nil {
	// 	return
	// }
	// play(response.Body)
}
