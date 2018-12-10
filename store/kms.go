package store

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	log "github.com/sirupsen/logrus"
)

type KMSStore struct {
	ciphertext string
	cipherTextBlob []byte
	svc *kms.KMS
}

func (k *KMSStore) Get() string {
	var input *kms.DecryptInput

	input = &kms.DecryptInput{
		CiphertextBlob: k.cipherTextBlob,
	}

	req := k.svc.DecryptRequest(input)
	result, err := req.Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				log.Errorf(kms.ErrCodeNotFoundException, aerr.Error())
			case kms.ErrCodeDisabledException:
				log.Errorf(kms.ErrCodeDisabledException, aerr.Error())
			case kms.ErrCodeInvalidCiphertextException:
				log.Errorf(kms.ErrCodeInvalidCiphertextException, aerr.Error())
			case kms.ErrCodeKeyUnavailableException:
				log.Errorf(kms.ErrCodeKeyUnavailableException, aerr.Error())
			case kms.ErrCodeDependencyTimeoutException:
				log.Errorf(kms.ErrCodeDependencyTimeoutException, aerr.Error())
			case kms.ErrCodeInvalidGrantTokenException:
				log.Errorf(kms.ErrCodeInvalidGrantTokenException, aerr.Error())
			case kms.ErrCodeInternalException:
				log.Errorf(kms.ErrCodeInternalException, aerr.Error())
			case kms.ErrCodeInvalidStateException:
				log.Errorf(kms.ErrCodeInvalidStateException, aerr.Error())
			default:
				log.Errorf(aerr.Error())
			}
		} else {
			log.Errorf(err.Error())
		}
		return ""
	}

	return string(result.Plaintext)
}

func NewKMSStore(conf map[string]string) (Store, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	region, ok := conf["region"]
	if ok {
		cfg.Region = region
	}

	svc := kms.New(cfg)

	ciphertext, ok := conf["ciphertext"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the AWS KMS datastore", "data"))
	}

	cipherTextBlob, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	return &KMSStore{
		cipherTextBlob: cipherTextBlob,
		ciphertext: ciphertext,
		svc: svc,
	}, nil
}
