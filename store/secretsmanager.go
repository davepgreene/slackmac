package store

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

type SecretsManagerStore struct {
	svc *secretsmanager.SecretsManager
	secretId string
}

func (s *SecretsManagerStore) Get() string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.secretId),
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

func NewSecretsManagerStore(conf map[string]string) (Store, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	region, ok := conf["region"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the AWS secrets manager datastore", "region"))
	}
	cfg.Region = region

	svc := secretsmanager.New(cfg)

	secretId, ok := conf["id"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the AWS secrets manager datastore", "id"))
	}

	return &SecretsManagerStore{
		svc: svc,
		secretId: secretId,
	}, nil
}
