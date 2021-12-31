package config

import(
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func Connect() (*s3.S3, *s3manager.Uploader, *s3manager.Downloader ){
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("dfsfsd7w492rh", "ebfiwubiu42353f", ""),
		Endpoint:         aws.String("http://localhost:4566"),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Create s3 client
	// svc := s3.New(sess, &aws.Config{
	// 	Endpoint: aws.String("http://localhost:4566"),
	// 	Region: aws.String("us-west-2"),
	// })
	svc := s3.New(sess)

	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)
	
	fmt.Println("conn +1")
	//fmt.Printf("%T",uploader)
	return svc,uploader,downloader
}