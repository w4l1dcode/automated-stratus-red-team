package pkg

import (
	"cloud-threat-emulation/pkg/stratus"
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/sirupsen/logrus"
)

// Client represents a generic AWS client with a logger and an ECR client.
type Client struct {
	ctx       context.Context
	logger    *logrus.Logger
	ecrClient *ecr.Client
}

// New initializes a new Client with an ECR client configured.
func New(ctx context.Context, l *logrus.Logger, region string) (*Client, error) {
	client := Client{logger: l, ctx: ctx}

	// Load the AWS configuration, specifying the region
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	// Create an ECR client from the configuration
	client.ecrClient = ecr.NewFromConfig(cfg)
	return &client, nil
}

func (c *Client) DetonateTTPs(dbAWSTactics *sql.DB, platform string) error {
	return stratus.DetonateTTPs(dbAWSTactics, platform, c.ecrClient, c.logger)
}
