package utils

import (
	"bytes"
	"io"
	"net/http"
)

func CallRESTAPI(url string, method string, body []byte) (respBody []byte, respHeader http.Header, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return responseBody, resp.Header, nil
}

func CallRESTAPIWithToken(url string, method string, body []byte, token string) (respBody []byte, respHeader int, err error) {
	if body == nil {
		body = []byte{}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return responseBody, resp.StatusCode, nil
}
