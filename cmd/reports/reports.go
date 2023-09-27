package reports

import (
	"sando/cmd/reports/create"

	"github.com/spf13/cobra"
)

const helpText = `Create a bandwidth report`

func NewCmdReports() *cobra.Command {
	cmd := cobra.Command{
		Use:     "reports",
		Short:   "Create a report",
		Long:    helpText,
		Aliases: []string{"report"},
		RunE:    report,
	}

	cc := create.NewCmdCreate()

	cmd.AddCommand(cc)

	create.SetFlags(cc)

	return &cmd
}

func report(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
