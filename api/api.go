// Package api implements functions to make http requests, with or without data, of multiple possible methods.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

func Call(method, route string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s://%s:%d",
		viper.GetString("api.protocol"),
		viper.GetString("api.host"),
		viper.GetInt("api.port"),
	)
	request := fmt.Sprintf("%s%s", endpoint, route)

	switch method {
	case GET:
		response, err := http.Get(request)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	case DELETE:
		// Create client
		client := &http.Client{}

		// Create request
		req, err := http.NewRequest("DELETE", request, nil)
		if err != nil {
			return nil, err
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	}
	return nil, errors.New(fmt.Sprintf("method must be 'GET' or 'DELETE', got '%s'", method))
}

func CallWithData(method, route string, data []byte) (string, error) {
	endpoint := fmt.Sprintf("%s://%s:%d",
		viper.GetString("api.protocol"),
		viper.GetString("api.host"),
		viper.GetInt("api.port"),
	)
	request := fmt.Sprintf("%s%s", endpoint, route)
	switch method {
	case POST:

		resp, err := http.Post(request, "application/json",
			bytes.NewBuffer(data))

		if err != nil {
			return "", err
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		return fmt.Sprintf("%s\n", res["message"]), nil
	case PUT:
		// initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, request, bytes.NewBuffer(data))
		if err != nil {
			return "", err
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)
		return fmt.Sprintf("%s\n", res["message"]), nil
	}

	return "", errors.New(fmt.Sprintf("method must be 'POST' or 'PUT', got '%s'", method))
}
