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
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/soap"

	"github.com/kubev2v/migration-planner/internal/gather_info"
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
	// Create a template output file

	// Create new client
	vCenterCredentials, err := NewVCenterClients(context.Background())
	if err != nil {
		return err
	}
	config := gather_info.Config{
		OutputFilePath: "",
		Client:         vCenterCredentials,
	}

	vmInfo := gather_info.VMInfo{
		Config:     &config,
		Log:        nil,
		OutputFile: "",
	}

	// gather VM information + update the output file (vm list , vm detailed, vm tools, vm networks,
	err = vmInfo.Run(context.Background())
	if err != nil {
		// gather all errors and continue
	}

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

func NewVCenterClients(ctx context.Context) (*gather_info.VCenterClients, error) {
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

	soapClient, err := govmomi.NewClient(ctx, u, creds.insecure)
	if err != nil {
		return nil, err
	}

	// Initialize the REST client
	restClient := rest.NewClient(soapClient.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		return nil, err
	}

	return &gather_info.VCenterClients{
		SOAPClient: soapClient,
		RESTClient: restClient,
	}, nil
}
