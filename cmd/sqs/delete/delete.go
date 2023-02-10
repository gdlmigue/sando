package delete

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

type DeleteMsg struct {
	ReceiptHandle *string
}

// Generated by https://quicktype.io
type SQSMessage struct {
	Records []Record `json:"Records"`
}

type Record struct {
	EventVersion      string            `json:"eventVersion"`
	EventSource       string            `json:"eventSource"`
	AwsRegion         string            `json:"awsRegion"`
	EventTime         string            `json:"eventTime"`
	EventName         string            `json:"eventName"`
	UserIdentity      ErIdentity        `json:"userIdentity"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements  ResponseElements  `json:"responseElements"`
	S3                S3                `json:"s3"`
}

type RequestParameters struct {
	SourceIPAddress string `json:"sourceIPAddress"`
}

type ResponseElements struct {
	XAmzRequestID string `json:"x-amz-request-id"`
	XAmzID2       string `json:"x-amz-id-2"`
}

type S3 struct {
	S3SchemaVersion string `json:"s3SchemaVersion"`
	ConfigurationID string `json:"configurationId"`
	Bucket          Bucket `json:"bucket"`
	Object          Object `json:"object"`
}

type Bucket struct {
	Name          string     `json:"name"`
	OwnerIdentity ErIdentity `json:"ownerIdentity"`
	Arn           string     `json:"arn"`
}

type ErIdentity struct {
	PrincipalID string `json:"principalId"`
}

type Object struct {
	Key       string `json:"key"`
	Size      int64  `json:"size"`
	ETag      string `json:"eTag"`
	VersionID string `json:"versionId"`
	Sequencer string `json:"sequencer"`
}

func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete",
		Run:   deleteAll,
	}
}

func getQueueURL() (*string, error) {
	svc := sqs.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	res, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("bc-sawmill-cdn-data-dlq"),
	})
	if err != nil {
		return nil, err
	}
	return res.QueueUrl, nil
}

func getMessages() error {
	queueURL, err := getQueueURL()
	if err != nil {
		return err
	}
	svc := sqs.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	for {

		msgRes, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            queueURL,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(10),
			VisibilityTimeout:   aws.Int64(5),
		})

		if err != nil {
			return err
		}

		// delData := []DeleteMsg{}
		fmt.Println(len(msgRes.Messages))

		for i := 0; i < len(msgRes.Messages); i++ {
			data := &SQSMessage{}
			body := msgRes.Messages[i].Body
			receiptHandle := msgRes.Messages[i].ReceiptHandle
			err = json.Unmarshal([]byte(*body), &data)
			if err != nil {
				return err
			}

			if data.Records[0].S3.ConfigurationID == "Ooyala Fastly" {
				fmt.Println("jeje")
				fmt.Println(msgRes)
				deleteMessage(*queueURL, receiptHandle)
			} else {
				fmt.Println("not jeje")
			}
		}
		continue
	}
}

func deleteMessage(queueURL string, messageHandle *string) error {
	svc := sqs.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	fmt.Println("Deleting ", *messageHandle)
	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: messageHandle,
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteAll(cmd *cobra.Command, _ []string) {
	getMessages()
}
