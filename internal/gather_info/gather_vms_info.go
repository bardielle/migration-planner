package gather_info

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubev2v/migration-planner/pkg/log"
	"github.com/vmware/govmomi/vapi/rest"
	"net/http"
	"net/url"
)

type VMInfo struct {
	Config     *Config
	Log        *log.PrefixLogger
	OutputFile string
}

// VM represents minimal information about a VM from the vSphere REST API
type VM struct {
	VM         string `json:"vm"`
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
}

// VMDetails represents detailed information about a VM from the vSphere REST API
type VMDetails struct {
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
	GuestOS    string `json:"guest_os"`
	CPU        struct {
		Count int `json:"count"`
	} `json:"cpu"`
	Memory struct {
		SizeMiB int `json:"size_MiB"`
	} `json:"memory"`
}

func (i *VMInfo) Run(ctx context.Context) error {
	//soapClient := i.Config.Client.SOAPClient
	restClient := i.Config.Client.RESTClient

	// List all VMs using the REST API
	i.listVMs(ctx, restClient)

	// Loop through each VM and retrieve details
	for _, vm := range vms {

	}

	return nil
}

func (i *VMInfo) listVMs(ctx context.Context, restClient *rest.Client) ([]VM, error) {
	httpClient := restClient.Client
	listVMsURL := restClient.URL().ResolveReference(&url.URL{Path: "/rest/vcenter/vm"})
	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listVMsURL.String(), nil)
	if err != nil {
		return nil, err
	}
	var vms []VM
	// Define a function to process the response
	handleResponse := func(res *http.Response) error {
		defer res.Body.Close()

		// Check for HTTP errors
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		// Decode the response JSON into a slice of VMs
		if err := json.NewDecoder(res.Body).Decode(&vms); err != nil {
			return err
		}

		// Print out the VM names
		for _, vm := range vms {
			fmt.Printf("VM Name: %s, Power State: %s\n", vm.Name, vm.PowerState)
		}

		return nil
	}

	// Perform the HTTP request
	err = httpClient.Do(ctx, req, handleResponse)
	if err != nil {
		return nil, err
	}
	return vms, nil
}
