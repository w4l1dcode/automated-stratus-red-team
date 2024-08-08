package main

import (
	"automated-stratus-red-team/config"
	"automated-stratus-red-team/pkg/stratus"
	"flag"
	_ "github.com/datadog/stratus-red-team/v2/pkg/stratus/loader" // Note: This import is needed
	"github.com/sirupsen/logrus"
	"os"
)

const (
	dbPath = "./cache/tactics.db"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.WithField("level", logrusLevel.String()).Info("set log level")
	logger.SetLevel(logrusLevel)

	// --

	dir := "cache"
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(dir, 0755)
		if err != nil {
			logrus.Fatalf("Error creating directory: %v\n", err)
			return
		}
		logrus.Info("Directory .cache created")
	}

	// --

	err = os.Setenv("AWS_ACCESS_KEY_ID", conf.AWS.AwsAccessKeyId)
	if err != nil {
		return
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", conf.AWS.AwsSecretAccessKey)
	if err != nil {
		return
	}
	err = os.Setenv("AWS_SESSION_TOKEN", conf.AWS.SessionToken)
	if err != nil {
		return
	}
	err = os.Setenv("AWS_REGION", "eu-west-1")
	if err != nil {
		return
	}

	// --

	db := stratus.InitDB(dbPath)

	logger.Info("Starting AWS tests...")

	awsErr := stratus.DetonateTTPs(db, "AWS")
	if awsErr != nil {
		logrus.Fatalf("failed to Execute stratus red team: %v", awsErr)
		return
	}
}
