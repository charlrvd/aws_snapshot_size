package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ebs"
	"github.com/aws/aws-sdk-go/service/ec2"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type resultJsonFormat struct {
	Date *time.Time `json:"snapshot-date"`
	Id   string     `json:"snapshot-id"`
	Size string     `json:"snapshot-size"`
}

// error handler function based on the aws golang doc
func aws_err(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			log.Error(aerr.Error())
		}
	} else {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("AWS API call error")
	}
}

// shamefull copy
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func snapshots_size(region, profileName, snapshot1, snapshot2 string) (int64, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error creating aws session")
		return 0, err
	}
	svc := ebs.New(sess)
	input := &ebs.ListChangedBlocksInput{
		FirstSnapshotId:  aws.String(snapshot1),
		SecondSnapshotId: aws.String(snapshot2),
	}
	result, err := svc.ListChangedBlocks(input)
	if err != nil {
		aws_err(err)
		return 0, err
	}
	size := len(result.ChangedBlocks)
	block_size := result.BlockSize
	return int64(size) * *block_size, nil
}

func get_snapshots(region, profileName, volumeId, output string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error creation aws session")
		return err
	}
	svc := ec2.New(sess)
	input := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("status"),
				Values: []*string{
					aws.String("completed"),
				},
			},
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(volumeId),
				},
			},
		},
		OwnerIds: []*string{
			aws.String("self"),
		},
	}
	result, err := svc.DescribeSnapshots(input)
	if err != nil {
		aws_err(err)
		return err
	}
	for i, snap := range result.Snapshots {
		if i < len(result.Snapshots)-1 {
			size, err := snapshots_size(region, profileName, *snap.SnapshotId, *result.Snapshots[i+1].SnapshotId)
			if err != nil {
				return err
			}
			var message string
			switch output {
			case "json":
				jsonBuf, _ := json.Marshal(resultJsonFormat{
					Date: snap.StartTime,
					Id:   *snap.SnapshotId,
					Size: ByteCountIEC(size)})
				message = string(jsonBuf)
			default:
				message = fmt.Sprintf("Date: %-35s - ID: %s - Size: %s",
					snap.StartTime, *snap.SnapshotId, ByteCountIEC(size))
			}
			fmt.Println(message)
		}
	}
	return nil
}

func main() {
	// setup log
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// cli flags
	var region string
	var profileName string
	var volumeId string
	var output string
	flag.StringVar(&region, "r", "ca-central-1", "AWS region name")
	flag.StringVar(&profileName, "p", "default", "AWS config profile name")
	flag.StringVar(&volumeId, "v", "", "AWS volume id")
	flag.StringVar(&output, "o", "text", "Output format")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			message := fmt.Sprintf("Missing argument (%s)", f.Name)
			log.Fatal(message)
		}
	})

	err := get_snapshots(region, profileName, volumeId, output)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Execution failure")
	}
}
