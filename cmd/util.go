package cmd

import "github.com/spf13/cobra"

type bootData struct{}

type (
	cobraFunc func(cmd *cobra.Command, args []string)
	bootFunc  func(cmd *cobra.Command, args []string, passData *bootData)
)

func boot(fn bootFunc) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		// todo
		// connect db and something else here
		data := &bootData{}

		fn(cmd, args, data)
	}
}
