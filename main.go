package main

import (
    "bytes"
    "log"
    "net/http"
    "fmt"
    "os"
    "flag"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

const (
    DEFAULT_S3_REGION = "us-east-2"
)

func main() {

    // Environmental Setup
    accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if len(accessKey) == 0 {
        fmt.Println("No environment variable found for: AWS_ACCESS_KEY_ID. Quiting.")
        os.Exit(1)
    }

    secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if len(secretKey) == 0 {
        fmt.Println("No environment variable found for: AWS_SECRET_ACCESS_KEY. Quiting.")
        os.Exit(1)
    }

    creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
    
    _, err := creds.Get() 
    if err != nil { 
        panic(err)
    }

    altAwsRegion := os.Getenv("AWS_REGION")
	if len(altAwsRegion) == 0 {
        altAwsRegion = DEFAULT_S3_REGION
    }

    altAwsBucket := os.Getenv("AWS_BUCKET")
	if len(altAwsBucket) == 0 {
        altAwsBucket = ""
    }

    // Flag Setup
    bucketPtr := flag.String("bucket", "", "Aws Bucket Name")
	destPtr := flag.String("destination", "", "Path to destination (must start with '/')")
	srcPtr := flag.String("source", "", "Path to file to upload")
	regionPtr := flag.String("region", "", "Region to use")

    flag.Parse()

    if len(*bucketPtr) == 0 {
        bucketPtr = &altAwsBucket
    }

    if len(*regionPtr) == 0 {
        fmt.Println("No environment variable found for: AWS_REGION. Using default: " + DEFAULT_S3_REGION)
        regionPtr = &altAwsRegion
    }

    // Confirm required Values
    if len(*bucketPtr) == 0 {
        fmt.Println("No bucket value provided. Please use the -bucket option. Quiting.")
        os.Exit(1)
    }

    if len(*destPtr) == 0 {
        fmt.Println("No destination value provided. Please use the -destination option. Quiting.")
        os.Exit(1)
    }

    if len(*srcPtr) == 0 {
        fmt.Println("No source value provided. Please use the -source option. Quiting.")
        os.Exit(1)
    }
    
    fmt.Println("Authenticating to AWS...")

    session, err := session.NewSession(&aws.Config{
        Region: aws.String(*regionPtr),
        Credentials: creds,
    })


    if err != nil {
        log.Fatal(err)
    }

    err = AddFileToS3(session, *srcPtr, *destPtr, *bucketPtr)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Upload complete")

}

func AddFileToS3(s *session.Session, src string, dest string, bucket string) error {
    file, err := os.Open(src)
    if err != nil {
        return err
    }
    defer file.Close()

    fileInfo, _ := file.Stat()
    var size int64 = fileInfo.Size()
    buffer := make([]byte, size)
    file.Read(buffer)

    fmt.Println("Uploading " + src + " to " + dest + " on bucket " + bucket)

    _, err = s3.New(s).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(bucket),
        Key:                  aws.String(dest),
        // ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(buffer),
        // ContentLength:        aws.Int64(size),
        ContentType:          aws.String(http.DetectContentType(buffer)),
        // ContentDisposition:   aws.String("attachment"),
        // ServerSideEncryption: aws.String("AES256"),
    })
    return err
}
