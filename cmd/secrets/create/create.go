package create

import (
	"encoding/json"
	"sando/internal/cmdcommon"
	"sando/internal/cmdutil"
	"sando/internal/query"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"
)

const (
	helpText = `Create an AWS secret for SFTP`
	examples = `$ sando secrets create
	#Create a LDS secret for SFTP so logs are sent from Akamai
	$ sando secrets create -u clientName-amd-CPCODE -p myawesomesecret`
)

func NewCmdCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a AWS secret",
		Long:    helpText,
		Example: examples,
		Run:     create,
	}
}

func SetFlags(cmd *cobra.Command) {
	cmdcommon.SetCreateFlags(cmd, "Issue")
}

func create(cmd *cobra.Command, _ []string) {
	params := parseFlags(cmd.Flags())
	err := func() error {
		s := cmdutil.Info("Creating secret...")
		defer s.Stop()

		secret := map[string]string{
			"Password":      params.password,
			"Role":          "arn:aws:iam::749779118921:role/Transfer_for_SFTP",
			"HomeDirectory": "/bc-cdn-logs-data/sources/sftp/user/" + params.username,
			"Policy": `{
				"Version": "2012-10-17",
				"Statement": [
				  {
					"Sid": "AllowListingOfUserFolder",
					"Action": [
					  "s3:ListBucket"
					],
					"Effect": "Allow",
					"Resource": [
					  "arn:aws:s3:::${transfer:HomeBucket}"
					],
					"Condition": {
					  "StringLike": {
						"s3:prefix": [
						  "${transfer:HomeFolder}/*",
						  "${transfer:HomeFolder}"
						]
					  }
					}
				  },
				  {
					"Sid": "HomeDirObjectAccess",
					"Effect": "Allow",
					"Action": [
					  "s3:PutObject",
					  "s3:GetObject",
					  "s3:DeleteObjectVersion",
					  "s3:DeleteObject",
					  "s3:GetObjectVersion",
					  "s3:GetObjectACL",
					  "s3:PutObjectACL"
					],
					"Resource": "arn:aws:s3:::${transfer:HomeDirectory}*"
				  }
				]
			  }
			  `,
		}
		res, _ := json.Marshal(secret)
		svc := secretsmanager.New(session.New(&aws.Config{
			Region: aws.String("us-east-1"),
		}))

		input := &secretsmanager.CreateSecretInput{
			Name:         aws.String("SFTP/" + params.username),
			Description:  aws.String("CDN SFTP user account for log delivery"),
			SecretString: aws.String(string(res)),
		}

		_, err := svc.CreateSecret(input)
		if err != nil {
			return err
		}
		return nil
	}()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Secret created\n")
}

type createParams struct {
	username string
	password string
}

func parseFlags(flags query.FlagParser) *createParams {
	username, err := flags.GetString("username")
	cmdutil.ExitIfError(err)

	password, err := flags.GetString("password")
	cmdutil.ExitIfError(err)

	return &createParams{
		username: username,
		password: password,
	}
}
