package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type MyEvent struct {
	ScreenName string `json:"screen_name"`
	PNGBase64  string `json:"pngbase64"`
}

type MyResponse struct {
	URI string `json:"uri:"`
	OK  bool   `json:"ok"`
}

func handler(request MyEvent) (MyResponse, error) {
	var BUCKET = os.Getenv("BUCKET")
	var KEY = fmt.Sprintf("/%v.png", request.ScreenName)
	fmt.Printf("loaded envvar\n")
	// extract image file from event
	decoded, err := b64.StdEncoding.DecodeString(request.PNGBase64)
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return MyResponse{URI: "", OK: false}, err
	}
	fmt.Printf("decoded\n")

	// save image into temporary file
	tmpFile, err := ioutil.TempFile("", "received*.png")
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return MyResponse{URI: "", OK: false}, err
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write(decoded)
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return MyResponse{URI: "", OK: false}, err
	}
	tmpFile.Sync()
	tmpFile.Close()
	fmt.Printf("wrote\n")

	// call primitive
	primitive := exec.Command("/primitive", "-n", "10", "-m", "1", "-i", tmpFile.Name(), "-o", "/tmp/result.png")
	err = primitive.Run()
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return MyResponse{URI: "", OK: false}, err
	}
	fmt.Printf("ran primitive\n")

	// load result image
	resultFile, err := ioutil.ReadFile("/tmp/result.png")
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return MyResponse{URI: "", OK: false}, err
	}
	fmt.Printf("loaded image\n")

	// upload file into S3
	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(endpoints.ApNortheast1RegionID),
	})
	_, errpo := svc.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(resultFile),
		Bucket: aws.String(BUCKET),
		Key:    aws.String(KEY),
		ACL:    aws.String("public-read"),
	})
	if errpo != nil {
		fmt.Printf("error occurred: %v\n", errpo)
		return MyResponse{URI: "", OK: false}, errpo
	}
	fmt.Printf("uploaded\n")
	// return image URI
	return MyResponse{URI: "", OK: true}, nil
}

func main() {
	lambda.Start(handler)
}
