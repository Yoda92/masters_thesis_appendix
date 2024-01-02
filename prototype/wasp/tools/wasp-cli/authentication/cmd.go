package authentication

import (
	"github.com/spf13/cobra"

	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

func initAuthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth <command>",
		Short: "Authentication tools",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Check(cmd.Help())
		},
	}
}

func Init(rootCmd *cobra.Command) {
	authCmd := initAuthCmd()
	loginCmd := initLoginCmd()
	infoCmd := initInfoCmd()

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(loginCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(infoCmd)

	loginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	loginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
}
