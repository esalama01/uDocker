package src

import (
	"encoding/json"
	"fmt"
	"strings"

	//"io"
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

// --------------------------------------------------------------Authentication---------------------------------------------------------------------------
type TokenResponse struct {
	Token string `json:"token"`
}

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
	authHeader := resp.Header.Get("Www-Authenticate")

	trimmedHeader := strings.TrimPrefix(authHeader, "Bearer ")
	pairs := strings.Split(trimmedHeader, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.Trim(strings.TrimSpace(kv[1]), "\"")
			m[key] = val
		}
	}
	return m
	//body, _ := io.ReadAll(resp.Body)
	//return string(body)
}

//If not, ask for a token.

func Request_token(attrs map[string]string) string {
	//The blueprint for requesting a token(supposi que 3and www headers.) --> The_val_of_bearer_realm + "?" +"service=" + The_val_of_service + "&" + "scope=" + val_of_repository
	full_path := fmt.Sprintf("%s?service=%s&scope=%s", attrs["realm"], attrs["service"], attrs["scope"])
	//i need to delete the "'s
	full_path = strings.ReplaceAll(full_path, `"`, "")
	resp, err := http.Get(full_path) //Check that the endpoint implements Docker Registry API V2.
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var tokenRes TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		panic(err)
	}

	return tokenRes.Token
}

func Authenticate(name string) int {
	if !strings.Contains(name, "/") {
		name = "library/" + name
	}
	m := Check_endpoint(name)
	token := Request_token(m)
	full_path := fmt.Sprintf("https://registry-1.docker.io/v2/%s/tags/list", name)
	req, err := http.NewRequest("GET", full_path, nil)
	if err != nil {
		panic(err)
	}

	// 2. Set the headers (equivalent to -H)
	req.Header.Set("Authorization", "Bearer "+token)
	// 3. Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp.StatusCode
}

func Manifest(name, ref string) int { //a function for pulling an image manifest.
	if !strings.Contains(name, "/") {
		name = "library/" + name
	}
	m := Check_endpoint(name)
	token := Request_token(m)
	full_path := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", name, ref)
	req, err := http.NewRequest("GET", full_path, nil)
	if err != nil {
		panic(err)
	}

	// 2. Set the headers (equivalent to -H)
	req.Header.Set("Authorization", "Bearer "+token)
	// 3. Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp.StatusCode
}

//i need to go back to https://distribution.github.io/distribution/spec/api/
