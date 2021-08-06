package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	ScreenName string `json:"screen_name"`
	PNGBase64  string `json:"pngbase64"`
}

type MyResponse struct {
	URI string `json:"uri:"`
}

func handler(request MyEvent) (MyResponse, error) {
	// extract image file from event
	// save image into temporary file
	// call primitive
	// load result image
	// upload file into S3
	// return image URI
	return MyResponse{URI: ""}, nil
}

func main() {
	lambda.Start(handler)
}
