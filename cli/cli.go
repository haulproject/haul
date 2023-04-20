package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	OutputStyleJSON       = "json"
	OutputStyleJSONPretty = "json_pretty"
	OutputStyleTabby      = "tabby" // Default
)

type Client struct {
	OutputStyle  string   // e.g. "json"
	TabbyHeaders []string // e.g. ["id", "name"]
}

func New() *Client {
	return &Client{OutputStyle: OutputStyleTabby}
}

func (c *Client) Output(message []byte) error {
	switch c.OutputStyle {
	case OutputStyleJSON:
		fmt.Println(string(message))

	case OutputStyleJSONPretty:
		var bytes_buffer bytes.Buffer

		if err := json.Indent(&bytes_buffer, message, "", "  "); err != nil {
			return err
		}
		_, err := fmt.Println(string(bytes_buffer.Bytes()))

		return err

	/*
		case OutputStyleTabby:
			var result interface{}

			json.Unmarshal(message, &result)

			fmt.Println("tabby")

			//t := tabby.New()

			switch result.(type) {
			case map[string]interface{}:
				fmt.Println("is a map[string]interface{}")
			case []interface{}:
				fmt.Println("is a []interface{}")
				result_slice, ok := result.([]interface{})
				if !ok {
					return fmt.Errorf("Cannot type result as []interface{}")
				}

				var fields []interface{}

				fmt.Println(fields)

				for _, object := range result_slice[0] {
					fmt.Println(object)
				}

			case []map[string]interface{}:
				fmt.Println("is a []map[string]interface{}")
			default:
				return errors.New("Unknown type for tabby message")

			}
	*/
	default:
		return errors.New("Invalid OutputStyle")
	}

	return nil
}
