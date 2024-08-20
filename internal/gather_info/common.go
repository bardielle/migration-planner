package gather_info

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
)

type VCenterClients struct {
	SOAPClient *govmomi.Client
	RESTClient *rest.Client
}

type Config struct {
	OutputFilePath string
	Client         *VCenterClients
}
