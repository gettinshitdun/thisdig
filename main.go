package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func getZoneName(zone string) (string, error) {
	switch zone {
	case "sg":
		return "asia-southeast1-b", nil
	case "eu":
		return "europe-west1-b", nil
	case "sc":
		return "us-east1-d", nil
	case "or":
		return "us-west1-b", nil
	default:
		return "", fmt.Errorf("unknown region: %s", zone)
	}
}

// func commonSuffix(strs []string) string {
//     if len(strs) == 0 {
//         return ""
//     }
//
//     // Helper: reverse a string
//     reverse := func(s string) string {
//         r := []rune(s)
//         for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
//             r[i], r[j] = r[j], r[i]
//         }
//         return string(r)
//     }
//
//     // Reverse all strings
//     reversed := make([]string, len(strs))
//     for i, s := range strs {
//         reversed[i] = reverse(s)
//     }
//
//     // Find common prefix among reversed strings
//     prefix := reversed[0]
//     for _, s := range reversed[1:] {
//         j := 0
//         for j < len(prefix) && j < len(s) && prefix[j] == s[j] {
//             j++
//         }
//         prefix = prefix[:j]
//     }
//
//     // Reverse the prefix to get common suffix
//     return reverse(prefix)
// }

func getGcloudInstance(ip, region string, hosts *Machines, waitGroup *sync.WaitGroup) {
	// Build the command
	zone, err := getZoneName(region)
	// fmt.Printf("Running gcloud search for %s %s\n", ip, zone)
	if err != nil {
		fmt.Printf("error1 while running gcloud command for %s\nerr: %v", ip, err)
		waitGroup.Done()
		return
	}
	cmd := exec.Command(
		"gcloud", "compute", "instances", "list",
		"--zones="+zone,
		"--filter=networkInterfaces[].accessConfigs[].natIP=\""+ip+"\"",
		"--format=value(name,zone)",
	)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error2 while running gcloud command for %s\nerr: %v\n", ip, string(out))
		waitGroup.Done()
		return
	}


	fields := strings.Fields(string(out))
	if len(fields) > 0 {
		buf := fields[0]
		hosts.add(buf, ip)
	}

	waitGroup.Done()
}

type Machines struct {
	mutex *sync.Mutex
	machines map[string]string
}

func (m *Machines) add(name, ip string) {
	m.mutex.Lock()
	m.machines[ip] = name
	m.mutex.Unlock()
}

func (m *Machines) String() string {
	return fmt.Sprint(m.machines)
}


func main() {
	domain := os.Args[1]
	region := os.Args[2]
	out, err := exec.Command("dig", domain, "+short").Output()
	if err != nil {
		_ = fmt.Errorf("while digging for %s\n Error: %v\n", domain, err)
	}

	ips := strings.Split(string(out), "\n")

	var hosts *Machines = &Machines{
		mutex:    &sync.Mutex{},
		machines: map[string]string{},
	}

	var waitGroup *sync.WaitGroup = &sync.WaitGroup{}

	fmt.Printf("Quering GCP for all the ips found behind %s\n", domain)
	for _, ip := range ips[:len(ips) - 1] {
		waitGroup.Add(1)
		go getGcloudInstance(ip, region, hosts, waitGroup)
	}

	waitGroup.Wait()
	

	fmt.Printf("%-15s %-20s\n", "IP", "Instance")
	fmt.Println(strings.Repeat("-", 40))

	for ip, host := range hosts.machines {
		fmt.Printf("%-15s %-20s\n", ip, host)
	}

	fmt.Println(strings.Repeat("-", 40))


	var regexFormats []string = []string{}
	var values []string = []string{}
	for _, host := range hosts.machines {
		values = append(values, host)
	}
	simpleRegex := strings.Join(values, "|")
	regexFormats = append(
		regexFormats,
		simpleRegex,
	)

	// commonSection := commonSuffix(hosts.machines)
	// fmt.Printf("common: (%s)\n",commonSection)
	// var uncommonSection string = ""
	// for idx, name := range hosts.machines {
	// 	buf := strings.Split(name, commonSection)[1]
	// 	uncommonSection += buf
	// 	if idx != len(hosts.machines) - 1 {
	// 		uncommonSection += "|" 
	// 	}
	// }
	//
	// regexFormats = append(
	// 	regexFormats,
	// 	fmt.Sprintf("~\"%s(%s)\"",
	// 		commonSection,
	// 		uncommonSection,
	// 	),
	// )

	fmt.Printf("Regex Formats:\n")
	for _, pattern := range regexFormats {
		fmt.Printf("~%s\n", pattern)
	}
}
