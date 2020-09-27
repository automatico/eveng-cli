package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"eveng-cli/src/api"
	"eveng-cli/src/utils"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Inspect your EVE-NG server",
	Long:  `eveng-cli server status`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		server := api.ServerConfig()
		cookie, err := utils.JSONCookieFileToStruct(".cookie.json")
		// client := utils.SetupHTTPClient(false)

		status := server.GetStatus(cookie)
		statusCode, err := api.HTTPReturnCodes(status)

		if err != nil && statusCode == 412 {
			cookie = server.Auth()
			utils.CookieToJSONFile(".cookie.json", cookie)
		}

		if args[0] == "status" {
			status := server.GetStatus(cookie)
			utils.PrintResponse(status)
		} else if args[0] == "config" {
			fmt.Println("Create credentials logic.")

		} else {
			fmt.Printf("Invalid argument: %s\n", args[0])
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
