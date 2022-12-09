package cmd

import (
	"github.com/9d4/semaphore/store"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbSeedCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database utilities",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run database seeder",
	Run: boot(func(cmd *cobra.Command, args []string, passData *bootData) {
		// seed?
		jww.INFO.Print("Seeding database...")
		store.Seed(passData.store)
		jww.INFO.Print("done.")
	}),
}
