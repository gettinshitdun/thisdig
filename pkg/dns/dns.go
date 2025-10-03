package dns

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gettinshitdun/thisdig/pkg/utils"
)

type DNSQuerier struct {
	domain string
}

func New(domain string) *DNSQuerier {
	return &DNSQuerier{
		domain: domain,
	}
}

type IPs []string

func (q *DNSQuerier) Query() IPs {
	out, err := exec.Command("dig", q.domain, "+short").Output()
	utils.HandleError(err, fmt.Sprintf("while executing dig for %s", q.domain), false)
	var ips IPs = strings.Split(string(out), "\n")
	return ips[:len(ips) - 1]
}
