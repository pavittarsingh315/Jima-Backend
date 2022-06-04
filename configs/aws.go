package configs

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	region               = "us-west-1"
	bucketName           = "nerajima"
	accessKey, secretKey = EnvAWSCredentials()
)

// "By default, the SDK detects AWS credentials set in your environment and uses them to sign requests to AWS. That way you don’t need to manage credentials in your applications."
// All we need to do is have "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY" environment variables set in our environment and SDK will do the rest.
// Read more here: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
var s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
	Region:      aws.String(region),
	Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
})))

// Generate an upload url that'll put a file in an S3 bucket in a specific directory.
func GenerateS3UploadUrl(directory string) (string, error) {
	fileName, genErr := generateRandFileName(64)
	if genErr != nil {
		return "", genErr
	}
	fileKey := fmt.Sprintf("%s/%s", directory, fileName)

	req, _ := s3Session.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	url, err := req.Presign(time.Minute * 1)

	if err != nil {
		return "", err
	}
	return url, nil
}

// Delete an S3 object located in the path of the S3 bucket.
func DeleteS3Object(filePath string) error {
	_, err := s3Session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}
	return nil
}

func generateRandFileName(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
