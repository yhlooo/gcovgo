package gcovgo

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/yhlooo/gcovgo/pkg/version"
)

const versionTemplate = `Version:   {{ .Version }}
GitCommit: {{ .GitCommit }}
GoVersion: {{ .GoVersion }}
Arch:      {{ .Arch }}
OS:        {{ .OS }}
`

// newVersionCommand 创建 version 子命令
func newVersionCommand() *cobra.Command {
	outputFormat := ""

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.GetVersion()

			switch outputFormat {
			case "yaml":
				raw, err := yaml.Marshal(info)
				if err != nil {
					return err
				}
				fmt.Println(string(raw))
			case "json":
				raw, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(raw))
			default:
				tpl, err := template.New("version").Parse(versionTemplate)
				if err != nil {
					return err
				}
				return tpl.Execute(os.Stdout, info)
			}

			return nil
		},
	}

	// 将选项绑定到命令行
	cmd.Flags().StringVarP(&outputFormat, "format", "f", outputFormat, "Output format. One of 'yaml' or 'json'.")

	return cmd
}
