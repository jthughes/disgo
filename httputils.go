package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func printHttpResponseHeaders(resp *http.Response) {
	fmt.Printf("Response Headers:\n")
	for k, v := range resp.Header {
		fmt.Printf("%v: %v\n", k, v)
	}
}

func printHttpResponse(resp *http.Response) {
	respDump, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("Response:\n%s", string(respDump))
}

func printHttpResponseBody(resp *http.Response, formatted bool) {
	body, err := io.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		fmt.Println(err)
		return
	}
	if formatted {
		var bufferOut bytes.Buffer
		err = json.Indent(&bufferOut, body, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(bufferOut.String())
	} else {
		fmt.Printf("%s\n", body)
	}
}
