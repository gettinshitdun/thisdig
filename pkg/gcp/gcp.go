package gcp

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/gettinshitdun/thisdig/pkg/dns"
	"github.com/gettinshitdun/thisdig/pkg/utils"
)

type InstanceType uint8


const (
	DIRECT InstanceType    = iota
	BEHIND_NLB
	BEHIND_GLB_INSTANCE
)

func (t InstanceType) String() string {
	return [...]string{
		"DIRECT",
		"BEHIND_NLB",
		"BEHIND_GLB",
	}[t]
}

type Instance struct {
	Typ      InstanceType
	IP       string
	HostName string
}

func (i * Instance) String() string {
	return fmt.Sprintf("{{ %s | %s | %s }}", i.Typ, i.IP, i.HostName)
}
type Instances struct {
	mutex *sync.Mutex
	instances map[string]*Instance // IP to Instance mapping
}

func (i * Instances) String() string {
	var out string = "\n{\n"

	for key, val := range i.instances {
		out += fmt.Sprintf("[ %s %s ]\n", key, val)
	}

	return out + "}"
}
 
func (i Instances) getRemainingIPs(allIPs dns.IPs) dns.IPs {
	keys := make([]string, 0, len(i.instances))
	for k := range i.instances {
		keys = append(keys, k)
	}
	if len(allIPs) == len(keys) {
		return dns.IPs{}
	}
	ipMap := make(map[string]struct{})
	for _, v := range keys {
		ipMap[v] = struct{}{}
	}

	var diff dns.IPs 
	for _, v := range allIPs {
		if _, found := ipMap[v]; !found {
			diff = append(diff, v)
		}
	}
	return diff
}

func (i Instances) addEntry(ip string, in *Instance) error {
	i.mutex.Lock()
	if _, exists := i.instances[ip]; !exists {
		i.instances[ip]	= in
		i.mutex.Unlock()
	} else {
		i.mutex.Unlock()
		return fmt.Errorf("mapping an ip: ip(%s) already mapped to a machine(%s)", ip, i.instances[ip].HostName)
	}
	return nil
}


type GCPQuerier struct {
	ips       dns.IPs
	instances *Instances
}

func New(ips dns.IPs) *GCPQuerier {
	return &GCPQuerier{
		ips:       ips,
		instances: &Instances{},
	}
}


func (q * GCPQuerier) GetMappedInstances() Instances {
	instances := Instances{
		mutex:     &sync.Mutex{},
		instances: map[string]*Instance{},
	}

	directProvider := &DirectInstanceProvider{}
	var remainingIps dns.IPs = directProvider.GetInstances(q.ips)
	fmt.Println(directProvider.instances)
	fmt.Println(remainingIps)
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

func (d * DirectInstanceProvider) GetInstances(ips dns.IPs) dns.IPs {
	var waitGroup * sync.WaitGroup = &sync.WaitGroup{}
	d.instances = &Instances{
		mutex:     &sync.Mutex{},
		instances: map[string]*Instance{},
	}

	for _, ip := range ips {
		waitGroup.Add(1)
		go searchAndAddInstance(ip, waitGroup, d.instances)
	}
	waitGroup.Wait()
	
	return d.instances.getRemainingIPs(ips)
}

func searchAndAddInstance(ip string, waitGroup * sync.WaitGroup, directInstances *Instances) {
	cmd := exec.Command(
		"gcloud", "compute", "instances", "list",
		"--zones=" + "us-east1-d",
		"--filter=networkInterfaces[].accessConfigs[].natIP=\""+ip+"\"",
		"--format=value(name,zone)",
	)
	out, err := cmd.Output()
	utils.HandleError(err, fmt.Sprintf("while running gcloud command for direct instance mode for '%s'", ip), true)
	fields := strings.Fields(string(out))
	hostName := fields[0]

	utils.HandleError(
		directInstances.addEntry(ip, &Instance{
			HostName: hostName,
			IP: ip,
			Typ: DIRECT,
		}),
		"while adding entry to directInstances",
		false,
	)
	waitGroup.Done()
}

func (d * DirectInstanceProvider) String() string {
	return fmt.Sprintf("Hello World %s", "hello world")
}

type NLBInstanceProvider struct {
	instances *Instances
}
func (n * NLBInstanceProvider) GetInstances() (*Instances, dns.IPs, error) {
	return nil, nil, nil
}
