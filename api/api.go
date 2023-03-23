package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

const (
	GET = "GET"
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
	}

	return nil, errors.New(fmt.Sprintf("method must be 'GET', got '%s'", method))
}
