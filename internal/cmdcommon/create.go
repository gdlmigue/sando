package cmdcommon

import "github.com/spf13/cobra"

// SetCreateFlags sets flags supported by create command.
func SetCreateFlags(cmd *cobra.Command, prefix string) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("username", "u", "", prefix+"test")
	cmd.Flags().StringP("password", "p", "", prefix+"test")
}

func SetCreateReportFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("startDate", "s", "", "")
	cmd.Flags().StringP("endDate", "e", "", "")
}

func SetCreateBatchFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().StringP("queue", "q", "", "")
}
