package cmdcommon

import "github.com/spf13/cobra"

// SetCreateFlags sets flags supported by create command.
func SetCreateFlags(cmd *cobra.Command, prefix string) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("username", "u", "", prefix+"test")
	cmd.Flags().StringP("password", "p", "", prefix+"test")

}
