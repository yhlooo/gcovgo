package gcovgo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	gcovraw "github.com/yhlooo/gcovgo/pkg/gcov/raw"
)

// newDumpCommand 创建 dump 子命令
func newDumpCommand() *cobra.Command {
	outputFormat := ""

	cmd := &cobra.Command{
		Use:   "dump PATH",
		Short: "Print coverage file contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read file %q error: %w", args[0], err)
			}

			raw := gcovraw.Raw{}
			if err := raw.UnmarshalBinary(content); err != nil {
				return fmt.Errorf("unmarshal gcov raw error: %w", err)
			}

			switch outputFormat {
			case "json":
				outputContent, err := json.MarshalIndent(raw, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal gcov raw to json error: %w", err)
				}
				fmt.Println(string(outputContent))
			default:
				// TODO: ...
			}

			return nil
		},
	}

	// 绑定选项到命令行参数
	cmd.Flags().StringVarP(&outputFormat, "format", "f", outputFormat, "Output format. One of 'yaml' or 'json'.")

	return cmd
}
