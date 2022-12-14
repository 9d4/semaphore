package cmd

import (
	"fmt"
	"github.com/9d4/semaphore/oauth"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"text/tabwriter"
)

func init() {
	rootCmd.AddCommand(oAuthCmd)
	oAuthCmd.AddCommand(oAuthAddCmd)
}

var oAuthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "OAuth utilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var oAuthAddCmd = &cobra.Command{
	Use:   "add [app-name]",
	Short: "Add new client app",
	Args:  cobra.ExactArgs(1),
	Run: boot(func(cmd *cobra.Command, args []string, passData *bootData) {
		oauthAppStore := oauth.NewStore(passData.db, passData.rdb)
		app, err := oauthAppStore.CreateAuto(args[0])
		if err != nil {
			jww.FATAL.Fatal(err)
		}
		tw := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
		fmt.Println("Created!")
		fmt.Fprintf(tw, "Name\t :%s\n", app.Name)
		fmt.Fprintf(tw, "ClientID\t :%s\n", app.ClientID)
		tw.Flush()
	}),
}
