package sqs

import (
	"sando/cmd/sqs/create"
	"sando/cmd/sqs/delete"

	"github.com/spf13/cobra"
)

const helpText = `Manage different SQS events`

func NewCmdSQS() *cobra.Command {
	cmd := cobra.Command{
		Use:   "sqs",
		Short: "Manage events",
		Long:  helpText,
		RunE:  queues,
	}

	dc := delete.NewCmdDelete()
	cc := create.NewCmdCreate()

	cmd.AddCommand(dc, cc)

	create.SetFlags(cc)

	return &cmd
}

func queues(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
