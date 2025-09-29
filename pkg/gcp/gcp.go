package gcp

import (
	"github.com/gettinshitdun/thisdig/pkg/dns"
)

const (
	DIRECT_INSTANCE     = iota
	BEHIND_NLB_INSTANCE
	BEHIND_GLB_INSTANCE
)

type Instance struct {
	IP       string
	HostName string
}
type Instances map[string]Instance // IP to Instance mapping


type GCPQuerier struct {
	ips [] string
	instances *Instances
}

func New(ips dns.IPs) *GCPQuerier {
	return &GCPQuerier{
		ips:       ips,
		instances: &Instances{},
	}
}

// cmd := exec.Command(
// 	"gcloud", "compute", "instances", "list",
// 	"--zones="+zone,
// 	"--filter=networkInterfaces[].accessConfigs[].natIP=\""+ip+"\"",
// 	"--format=value(name,zone)",
// )

func (q * GCPQuerier) GetMappedInstances() Instances {
	instances := make(map[string]Instance)
	return instances
}

func (q * GCPQuerier) String() string { // will hold the final format of the strings
	return ""
}


type InstanceProvider interface {
	*Instances
	GetInstances() (*Instances, error)
}

type DirectInstanceProvider struct {
	instances *Instances
}

func (d * DirectInstanceProvider) GetInstances() (*Instances, error){
	return nil, nil
}


type NLBInstanceProvider struct {
	instances *Instances
}
func (n * NLBInstanceProvider) GetInstances() (*Instances, error) {
	return nil, nil
}
