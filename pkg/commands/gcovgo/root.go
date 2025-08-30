package gcovgo

import "github.com/spf13/cobra"

// NewCommand 创建根命令
func NewCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: "GCC code coverage tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newDumpCommand(),
		newVersionCommand(),
	)

	return cmd
}
