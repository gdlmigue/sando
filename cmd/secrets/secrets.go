package secrets

import (
	"sando/cmd/secrets/create"

	"github.com/spf13/cobra"
)

const helpText = `Create a new secret`

func NewCmdSecrets() *cobra.Command {
	cmd := cobra.Command{
		Use:     "secrets",
		Short:   "Create secrets",
		Long:    helpText,
		Aliases: []string{"secret"},
		RunE:    secret,
	}

	cc := create.NewCmdCreate()

	cmd.AddCommand(cc)

	create.SetFlags(cc)

	return &cmd
}

func secret(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
