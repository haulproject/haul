/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// componentCreateCmd represents the componentCreate command
var componentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a component in the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Must specify at least 1 component to create")
		}
		if len(args) > 1 {
			log.Fatal("Can currently only insert 1 component at a time")
		}

		var components []types.Component

		for _, arg := range args {
			var component types.Component

			err := json.Unmarshal([]byte(arg), &component)
			if err != nil {
				log.Println(err)
				continue
			}

			components = append(components, component)

		}

		_, err := json.Marshal(components)
		if err != nil {
			log.Fatal(err)
		}

		for _, component := range components {
			if component.Name == "" {
				os.Stderr.WriteString("component.Name cannot be empty")
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		currentComponent, err := json.Marshal(components[0])

		endpoint := fmt.Sprintf("%s://%s:%d",
			viper.GetString("api.protocol"),
			viper.GetString("api.host"),
			viper.GetInt("api.port"),
		)
		request := fmt.Sprintf("%s%s", endpoint, "/v1/component")

		resp, err := http.Post(request, "application/json",
			bytes.NewBuffer(currentComponent))

		if err != nil {
			log.Fatal(err)
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		fmt.Println(res["message"])
	},
}

func init() {
	componentCmd.AddCommand(componentCreateCmd)
}
