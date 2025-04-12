package repo

import (
	"mime/multipart"
	"strconv"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	VkCloudHotboxEndpoint = "https://hb.ru-msk.vkcloud-storage.ru"
	defaultRegion         = "ru-msk"
)

type awsRepository struct {
	sess *session.Session
	cfg  config.AWSConfig
}

func NewAWSRepository(sess *session.Session, cfg config.AWSConfig) awsRepository {
	return awsRepository{
		sess: sess,
		cfg:  cfg,
	}
}

func (ar *awsRepository) Update(id int, file *multipart.File) error {

	svc := s3.New(ar.sess, aws.NewConfig().WithEndpoint(VkCloudHotboxEndpoint).WithRegion(defaultRegion))

	key := "/1/" + strconv.Itoa(id)

	if _, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(ar.cfg.BucketName),
		Key:    aws.String(key),
		Body:   *file,
		ACL:    aws.String("public-read"),
	}); err != nil {
		return err
	}

	return nil
}
