package src

import (
	"io"
	"net/http"
)

func Check_endpoint() string {
	resp, err := http.Get("https://registry-1.docker.io/v2/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

//i need to go back to https://distribution.github.io/distribution/spec/api/
