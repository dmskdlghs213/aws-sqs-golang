package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var svc *sqs.SQS
var queueURL string

func init() {
	cred := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")

	region := os.Getenv("AWS_REGION")
	conf := &aws.Config{
		Credentials: cred,
		Region:      &region,
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		log.Error(err)
	}

	svc = sqs.New(sess)
	queueURL = os.Getenv("TEST_SQS_URL")
}

func main() {

	e := echo.New()

	e.POST("/send", Send)
	e.GET("/receive", Receive)

	e.Logger.Fatal(e.Start(":80")) // listen and serve on :8080
}

func Receive(c echo.Context) error {

	var max int64 = 1
	var wait int64 = 10
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: &max,
		WaitTimeSeconds:     &wait,
	}

	resp, err := svc.ReceiveMessage(receiveParams)
	if err != nil {
		return err
	}

	var receiptHandle string
	for i := range resp.Messages {
		receiptHandle = *resp.Messages[i].ReceiptHandle
	}

	params := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	}

	_, err = svc.DeleteMessage(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func Send(c echo.Context) error {

	body := "1 :Hello !!"
	sendParams := &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &body,
	}

	_, err := svc.SendMessage(sendParams)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "success !!")
}
