package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/cli"
	"net/url"
	"os"
	"strings"
)

type InstancesCommand struct {
	Role string
	Env  string
	Ui   cli.Ui
}

func instancesCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &InstancesCommand{
		Role: "",
		Env:  "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorBlue,
		},
	}, nil
}

func (c *InstancesCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("instances", flag.ContinueOnError)
	cmdFlags.StringVar(&c.Role, "role", "", "role tag value to match")
	cmdFlags.StringVar(&c.Env, "environment", "", "environment tag value to match")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.Parse(args)
	var numFlags = 0

	config := aws.NewConfig().WithRegion("ap-southeast-2")
	sess := session.New(config)
	svc := ec2.New(sess)

	var filters []*ec2.Filter
	filters = make([]*ec2.Filter, 0, 10)

	if c.Role != "" {
		numFlags++
		filter := ec2.Filter{
			Name: aws.String("tag:Role"),
			Values: []*string{
				aws.String(c.Role),
			},
		}
		filters = append(filters, &filter)
	}
	if c.Env != "" {
		numFlags++
		filter := ec2.Filter{
			Name: aws.String("tag:Environment"),
			Values: []*string{
				aws.String(c.Env),
			},
		}
		filters = append(filters, &filter)
	}
	if numFlags == 0 {
		//		cmdFlags.Usage()
		return 1
	}

	params := &ec2.DescribeInstancesInput{
		DryRun:  aws.Bool(false),
		Filters: filters,
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return 1
	}

	// Pretty-print the response data.
	//	fmt.Println(resp.Reservations)
	c.printInstanceInfo(resp)

	return 0
}

func (c *InstancesCommand) Help() string {
	return "myaws instances: find instances by tag Role and/or Environment"
}

func (c *InstancesCommand) Synopsis() string {
	return "find all instances with tag Role and Environment matching values given on cmd line"
}

func (c *InstancesCommand) printInstanceInfo(resp *ec2.DescribeInstancesOutput) {

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {

			// We need to see if the Name is one of the tags. It's not always
			// present and not required in Ec2.
			name := "nil"
			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					name = url.QueryEscape(*keys.Value)
				}
			}

			important_vals := []*string{
				inst.InstanceId,
				&name,
				inst.PrivateIpAddress,
				inst.InstanceType,
				inst.PublicIpAddress,
			}

			// Convert any nil value to a printable string in case it doesn't
			// doesn't exist, which is the case with certain values
			output_vals := []string{}
			for _, val := range important_vals {
				if val != nil {
					output_vals = append(output_vals, *val)
				} else {
					output_vals = append(output_vals, "nil")
				}
			}
			// The values that we care about, in the order we want to print them
			fmt.Println(strings.Join(output_vals, ","))
		}
	}
}
