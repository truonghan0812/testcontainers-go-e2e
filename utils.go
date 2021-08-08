package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func JsonReader(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func ReadResp(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		panic(err)
	}

}
