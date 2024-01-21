package flags

import (
	"github.com/spf13/cobra"
)

var NAME, PASS string

func Load() {
	rootCmd := cobra.Command{
		Use:   "cmngb",
		Short: "Start server",
	}
	rootCmd.Flags().StringVarP(&NAME, "username", "u", "", "set db username")
	rootCmd.Flags().StringVarP(&PASS, "password", "p", "", "set db password")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.Execute()
}
