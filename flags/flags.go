package flags

import (
	"github.com/spf13/cobra"
)

var (
	DBUSER, DBPASS, SRVHOST, SRVPORT, DBHOST, DBPORT string
)

func Flags() {
	rootCmd := &cobra.Command{
		Use: "cmngb",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.Flags().StringVarP(&DBUSER, "username", "u", "", "set usernme")
	rootCmd.Flags().StringVarP(&DBPASS, "secret", "s", "", "set secret for database access")
	rootCmd.Flags().StringVarP(&SRVHOST, "host", "i", "localhost", "set server ip host")
	rootCmd.Flags().StringVarP(&SRVPORT, "port", "p", ":8080", "set server port")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("username")
	rootCmd.Execute()
}
