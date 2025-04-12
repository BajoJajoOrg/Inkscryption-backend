package aws

import (
	"log/slog"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcloud-storage.ru"
	defaultRegion         = "ru-msk"
)

func NewClient(cfg config.AWSConfig, logger *slog.Logger) (*session.Session, error) {
	// awsConfig := aws.Config{

	// }
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))

	if res, err := svc.ListBuckets(nil); err != nil {
		return nil, err
	} else {
		for _, b := range res.Buckets {
			logger.Info("* %s created on %s \n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	}
	return sess, nil
}
