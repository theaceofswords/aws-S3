package main

import(
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"code.qburst.com/navaneeth.k/aws-s3/config"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type repo struct {
	svc  *s3.S3
	uploader *s3manager.Uploader
	downloader *s3manager.Downloader
}

type intrf interface {
	listBuckets()
	create()
	upload()
	download()
}

var bucketName string 

func main(){
	svc,uploader,downloader := config.Connect()
	repo := CreateRepository(svc,uploader,downloader)
	bucketName = "fr-odh-us-east-1-dev"
//	repo.create()
    repo.listBuckets()
//	repo.upload()
//	repo.download()
	

}
func CreateRepository(svci *s3.S3, upldr *s3manager.Uploader, dwnl *s3manager.Downloader) intrf {
	return &repo{
		svc:  svci,
		uploader: upldr,
		downloader: dwnl,
	}
}
func (r *repo) create(){
	
	var err error
	
	result, err := r.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		
		// CreateBucketConfiguration: &s3.CreateBucketConfiguration{
		// 	LocationConstraint: aws.String("us-east-1"),
		//},

	})
	if err != nil {
		fmt.Printf("Unable to create bucket %q, %v\n", bucketName, err)
	}
	
	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)
	
	err = r.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Printf("Error occurred while waiting for bucket to be created, %v\n", bucketName)
	}
	
	fmt.Printf("Bucket %q successfully created\n", bucketName)
	fmt.Println(result)
	
	
	listResult, err := r.svc.ListBuckets(nil)
    if err != nil {
        fmt.Printf("Unable to list buckets, %v\n", err)
    }

    fmt.Println("Buckets:")
	

    for _, b := range listResult.Buckets {
        fmt.Printf("* %s created on %s\n",
        aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
    }
	
	
}

func (r *repo) listBuckets(){

	result, err := r.svc.ListBuckets(nil)
    if err != nil {
        fmt.Printf("Unable to list buckets, %v\n", err)
    }

    fmt.Println("Buckets:")

    for _, b := range result.Buckets {
        fmt.Printf("* %s created on %s\n",
        aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
    }

}

func (r *repo) upload(){
	
    //filename := "fileA"

    // file, err := os.Open("./fileA.json")
	// if err != nil {
	// 	fmt.Printf("Unable to open file %q, %v", err)
	// }
	
	// defer file.Close()

	input := &s3.PutObjectInput{
		Body:    aws.ReadSeekCloser(strings.NewReader("./fileA.zip")),
		Bucket:  aws.String(bucketName),
		Key:     aws.String("fileA.zip"),
	}
	
	result, err := r.svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		
	}
	
	fmt.Println(result)

	// _, err = r.uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(bucket),
	// 	Key: aws.String(filename),
	// 	Body: file,
	// })
	// if err != nil {
	// 	// Print the error and exit.
	// 	fmt.Printf("Unable to upload %q to %q, %v\n", filename, bucket, err)
	// }else{
	// 	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	// }
	
	
	
}

func (r *repo) download(){
	
    //item := "fileA"
	//file, err := os.Create(item)
	
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String("bucketname1"),
		Key:    aws.String("fileA.zip"),
	}

	metaData, err := r.svc.HeadObject(headInput)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
        		default:
            		fmt.Println(aerr.Error())
     		}
    	} else {
        fmt.Println(err.Error())
   	 	}	
    
	}

	fmt.Println(*metaData.ContentLength)


	input := &s3.GetObjectInput{
	    Bucket: aws.String("bucketname1"),
	    Key:    aws.String("fileA.zip"),
		//Range: aws.String("0,1000"),
	}
	result, err := r.svc.GetObject(input)
	if err != nil {
	    if aerr, ok := err.(awserr.Error); ok {
	        switch aerr.Code() {
	        case s3.ErrCodeNoSuchKey:
	            fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
	        case s3.ErrCodeInvalidObjectState:
	            fmt.Println(s3.ErrCodeInvalidObjectState, aerr.Error())
	        default:
	            fmt.Println(aerr.Error())
	        }
	    } else {
	        // Print the error, cast err to awserr.Error to get the Code and
	        // Message from an error.
	        fmt.Println(err.Error())
	    }
	    
	}

	// numBytes, err := r.downloader.Download(file,
	// 	&s3.GetObjectInput{
	// 		Bucket: aws.String(bucketName),
	// 		Key:    aws.String(item),
	// 	})
	// if err != nil {
	// 	fmt.Printf("Unable to download item %q, %v\n", item, err)
	// }else{
	// 	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	// }
	
	fmt.Println(result.ContentLength)


}