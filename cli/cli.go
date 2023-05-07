package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"codeberg.org/haulproject/haul/types"
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

// OutputObject prints indented or unindented json, or an ascii table (tabby)
// (depending on *Client#OutputStyle) to stdout
func (c *Client) OutputObject(tabby_printer types.TabbyPrinter) error {
	switch c.OutputStyle {
	case OutputStyleJSON:
		message, err := json.Marshal(tabby_printer)
		if err != nil {
			return err
		}

		fmt.Println(string(message))

	case OutputStyleJSONPretty:
		message, err := json.MarshalIndent(tabby_printer, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(message))
	case OutputStyleTabby:
		tabby_printer.TabbyPrint()
	default:
		return fmt.Errorf("Unknown output style %s", c.OutputStyle)
	}

	return nil
}

// Output prints indented or unindented (depending on *Client#OutputStyle)
// json from a []byte to stdout.
//
// Deprecated: Does not support tabby output, please use *Client#OutputObject
// instead for full features.
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
	default:
		return errors.New("Invalid OutputStyle")
	}

	return nil
}
