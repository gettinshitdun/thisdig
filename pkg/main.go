package main

import (
	"fmt"
	"os"

	"github.com/gettinshitdun/thisdig/pkg/dns"
	"github.com/gettinshitdun/thisdig/pkg/gcp"
)
// func getGcloudInstance(ip, region string, hosts *Machines, waitGroup *sync.WaitGroup) {
// 	// Build the command
// 	zone, err := getZoneName(region)
// 	// fmt.Printf("Running gcloud search for %s %s\n", ip, zone)
// 	if err != nil {
// 		fmt.Printf("error1 while running gcloud command for %s\nerr: %v", ip, err)
// 		waitGroup.Done()
// 		return
// 	}
// 	cmd := exec.Command(
// 		"gcloud", "compute", "instances", "list",
// 		"--zones="+zone,
// 		"--filter=networkInterfaces[].accessConfigs[].natIP=\""+ip+"\"",
// 		"--format=value(name,zone)",
// 	)
// 	out, err := cmd.Output()
// 	if err != nil {
// 		fmt.Printf("error2 while running gcloud command for %s\nerr: %v\n", ip, string(out))
// 		waitGroup.Done()
// 		return
// 	}
//
//
// 	fields := strings.Fields(string(out))
// 	if len(fields) > 0 {
// 		buf := fields[0]
// 		hosts.add(buf, ip)
// 	}
//
// 	waitGroup.Done()
// }
//



func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("'thisdig' usage:\n\nthisdig <domain>\n\n"))
	}

	domain := os.Args[1]

	gcpQuerier := gcp.New(dns.New(domain).Query())

	machines :=  gcpQuerier.GetMappedInstances()

	fmt.Print(machines)
}
