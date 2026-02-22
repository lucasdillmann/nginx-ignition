package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var out bytes.Buffer
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/netbirdio/netbird")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to locate netbird module: %v\n", err)
		os.Exit(0)
	}

	modDir := strings.TrimSpace(out.String())
	if modDir == "" {
		os.Exit(0)
	}

	targetFile := filepath.Join(modDir, "client", "firewall", "uspfilter", "forwarder", "forwarder.go")
	content, err := os.ReadFile(targetFile)
	if err != nil {
		os.Exit(0)
	}

	oldStr := "udpForwarder := udp.NewForwarder(s, f.handleUDP)"
	newStr := `udpForwarder := udp.NewForwarder(s, func(request *udp.ForwarderRequest) {
		f.handleUDP(request)
	})`

	if strings.Contains(string(content), oldStr) {
		newContent := strings.ReplaceAll(string(content), oldStr, newStr)

		info, err := os.Stat(targetFile)
		if err == nil {
			err = os.Chmod(targetFile, info.Mode()|0200)
			if err != nil {
				fmt.Printf("Warning: failed to make file writable: %v\n", err)
			}
		}

		err = os.WriteFile(targetFile, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("Failed to write patched file: %v\n", err)
			os.Exit(1)
		}
	}
}
