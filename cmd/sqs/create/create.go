package create

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sando/internal/cmdcommon"
	"sando/internal/cmdutil"
	"sando/internal/query"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

// S3Event struct
type S3Event struct {
	Records []RecordsData `json:"Records"`
}

// RecordsData struct
type RecordsData struct {
	EventVersion string `json:"eventVersion"`
	EventSource  string `json:"eventSource"`
	AwsRegion    string `json:"awsRegion"`
	EventTime    string `json:"eventTime"`
	EventName    string `json:"eventName"`
	S3           S3Data `json:"s3"`
}

// S3Data struct
type S3Data struct {
	SchemaVersion string     `json:"schemaVersion"`
	Bucket        BucketData `json:"bucket"`
	Object        ObjectData `json:"object"`
}

// BucketData struct
type BucketData struct {
	Name string `json:"name"`
	Arn  string `json:"arn"`
}

// ObjectData struct
type ObjectData struct {
	Key  string `json:"key"`
	Size string `json:"size"`
	Etag string `json:"eTag"`
}

const (
	helpText = `Create a batch of event notifications for a SQS`
	examples = `$ sando sqs create
	#Create a batch of events for a SQS using a file
	$ sando sqs create -f example_file.csv`
)

func NewCmdCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a batch of new SQS events",
		Long:    helpText,
		Example: examples,
		Run:     createBatch,
	}
}

func SetFlags(cmd *cobra.Command) {
	cmdcommon.SetCreateBatchFlags(cmd)
}

func createBatch(cmd *cobra.Command, _ []string) {
	params := parseFlags(cmd.Flags())
	err := func() error {
		s := cmdutil.Info("Creating new events...\n")
		defer s.Stop()
		batch, err := createEventList(params.file)
		if err != nil {
			return err
		}
		queueUrl, err := getQueueURL(params.queue)
		if err != nil {
			return err
		}
		if err := sendBatch(batch, queueUrl); err != nil {
			return err
		}
		return nil
	}()
	cmdutil.ExitIfError(err)
	cmdutil.Success("Messages sent\n")
}

func sendBatch(batch []*sqs.SendMessageBatchRequestEntry, queueURL *string) error {
	svc := sqs.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	batchEntry := []*sqs.SendMessageBatchRequestEntry{}
	for i, entry := range batch {
		if i%10 == 0 && i != 0 {
			cmdutil.Warn("Sending a batch of 10\n")
			out, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries:  batchEntry,
				QueueUrl: queueURL,
			})
			if err != nil {
				return err
			}
			batchEntry = nil
			fmt.Println(out)
		}
		batchEntry = append(batchEntry, entry)
	}
	return nil
}

func getQueueURL(queue string) (*string, error) {
	svc := sqs.New(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	res, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil {
		return nil, err
	}
	return res.QueueUrl, nil
}

func createEventList(file string) ([]*sqs.SendMessageBatchRequestEntry, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	batchEntry := []*sqs.SendMessageBatchRequestEntry{}
	for i, rec := range lines {
		event := S3Event{}
		event.Records = append(event.Records, RecordsData{
			EventVersion: "2.2",
			EventSource:  "aws:s3",
			AwsRegion:    rec[3],
			EventTime:    rec[1],
			EventName:    "recover-log",
			S3: S3Data{
				SchemaVersion: "1.0",
				Bucket: BucketData{
					Name: rec[4],
					Arn:  "arn:aws:s3:::" + rec[4],
				},
				Object: ObjectData{
					Key:  rec[0],
					Etag: rec[5],
					Size: rec[2],
				},
			},
		})
		eventBytes, _ := json.Marshal(event)
		entry := &sqs.SendMessageBatchRequestEntry{
			Id:           aws.String(fmt.Sprintf("%x", i)),
			DelaySeconds: aws.Int64(1),
			MessageBody:  aws.String(string(eventBytes)),
		}
		batchEntry = append(batchEntry, entry)
	}
	return batchEntry, nil
}

type createParams struct {
	file  string
	queue string
}

func parseFlags(flags query.FlagParser) *createParams {
	file, err := flags.GetString("file")
	cmdutil.ExitIfError(err)

	queue, err := flags.GetString("queue")
	cmdutil.ExitIfError(err)

	return &createParams{file: file, queue: queue}
}
