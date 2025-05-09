package pkg

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func CollectMetrics(lastNetIn, lastNetOut uint64, intervalSec float64) (map[string]interface{}, uint64, uint64, error) {
	// CPU
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, lastNetIn, lastNetOut, err
	}

	// Memory
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, lastNetIn, lastNetOut, err
	}

	// Disk
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, lastNetIn, lastNetOut, err
	}

	// Network
	netIOs, err := net.IOCounters(false)
	if err != nil || len(netIOs) == 0 {
		return nil, lastNetIn, lastNetOut, err
	}
	netIn := netIOs[0].BytesRecv
	netOut := netIOs[0].BytesSent

	var inRate, outRate float64
	if lastNetIn > 0 && lastNetOut > 0 {
		inRate = float64(netIn-lastNetIn) / intervalSec
		outRate = float64(netOut-lastNetOut) / intervalSec
	}

	// Public IP
	publicIP, err := GetPublicIP()
	if err != nil {
		publicIP = "unknown"
	}

	// Hostname
	hostname := GetHostname()

	data := map[string]interface{}{
		"hostname":          hostname,
		"public_ip":         publicIP,
		"cpu_total":         100,
		"cpu_used_percent":  cpuPercent[0],
		"mem_total_mb":      vm.Total / 1024 / 1024,
		"mem_used_mb":       (vm.Total - vm.Available) / 1024 / 1024,
		"mem_used_percent":  vm.UsedPercent,
		"disk_total_gb":     diskStat.Total / 1024 / 1024 / 1024,
		"disk_used_gb":      diskStat.Used / 1024 / 1024 / 1024,
		"disk_used_percent": diskStat.UsedPercent,
		"net_in_bytes":      netIn,
		"net_out_bytes":     netOut,
		"net_in_rate_bps":   inRate,
		"net_out_rate_bps":  outRate,
		"timestamp":         time.Now().Format(time.RFC3339),
	}

	return data, netIn, netOut, nil
}
