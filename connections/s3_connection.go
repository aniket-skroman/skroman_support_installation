package connections

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Connection interface {
	MakeNewSession() (*session.Session, error)
	GetBucketName() string
	UploadDeviceImage(file_path string) (*s3.PutObjectOutput, string, error)
}

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

func (s3_bucket *s3_connection) MakeNewSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(s3_bucket.AccessKey,
		s3_bucket.SecretKey, "",
	)

	_, err := creds.Get()

	if err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion(s3_bucket.Region).WithCredentials(creds)
	sess, _ := session.NewSession(cfg)

	return sess, nil
}

func (s3_bucket *s3_connection) GetBucketName() string {
	return s3_bucket.BucketName
}

func (s3_bucket *s3_connection) UploadDeviceImage(file_path string) (*s3.PutObjectOutput, string, error) {
	sess, err := s3_bucket.MakeNewSession()

	if err != nil {
		return nil, "", err
	}

	svc := s3.New(sess)

	file, err := os.Open(file_path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := file.Name()
	params := &s3.PutObjectInput{
		Bucket:        aws.String(s3_bucket.GetBucketName()),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	result, err := svc.PutObject(params)
	return result, path, err
}
