package store

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

type secretsManagerStore struct {
	svc      *secretsmanager.SecretsManager
	secretID string
}

// Get retrieves data from AWS secrets manager
func (s *secretsManagerStore) Get() string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.secretID),
	}

	req := s.svc.GetSecretValueRequest(input)
	result, err := req.Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				log.Errorf(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				log.Errorf(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				log.Errorf(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				log.Errorf(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				log.Errorf(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				log.Errorf(aerr.Error())
			}
		} else {
			log.Errorf(err.Error())
		}
		return ""
	}

	return *result.SecretString
}

func newSecretsManagerStore(conf map[string]string) (Store, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	region, ok := conf["region"]
	if !ok {
		return nil, fmt.Errorf("%s is required for the AWS secrets manager datastore", "region")
	}
	cfg.Region = region

	svc := secretsmanager.New(cfg)

	secretID, ok := conf["id"]
	if !ok {
		return nil, fmt.Errorf("%s is required for the AWS secrets manager datastore", "id")
	}

	return &secretsManagerStore{
		svc:      svc,
		secretID: secretID,
	}, nil
}
