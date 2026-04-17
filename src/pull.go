package src

import (
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

//--------------------------------------------------------------Authentication---------------------------------------------------------------------------

// check if you re authorized to authenticate.
func Check_endpoint(name string) {
	full_path := fmt.Sprintf("https://registry-1.docker.io/v2/%s/tags/list", name)
	resp, err := http.Get(full_path) //Check that the endpoint implements Docker Registry API V2.
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//mimic the -i flag result. --> outputs a map with keys
	//for name, values := range resp.Header {
	for _, value := range resp.Header["Www-Authenticate"] { //p.Header["Www-Authenticate"] is an array of strings.
		m := make(map[string]string)
		pairs := strings.Split(value, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
		fmt.Println(m)
	}
	//}
	fmt.Println()
	//body, _ := io.ReadAll(resp.Body)
	//return string(body)
}

//If not, ask for a token.

func Request_token() {

}

//i need to go back to https://distribution.github.io/distribution/spec/api/
