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
	outputFormat := "json"
	outputFile := ""

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

			// 打开输出文件
			w := os.Stdout
			if outputFile != "" {
				var err error
				w, err = os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					return fmt.Errorf("open output file %q error: %w", outputFile, err)
				}
				defer func() { _ = w.Close() }()
			}

			var outputContent []byte
			switch outputFormat {
			case "json":
				outputContent, err = json.MarshalIndent(raw, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal gcov raw to json error: %w", err)
				}
			case "gcov":
				// TODO
				fallthrough
			default:
				return fmt.Errorf("unknown output format: %q", outputFormat)
			}
			if _, err = fmt.Fprintln(w, string(outputContent)); err != nil {
				return fmt.Errorf("write output error: %w", err)
			}

			return nil
		},
	}

	// 绑定选项到命令行参数
	fs := cmd.Flags()
	fs.StringVarP(&outputFormat, "format", "f", outputFormat, "Output format, one of 'gcov' or 'json'")
	fs.StringVarP(&outputFile, "output", "o", outputFile, "Write output to file instead of stdout")

	return cmd
}
