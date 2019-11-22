package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	// Create aws Session uses default AWS creds /.aws
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: "sandbox",
		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Starts a new ec2 client
	svc := ec2.New(sess)

	//List all instances with tag
	if os.Args[1] == "LIST" {
		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:environment"),
					Values: []*string{
						aws.String("production"),
					},
				},
			},
		}

		result, err := svc.DescribeInstances(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println(result)

	}

	// Start Instance
	if os.Args[1] == "START" {
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]), aws.String(os.Args[3]),
			},
			DryRun: aws.Bool(true),
		}

		result, err := svc.StartInstances(input)
		awsErr, ok := err.(awserr.Error)

		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StartInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success", result.StartingInstances)
			}
		} else {
			fmt.Println("Error", err)
		}

	}

	// Stop Instance
	if os.Args[1] == "STOP" {
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(os.Args[2]), aws.String(os.Args[3]),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StopInstances(input)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success", result.StoppingInstances)
			}
		} else {
			fmt.Println("Error", err)
		}

	}
}
