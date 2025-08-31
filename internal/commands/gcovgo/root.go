package gcovgo

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand 创建根命令
func NewCommand(name string) *cobra.Command {
	verbosity := 0
	cpuProfile := ""

	var cpuProfileOutput *os.File
	cmd := &cobra.Command{
		Use:          name,
		Short:        "GCC code coverage tool",
		SilenceUsage: true,
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
			return cmd.Help()
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
	fs := cmd.PersistentFlags()
	fs.IntVarP(&verbosity, "verbose", "v", verbosity, "Number for the log level verbosity (0, 1, or 2)")
	fs.StringVar(&cpuProfile, "cpu-profile", cpuProfile, "Write a CPU profile to the specified file")

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
