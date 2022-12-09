package cmd

import (
	"fmt"
	"github.com/9d4/semaphore/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(keyCmd)
	keyCmd.AddCommand(keyGenerateCmd)
}

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Key utilities",
	RunE:  func(cmd *cobra.Command, args []string) error { return cmd.Help() },
}

var keyGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate 256 bit key",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.GenerateKey())
	},
}
