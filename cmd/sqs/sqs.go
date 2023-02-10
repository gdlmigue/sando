package sqs

import (
	"sando/cmd/sqs/delete"

	"github.com/spf13/cobra"
)

func NewCmdSQS() *cobra.Command {
	cmd := cobra.Command{
		Use:  "sqs",
		RunE: queues,
	}

	dc := delete.NewCmdDelete()

	cmd.AddCommand(dc)

	return &cmd
}

func queues(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
