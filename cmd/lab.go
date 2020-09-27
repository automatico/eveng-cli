package cmd

import (
	"crypto/tls"
	"eveng-cli/src/api"
	"eveng-cli/src/utils"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// labCmd represents the lab command
var labCmd = &cobra.Command{
	Use:   "lab",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		fmt.Println("lab called")
		server := api.ServerConfig()
		cookie, err := utils.JSONCookieFileToStruct(".cookie.json")
		status := server.GetStatus(cookie)
		statusCode, err := api.HTTPReturnCodes(status)
		if err != nil && statusCode == 412 {
			cookie = server.Auth()
			utils.CookieToJSONFile(".cookie.json", cookie)
		}

		utils.PrintResponse(status)
	},
}

func init() {
	rootCmd.AddCommand(labCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// labCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
