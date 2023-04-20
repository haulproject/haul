/*
 */
package cmd

import (
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"github.com/spf13/cobra"
)

// kitListCmd represents the kitList command
var kitListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Prints values of all kits",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		kits_bytes, err := api.Call(http.MethodGet, "/v1/kit")
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		/* This allows printing a list of kits in a tabby Table. Not reusable enough.
		if client.OutputStyle == cli.OutputStyleTabby {
			// tabby
			t := tabby.New()

			t.AddHeader("id", "name", "status", "tags")

			var kits []types.KitWithID

			if err := json.Unmarshal(kits_bytes, &kits); err != nil {
				log.Fatal(err)
			}

			for _, kit := range kits {
				tags, err := json.Marshal(kit.Tags)
				if err != nil {
					log.Fatal(err)
				}

				t.AddLine(kit.ID, kit.Name, kit.Status, string(tags))
			}

			t.Print()
			os.Exit(0)
		}
		*/

		err = client.Output(kits_bytes)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	kitCmd.AddCommand(kitListCmd)
}
