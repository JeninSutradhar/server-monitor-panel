package api

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Metric struct
type Metrics struct { // Public struct (by pascal casing (Uppercase )) to expose type variable information
	CPU               CPUStats              `json:"cpu"`
	Memory            MemoryStats           `json:"memory"`
	Disk              DiskStats             `json:"disk"`
	Net               NetworkStats          `json:"network"`
	Host              HostStats             `json:"host"`
	Load              LoadStats             `json:"load"`
	Processes         ProcessStats          `json:"process"`
	NetworkInterfaces NetworkInterfaceStats `json:"network_interfaces"`
}

type CPUStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Usage float64 `json:"usage"`
}

type MemoryStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}
type NetworkInterfaceStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Interfaces []NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Name       string   `json:"name"`
	MacAddress string   `json:"mac"`
	IPs        []string `json:"ips"`
}

type NetworkStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
}
type HostStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Uptime   uint64 `json:"uptime"`
	HostName string `json:"host_name"`
	Os       string `json:"os"`
}
type LoadStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
type ProcessStats struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Processes []Process `json:"processes"`
}

type Process struct { //Public struct (by pascal casing (Uppercase )) to expose type variable information
	Pid      int32   `json:"pid"`
	Name     string  `json:"name"`
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float32 `json:"mem_usage"`
	Username string  `json:"username"`
}

// GetCPUStats obtains the CPU information using Gopsutil lib ( public func or method using Camelcase!)
func GetCPUStats() (CPUStats, error) {
	perCPUUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return CPUStats{}, err
	}

	if len(perCPUUsage) > 0 {
		return CPUStats{
			Usage: perCPUUsage[0],
		}, nil
	}
	return CPUStats{}, fmt.Errorf("not enough info provided on usage by cpu percent")
}

// GetMemoryStats obtains memory information from using the gopsutil library ( public func or method using Camelcase!)
func GetMemoryStats() (MemoryStats, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return MemoryStats{}, err
	}

	return MemoryStats{
		Total:       memory.Total,
		Available:   memory.Available,
		Used:        memory.Used,
		UsedPercent: memory.UsedPercent,
	}, nil
}

// GetDiskStats retrieves total and used information by querying partitions and aggregting the totals.( public func or method using Camelcase!)
func GetDiskStats() (DiskStats, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskStats{}, err
	}

	return DiskStats{
		Total:       diskStat.Total,
		Free:        diskStat.Free,
		Used:        diskStat.Used,
		UsedPercent: diskStat.UsedPercent,
	}, nil
}

// GetNetworkStats retrieves current bytes sent/received information ( public func or method using Camelcase!)
func GetNetworkStats() (NetworkStats, error) {
	netInfo, err := net.IOCounters(false)
	if err != nil {
		return NetworkStats{}, err
	}
	if len(netInfo) > 0 {
		return NetworkStats{
			BytesSent: netInfo[0].BytesSent,
			BytesRecv: netInfo[0].BytesRecv,
		}, nil
	}

	return NetworkStats{}, fmt.Errorf("not enough info provided from net interface counters")

}

// GetNetworkInterfaces from all interfaces with the current implementation ( public func or method using Camelcase!)
func GetNetworkInterfaces() (NetworkInterfaceStats, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return NetworkInterfaceStats{}, err
	}
	var collectedInterfaces []NetworkInterface

	for _, value := range interfaces {

		addresses := make([]string, 0)

		for _, addr := range value.Addrs {
			addresses = append(addresses, addr.Addr)
		}
		networkInterface := NetworkInterface{
			Name:       value.Name,
			MacAddress: value.HardwareAddr,
			IPs:        addresses,
		}
		collectedInterfaces = append(collectedInterfaces, networkInterface)
	}

	return NetworkInterfaceStats{
		Interfaces: collectedInterfaces,
	}, nil
}

// GetHostStats obtains total host uptime in seconds ( public func or method using Camelcase!)
func GetHostStats() (HostStats, error) {
	h, err := host.Info()
	if err != nil {
		return HostStats{}, err
	}
	return HostStats{
		Uptime:   h.Uptime,
		HostName: h.Hostname,
		Os:       h.OS,
	}, nil
}

// GetLoadStats Obtains information related to CPU load. ( public func or method using Camelcase!)
func GetLoadStats() (LoadStats, error) {

	loadStat, err := load.Avg()
	if err != nil {
		return LoadStats{}, err
	}

	return LoadStats{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}, nil
}

// GetProcessStats list process of all available  ( public func or method using Camelcase!)
func GetProcessStats() (ProcessStats, error) {
	allProcesses, err := process.Processes()
	if err != nil {
		return ProcessStats{}, err
	}

	var processes []Process
	for _, proc := range allProcesses {

		name, err := proc.Name()
		if err != nil {
			//for now i want ignore process if it cant obtain their info.
			continue
		}
		memInfo, err := proc.MemoryPercent()

		if err != nil {
			continue

		}
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}
		username, err := proc.Username()
		if err != nil {
			continue
		}

		process := Process{
			Pid:      proc.Pid,
			Name:     name,
			MemUsage: memInfo,
			CPUUsage: cpuPercent,
			Username: username,
		}
		processes = append(processes, process)

	}

	return ProcessStats{
		Processes: processes,
	}, nil
}

// GetMetrics collects all system metrics and stores in the struct of type Metrics  ( public func or method using Camelcase!)
func GetMetrics() (Metrics, error) {
	cpu, err := GetCPUStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("Get CPU Usage: %v", err)
	}
	mem, err := GetMemoryStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("Get Memory Usage: %v", err)
	}
	disk, err := GetDiskStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("Get Disk Usage: %v", err)
	}
	net, err := GetNetworkStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("Get Network Usage: %v", err)
	}

	host, err := GetHostStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("get Host: %v", err)
	}

	load, err := GetLoadStats()
	if err != nil {
		return Metrics{}, fmt.Errorf("get Load Stats: %v", err)
	}

	processes, err := GetProcessStats()

	if err != nil {
		return Metrics{}, fmt.Errorf("get Process Stats %v", err)

	}

	interfaces, err := GetNetworkInterfaces()

	if err != nil {
		return Metrics{}, fmt.Errorf("get network interfaces: %v", err)
	}
	return Metrics{
		CPU:               cpu,
		Memory:            mem,
		Disk:              disk,
		Net:               net,
		Host:              host,
		Load:              load,
		Processes:         processes,
		NetworkInterfaces: interfaces,
	}, nil
}
