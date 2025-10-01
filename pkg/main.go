package main

import (
	"fmt"
	"os"

	"github.com/gettinshitdun/thisdig/pkg/dns"
	"github.com/gettinshitdun/thisdig/pkg/gcp"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("'thisdig' usage:\n\nthisdig <domain>\n\n"))
	}

	domain := os.Args[1]

	gcpQuerier := gcp.New(dns.New(domain).Query())

	_ = gcpQuerier.GetMappedInstances()
}
