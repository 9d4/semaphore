package cmd

import (
	"context"
	"fmt"
	"github.com/9d4/semaphore/oauth2/models"
	"github.com/9d4/semaphore/oauth2/store"
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
	Use:   "add [app-id] [app-domain] [app-secret]",
	Short: "Add new client app",
	Args:  cobra.ExactArgs(3),
	Run: boot(func(cmd *cobra.Command, args []string, passData *bootData) {
		clientStore := store.NewClientStoreRedis(passData.rdb)
		err := clientStore.Set(args[0], &models.Client{
			ID:     args[0],
			Secret: args[2],
			Domain: args[1],
		})
		if err != nil {
			jww.FATAL.Fatal(err)
			return
		}

		cli, err := clientStore.GetByID(context.Background(), args[0])
		if err != nil {
			jww.FATAL.Fatal(err)
			return
		}

		tw := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
		fmt.Println("Created!")
		fmt.Fprintf(tw, "ClientID\t :%s\n", cli.GetID())
		fmt.Fprintf(tw, "Secret\t :%s\n", cli.GetSecret())
		fmt.Fprintf(tw, "Domain\t :%s\n", cli.GetDomain())
		tw.Flush()
	}),
}
