package main

import (
	"fmt"
	"github.com/mostlygeek/arp"
)

func main() {
	for ip, mac := range arp.Table() {
		fmt.Printf("%s : %s\n", ip, mac)
	}
}
