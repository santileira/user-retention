package script

import (
	"github.com/santileira/user-retention/domain/userretention/calculator"
	handler "github.com/santileira/user-retention/domain/userretention/handler"
	validator "github.com/santileira/user-retention/domain/userretention/validator"
	"github.com/santileira/user-retention/filereader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		"file with the application's user activity")

	cmd.Run = func(_ *cobra.Command, _ []string) {
		r.configureLog(options.LogLevel)
		r.calculateUserRetention(options)
	}
	return cmd
}

func (r *Runnable) calculateUserRetention(options *Options) {
	userRetentionValidator := validator.NewUserRetentionValidatorImpl()
	fileReader := filereader.NewFileReaderImpl()
	userRetentionCalculator := calculator.NewUserRetentionCalculatorImpl()
	userRetentionHandler := handler.NewUserRetentionHandler(userRetentionValidator, fileReader, userRetentionCalculator)

	result, err := userRetentionHandler.HandleRequest(options.FilePath)
	if err != nil {
		logrus.Errorf("Error calculating the user retention, err: %s", err.Error())
		return
	}

	logrus.Debugf("The result is: \n%s", result)
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
