package src

import (
	"fmt"
	"strings"

	"io"
	"net/http"
)

/*

Pulling an Image Manifest:
		The image manifest can be fetched with the following url:
			GET /v2/<name>/manifests/<reference>

Pulling a Layer:
		Layers are stored in the blob portion of the registry, keyed by digest. Pulling a layer is carried out by a standard http request. The URL is as follows:
			GET /v2/<name>/blobs/<digest> ---> digest is of the form sha256

*/

//--------------------------------------------------------------Authentication---------------------------------------------------------------------------

// check if you re authorized to authenticate.
func Check_endpoint(name string) map[string]string {
	full_path := fmt.Sprintf("https://registry-1.docker.io/v2/%s/tags/list", name)
	resp, err := http.Get(full_path) //Check that the endpoint implements Docker Registry API V2.
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//mimic the -i flag result. --> outputs a map with keys Bearer realm, scope and service.
	//for name, values := range resp.Header {
	m := make(map[string]string)
	for _, value := range resp.Header["Www-Authenticate"] { //p.Header["Www-Authenticate"] is an array of strings.
		pairs := strings.Split(value, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
	}
	//}
	return m
	//body, _ := io.ReadAll(resp.Body)
	//return string(body)
}

//If not, ask for a token.

func Request_token(attrs map[string]string) string {
	//The blueprint for requesting a token(supposi que 3and www headers.) --> The_val_of_bearer_realm + "?" +"service=" + The_val_of_service + "&" + "scope=" + val_of_repository
	full_path := fmt.Sprintf("%s?service=%s&scope=%s", attrs["Bearer realm"], attrs["service"], attrs["scope"])
	//i need to delete the "'s
	full_path = strings.ReplaceAll(full_path, `"`, "")
	resp, err := http.Get(full_path) //Check that the endpoint implements Docker Registry API V2.
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	full_token := string(body)
	//i need to filter the access token from the response.
	before, _, _ := strings.Cut(full_token, `","access_token`)
	_, after, _ := strings.Cut(before, `{"token":"`)

	/*
		print(full_token)
		print(before)
		print("\n------------------------------------------------------------------------\n")
		fmt.Print(ba3d)
		print("\n------------------------------------------------------------------------\n")
	*/

	return after
}

//i need to go back to https://distribution.github.io/distribution/spec/api/
