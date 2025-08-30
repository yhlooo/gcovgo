package gcovgo

import "github.com/spf13/cobra"

// newDumpCommand 创建 dump 子命令
func newDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Print coverage file contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	return cmd
}
