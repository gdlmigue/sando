package root

import (
	"sando/cmd/reports"
	"sando/cmd/secrets"
	"sando/cmd/sqs"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sando-cli",
	Short: "Getting out of trouble on cdn stuff",
}

func NewCmdRoot() *cobra.Command {
	cmd := cobra.Command{
		Use:   "sando <command> <subcommand>",
		Short: "Interactive CDN CLI",
		Long:  "Interactive CDN command line to get out of trouble",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	addChildCommands(&cmd)
	return &cmd
}

func addChildCommands(cmd *cobra.Command) {
	cmd.AddCommand(secrets.NewCmdSecrets())
	cmd.AddCommand(sqs.NewCmdSQS())
	cmd.AddCommand(reports.NewCmdReports())
}
