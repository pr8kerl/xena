package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/mitchellh/cli"
	"net/url"
	"os"
	"strings"
)

type SnapshotsCommand struct {
	Name    string
	Latest  bool
	Summary bool
	Account string
	Ui      cli.Ui
}

func snapshotsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &SnapshotsCommand{
		Name:    "",
		Latest:  false,
		Summary: false,
		Account: "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorBlue,
		},
	}, nil
}

func (c *SnapshotsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("snapshots", flag.ContinueOnError)
	cmdFlags.StringVar(&c.Name, "name", "", "snapshot name value to match")
	cmdFlags.BoolVar(&c.Latest, "latest", false, "only show the latest snapshot")
	cmdFlags.BoolVar(&c.Summary, "summary", false, "only show snapshot id")
	cmdFlags.StringVar(&c.Account, "account", "", "the owner account id to filter snapshots with")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.Parse(args)

	if c.Account == "" {
		fmt.Printf("missing required account id option\n")
		return 1
	}

	//config := aws.NewConfig().WithRegion("ap-southeast-2")
	config := aws.NewConfig()
	sess := session.New(config)
	svc := ec2.New(sess)

	var filters []*ec2.Filter
	filters = make([]*ec2.Filter, 0, 10)

	filterState := ec2.Filter{
		Name: aws.String("status"),
		Values: []*string{
			aws.String("completed"),
		},
	}
	filterOwner := ec2.Filter{
		Name: aws.String("owner-id"),
		Values: []*string{
			aws.String(c.Account),
		},
	}
	filters = append(filters, &filterState)
	filters = append(filters, &filterOwner)

	params := &ec2.DescribeSnapshotsInput{
		DryRun:  aws.Bool(false),
		Filters: filters,
	}
	resp, err := svc.DescribeSnapshots(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return 1
	}

	// Pretty-print the response data.
	//fmt.Println(resp)
	c.printSnapshotInfo(resp)

	return 0
}

func (c *SnapshotsCommand) Help() string {
	return fmt.Sprintf("myaws snapshots --account <account id> [ --name <match snapshot name> ] [ --summary ] [ --latest]\n\t\tfind snapshots by name")
}

func (c *SnapshotsCommand) Synopsis() string {
	return "find all snapshots with tag Name matching values given on cmd line"
}

func (c *SnapshotsCommand) printSnapshotInfo(resp *ec2.DescribeSnapshotsOutput) {

	var latestSnap *ec2.Snapshot

SNAP:
	for _, snap := range resp.Snapshots {

		name := "nil"
		for _, keys := range snap.Tags {
			if *keys.Key == "Name" {
				name = url.QueryEscape(*keys.Value)
			}
		}
		if c.Name != "" {
			if !strings.Contains(name, c.Name) {
				// this snap name doesn't match
				continue SNAP
			}
		}

		if c.Latest {

			if latestSnap == nil {
				latestSnap = snap
			}
			if snap.StartTime.After(*latestSnap.StartTime) {
				latestSnap = snap
			}

		} else {
			c.PrintSnapshot(snap)
		}
	}
	if c.Latest && latestSnap != nil {
		c.PrintSnapshot(latestSnap)
	}
}

func (c *SnapshotsCommand) GetAccountAlias() (error, string) {

	svc := iam.New(session.New())

	params := &iam.ListAccountAliasesInput{}
	resp, err := svc.ListAccountAliases(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err, ""
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return nil, *resp.AccountAliases[0]

}

func (c *SnapshotsCommand) PrintSnapshot(snap *ec2.Snapshot) {

	var important_vals []*string
	name := ""
	snaptime := snap.StartTime.String()

	if c.Summary {
		important_vals = []*string{
			snap.SnapshotId,
		}
	} else {
		name = c.getSnapName(snap)
		important_vals = []*string{
			snap.SnapshotId,
			&name,
			&snaptime,
		}
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

func (c *SnapshotsCommand) getSnapName(snap *ec2.Snapshot) string {
	name := "nil"
	for _, keys := range snap.Tags {
		if *keys.Key == "Name" {
			name = url.QueryEscape(*keys.Value)
		}
	}
	return name
}
