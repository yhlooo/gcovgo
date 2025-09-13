package gcovgo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yhlooo/gcovgo/pkg/gcov"
)

// NewCommand 创建根命令
func NewCommand(name string) *cobra.Command {
	verbosity := 0
	cpuProfile := ""
	outputFormat := "human-readable"
	outputFile := ""

	var cpuProfileOutput *os.File
	cmd := &cobra.Command{
		Use:          name + " {SOURCE|OBJ}...",
		Short:        "GCC code coverage tool",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 设置日志
			setLogger(cmd, verbosity)
			// 输出 CPU 性能数据
			if cpuProfile != "" {
				var err error
				cpuProfileOutput, err = os.OpenFile(cpuProfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					return fmt.Errorf("open cpu profile file %q error: %w", cpuProfile, err)
				}
				if err := pprof.StartCPUProfile(cpuProfileOutput); err != nil {
					return fmt.Errorf("start cpu profile error: %w", err)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logr.FromContextOrDiscard(ctx)

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

			resolvedNoteFiles := map[string]bool{}
			for _, fileName := range args {
				fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
				noteFileName := fileName + ".gcno"
				if resolvedNoteFiles[noteFileName] {
					continue
				}
				resolvedNoteFiles[noteFileName] = true

				dataFileName := fileName + ".gcda"
				if _, err := os.Stat(dataFileName); err != nil {
					if !os.IsNotExist(err) {
						logger.Error(err, fmt.Sprintf("get data file %q info error", dataFileName))
						continue
					}
					dataFileName = ""
				}

				ret, err := gcov.ResolveBinaryFile(noteFileName, dataFileName)
				if err != nil {
					logger.Error(err, fmt.Sprintf("resolve %q error", noteFileName))
					continue
				}
				ret.DataFile = fileName
				ret.GcovNoteFile = noteFileName
				ret.GcovDataFile = dataFileName

				var outputContent []byte
				switch outputFormat {
				case "text":
					outputContent = []byte(ret.IntermediateText(ctx))
				case "json":
					outputContent, err = json.MarshalIndent(ret, "", "  ")
					if err != nil {
						return fmt.Errorf("marshal result to json error: %w", err)
					}
				case "human-readable":
					outputContent = []byte(ret.HumanReadableText(ctx))
				default:
					return fmt.Errorf("unknown output format: %q", outputFormat)
				}
				if _, err = fmt.Fprintln(w, string(outputContent)); err != nil {
					return fmt.Errorf("write output error: %w", err)
				}
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cpuProfileOutput != nil {
				pprof.StopCPUProfile()
				if err := cpuProfileOutput.Close(); err != nil {
					return fmt.Errorf("close cpu profile file %q error: %w", cpuProfile, err)
				}
			}
			return nil
		},
	}

	// 绑定选项到命令行参数
	globalFlags := cmd.PersistentFlags()
	globalFlags.IntVarP(&verbosity, "verbose", "v", verbosity, "Number for the log level verbosity (0, 1, or 2)")
	globalFlags.StringVar(&cpuProfile, "cpu-profile", cpuProfile, "Write a CPU profile to the specified file")
	fs := cmd.Flags()
	fs.StringVarP(&outputFormat, "format", "f", outputFormat, `Output format, one of:
  human-readable : human readable format
  text           : intermediate text format
  json           : intermediate JSON format
`)
	fs.StringVarP(&outputFile, "output", "o", outputFile, "Write output to file instead of stdout")

	// 添加子命令
	cmd.AddCommand(
		newDumpCommand(),
		newVersionCommand(),
	)

	return cmd
}

// setLogger 设置命令日志，并返回 logr.Logger
func setLogger(cmd *cobra.Command, verbosity int) logr.Logger {
	// 设置日志级别
	logrusLogger := logrus.New()
	switch verbosity {
	case 1:
		logrusLogger.SetLevel(logrus.DebugLevel)
	case 2:
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}
	// 将 logger 注入上下文
	logger := logrusr.New(logrusLogger)
	cmd.SetContext(logr.NewContext(cmd.Context(), logger))

	return logger
}
