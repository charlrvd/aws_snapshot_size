package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/ebs"

    "fmt"
    "flag"
    "os"
)

// error handler function based on the aws golang doc
func aws_err(err error) {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        fmt.Println(err.Error())
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
        fmt.Println("Error creation aws session")
        fmt.Println(err)
        return 0, err
    }
    svc := ebs.New(sess)
    input := &ebs.ListChangedBlocksInput{
        FirstSnapshotId: aws.String(snapshot1),
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

func get_snapshots(region, profileName, volumeId string) error {
    sess, err := session.NewSessionWithOptions(session.Options{
                            Profile: profileName,
                            Config: aws.Config{
                                Region: aws.String(region),
                            },
                        })
    if err != nil {
        fmt.Println("Error creation aws session")
        fmt.Println(err)
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
    //fmt.Println(result)
    for i, snap := range result.Snapshots {
        if i < len(result.Snapshots)-1 {
            size, err := snapshots_size(region, profileName, *snap.SnapshotId, *result.Snapshots[i+1].SnapshotId)
            if err != nil {
                return err
            }
            message := fmt.Sprintf("Date: %-35s - ID: %s - Size: %s",
                                   snap.StartTime, *snap.SnapshotId, ByteCountIEC(size))
            fmt.Println(message)
        }
    }
    return nil
}

func main() {
    var region string
    var profileName string
    var volumeId string
    flag.StringVar(&region, "r", "ca-central-1", "AWS region name")
    flag.StringVar(&profileName, "p", "default", "AWS config profile name")
    flag.StringVar(&volumeId, "v", "", "AWS volume id")
    flag.Parse()
    flag.VisitAll(func (f *flag.Flag) {
        if f.Value.String() == "" {
            message := fmt.Sprintf("Missing argument (%s)", f.Name)
            fmt.Println(message)
            os.Exit(1)
        }
    })
    err := get_snapshots(region, profileName, volumeId)
    if err != nil {
        os.Exit(2)
    }
}
