package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kubev2v/migration-planner/pkg/log"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

func main() {
	command, err := NewGatherInformationCommand()
	if err != nil {
		os.Exit(1)
	}
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

type VCenterCredentials struct {
	url      string
	username string
	password string
	insecure bool
}

type gatherInfoCmd struct {
	log         *log.PrefixLogger
	credentials *VCenterCredentials
}

func NewGatherInformationCommand() (*gatherInfoCmd, error) {

	g := &gatherInfoCmd{
		log: log.NewPrefixLogger(""),
	}
	credentials, err := ReadCredentialFile()
	if err != nil {
		return nil, err
	}
	g.credentials = credentials

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Println("This program gather information about a vCenter.")
	}

	flag.Parse()

	//g.log.SetLevel()

	return g, nil
}

func (a *gatherInfoCmd) Execute() error {
	// Create new client
	// Create a template output file
	// gather VM information + update the output file (vm list , vm detailed, vm tools, vm networks,
	// gather ESXi information + update the output file (esxi_hosts  , powered on hosts)
	// gather dataceneter information + update the output file
	// gather folders information + update the output file
	// gather clusters information + update the output file
	// gather networks information + update the output file
	// gather portgroup and dvswitch information + update the output file
	// gather resource pool  information + update the output file
	// gather datastore information + update the output file
	// run the validator
	// aggregate report

	return nil
}

func ReadCredentialFile() (*VCenterCredentials, error) {
	file, err := os.Open(".creds")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	vcenterCreds := &VCenterCredentials{
		url:      "",
		username: "",
		password: "",
		insecure: false,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Remove whitespace and check if the line contains key="value"
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"`) // Remove quotes around the value

			// Map keys to struct fields
			switch key {
			case "url":
				vcenterCreds.url = value
			case "username":
				vcenterCreds.username = value
			case "password":
				vcenterCreds.password = value
			case "insecure":
				insecure, err := strconv.ParseBool(value)
				if err != nil {
					return nil, fmt.Errorf("invalid value for insecure: %v", err)
				}
				vcenterCreds.insecure = insecure
			}
		}
	}
	return vcenterCreds, nil
}

// NewClient creates a vim25.Client for use in the examples
func NewClient(ctx context.Context) (*vim25.Client, error) {
	creds, err := ReadCredentialFile()
	if err != nil {
		return nil, err
	}
	// Parse URL from string
	u, err := soap.ParseURL(creds.url)
	if err != nil {
		return nil, err
	}
	u.User = url.UserPassword(creds.username, creds.password)

	insecureDescription := fmt.Sprintf("Don't verify the server's certificate chain [%s]", creds.insecure)
	insecureFlag := flag.Bool("indecure", creds.insecure, insecureDescription)

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: *insecureFlag,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
