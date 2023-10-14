package connections

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Connection interface{}

type s3_connection struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
}

func NewS3Connection() S3Connection {
	return &s3_connection{
		AccessKey:  "AKIA3VMV3LWIQ6EL63WU",
		SecretKey:  "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",
		Region:     "ap-south-1",
		BucketName: "skromansupportbucket",
	}
}

func (s3_bucket *s3_connection) MakeNewSession() (*s3.S3, error) {
	creds := credentials.NewStaticCredentials(s3_bucket.AccessKey,
		s3_bucket.SecretKey, "",
	)

	_, err := creds.Get()

	if err != nil {
		return nil, err
	}
}
