package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"

    "fmt"
    "flag"
    "os"
//    "time"
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

func get_snapshots(region, profileName, volumeId string) error {
    sess, err := session.NewSessionWithOptions(session.Options{
                            Profile: profileName,
                            Config: aws.Config{
                                Region: aws.String(region),
                            },
                            //SharedCoonfigState: session.SharedConfigEnable,
                        })
    if err != nil {
        fmt.Println("Error creation aws session")
        fmt.Println(err)
        return err
    }
    fmt.Printf("Volume id [%s]", volumeId)
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
    fmt.Println(result)
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
