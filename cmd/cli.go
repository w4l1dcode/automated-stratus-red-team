package main

import (
	"cloud-threat-emulation/config"
	"cloud-threat-emulation/pkg"
	"cloud-threat-emulation/pkg/stratus"
	"context"
	"flag"
	"fmt"
	_ "github.com/datadog/stratus-red-team/v2/pkg/stratus/loader" // Note: This import is needed
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"sync"
)

const (
	dbPathAWSTactics = "./cache/tactics_aws.db"
	dbPathK8sTactics = "./cache/tactics_k8s.db"
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
	err = os.Setenv("AWS_REGION", conf.AWS.Region)
	if err != nil {
		return
	}

	// Initialize the AWS client using Client
	awsClient, err := pkg.New(context.Background(), logger, conf.AWS.Region)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize AWS client")
	}

	// --

	collectErrors := make(chan error)
	collectWG := &sync.WaitGroup{}

	collectWG.Add(1)
	go func() {
		logger.Info("Executing AWS tests...")

		defer collectWG.Done()

		dbAWSTactics, awsErr := stratus.InitDB(dbPathAWSTactics)
		if awsErr != nil {
			collectErrors <- fmt.Errorf("failed to intialize database: %v", awsErr)
			return
		}

		awsErr = awsClient.DetonateTTPs(dbAWSTactics, "AWS")
		if awsErr != nil {
			collectErrors <- fmt.Errorf("failed to Execute stratus red team: %v", awsErr)
			return
		}

	}()

	// K8s Tests
	collectWG.Add(1)
	go func() {
		logger.Info("Executing K8s tests...")

		defer collectWG.Done()

		// Authenticate with Kubernetes
		downloadKubeConfig := exec.Command("sh", "-c", fmt.Sprintf("aws eks update-kubeconfig --region %s --name %s", conf.Kubernetes.Region, conf.Kubernetes.ClusterName))
		downloadKubeConfig.Stdout = os.Stdout
		downloadKubeConfig.Stderr = os.Stderr

		if err := downloadKubeConfig.Run(); err != nil {
			logger.WithError(err).Fatal("failed to update kube config")
		} else {
			logger.Info("Kube config updated successfully")
		}

		dbK8sTactics, k8sErr := stratus.InitDB(dbPathK8sTactics)
		if k8sErr != nil {
			collectErrors <- fmt.Errorf("failed to initialize K8s database: %v", k8sErr)
			return
		}

		k8sErr = awsClient.DetonateTTPs(dbK8sTactics, "kubernetes")
		if k8sErr != nil {
			collectErrors <- fmt.Errorf("failed to execute K8s stratus red team: %v", k8sErr)
			return
		}
	}()

	// --

	collectDone := make(chan struct{})
	go func() {
		collectWG.Wait()
		close(collectDone)
	}()

	logger.Info("Waiting for test executions to finish")
	select {
	case err := <-collectErrors:
		logger.WithError(err).Fatal("Failed to execute red team tests")
	case <-collectDone:
		logger.Info("Finished executing tests")
	}
}
