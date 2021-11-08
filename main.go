package main

import (
	"log"
	"os/exec"
	"strconv"
)

var (
	inChain          = "PROXY_INIT_REDIRECT"
	outChain         = "PROXY_INIT_OUTPUT"
	proxyOutPort     = 62082
	proxyInPort      = 62081
	controlPlanePort = 62000
	userID           = 2102
)

func outbound() []*exec.Cmd {
	commands := make([]*exec.Cmd, 0)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-N", inChain,
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", "PREROUTING", "-j", inChain,
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", inChain, "-p", "tcp", "--match", "multiport", "--dports", "62000", "-j", "RETURN",
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", inChain, "-p", "tcp", "-j", "REDIRECT", "--to-port", strconv.Itoa(proxyInPort),
		),
	)
	return commands
}

func inbound() []*exec.Cmd {
	commands := make([]*exec.Cmd, 0)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-N", outChain,
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", "OUTPUT", "-j", outChain,
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", inChain, "-m", "owner", "--uid-owner", strconv.Itoa(userID), "-j", "RETURN",
		),
	)
	commands = append(commands,
		exec.Command(
			"iptables", "-t", "nat", "-A", outChain, "-p", "tcp", "-j", "REDIRECT", "--to-port", strconv.Itoa(proxyOutPort),
		),
	)
	return commands
}

func main() {
	outCommands := outbound()
	inCommands := inbound()

	for _, cmd := range outCommands {
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Error:", err.Error())
			continue
		}
		if len(out) > 0 {
			log.Println("Output:", out)
		}
	}
	for _, cmd := range inCommands {
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Error:", err.Error())
			continue
		}
		if len(out) > 0 {
			log.Println("Output:", out)
		}
	}
}
