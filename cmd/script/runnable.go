package script

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/santileira/user-retention/calculator"
)

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

func (r *Runnable) Cmd() *cobra.Command {
	options := &Options{}

	var cmd = &cobra.Command{
		Use:   "script",
		Short: "Script to calculate application's user retention",
		Long:  `Script to calculate application's user retention`,
	}

	cmd.Flags().StringVar(&options.LogLevel, "log-level", defaultLogLevel, "log leve to use")
	cmd.Flags().StringVar(&options.FilePath, "file-path", defaultFilePath,
		"file with the user activity for the application")

	cmd.Run = func(_ *cobra.Command, _ []string) {
		r.configureLog(options.LogLevel)
		r.calculateUserRetention(options)
	}
	return cmd
}

func (r *Runnable) calculateUserRetention(options *Options) {
	userRetentionCalculator := calculator.NewUserRetentionCalculator(options.FilePath, 0)
	result, err := userRetentionCalculator.Calculate()
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func (r *Runnable) configureLog(logLevel string) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warnf("Error passing the log level: %s", logLevel)
		return
	}

	logrus.Infof("Setting log level: %s", lvl.String())
	logrus.SetLevel(lvl)
}
