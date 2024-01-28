package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type NetAdapter struct {
	Name                 string
	InterfaceDescription string
	IfIndex              int `json:"ifIndex"`
	Status               string
	MacAddress           string
	LinkSpeed            string
}

type NetAdapterStatistics struct {
	Name                   string
	ReceivedBytes          int
	ReceivedUnicastPackets int
	SentBytes              int
	SentUnicastPackets     int
}

func getNetAdapter(name string) (NetAdapter, error) {
	var adapter NetAdapter
	command := fmt.Sprintf("Get-NetAdapter -Name %s", name)
	err := runPowerShellCommandAndParseJson(command, &adapter)
	return adapter, err
}

func getNetAdapterStatistics(name string) (NetAdapterStatistics, error) {
	var statistics NetAdapterStatistics
	command := fmt.Sprintf("Get-NetAdapterStatistics -Name %s", name)
	err := runPowerShellCommandAndParseJson(command, &statistics)
	return statistics, err
}

func runPowerShellCommandAndParseJson(command string, v interface{}) error {
	output, err := runPowerShellCommand(command, "|", "ConvertTo-Json")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(output), v)
	if err != nil {
		return err
	}

	return nil
}

func runPowerShellCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("powershell", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a network interface name.")
		os.Exit(1)
	}

	interfaceName := os.Args[1]

	adapter, err := getNetAdapter(interfaceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if adapter.Status != "Up" {
		fmt.Printf("Interface %s is not up.\n", interfaceName)
		os.Exit(1)
	}

	fmt.Printf("Adapter: %v\n", adapter)

	var receivedBytes, sentBytes int

	for {
		statistics, err := getNetAdapterStatistics(interfaceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("[%s] Received: %.2f Mb\n", time.Now().Format("2006-01-02 15:04:05"), float64(statistics.ReceivedBytes-receivedBytes)/125000)
		fmt.Printf("[%s] Sent: %.2f Mb\n", time.Now().Format("2006-01-02 15:04:05"), float64(statistics.SentBytes-sentBytes)/125000)

		receivedBytes = statistics.ReceivedBytes
		sentBytes = statistics.SentBytes

		time.Sleep(1 * time.Second)
	}
}
