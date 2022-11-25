package cmd

import (
	"log"

	"github.com/9d4/semaphore/server"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "semaphore",
	Short: "Start semaphore server.",
	Long:  "Semaphore is blablabla..........",
	Run: boot(func(cmd *cobra.Command, args []string, passData *bootData) {
		log.Fatal(server.Start())
	}),
}
