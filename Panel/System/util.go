package System

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os/exec"
)

// getLocalIP 获取本地IP 数组
func getLocalIP() ([]string, error) {
	res := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("GetOSData() Error: ", err.Error())
	}
	for _, iface := range ifaces {
		adders := iface.Addrs
		if err != nil {
			fmt.Println("GetOSData() Error: ", err.Error())
		}
		for _, addr := range adders {
			res = append(res, addr.Addr)
		}
	}
	return res, nil
}

// getPublicIP 获取公网IP
func getPublicIP() (string, error) {
	out, err := exec.Command("curl", "ifconfig.me").Output()
	return string(out), err
}

// getRAM 获取最大内存
func getRAM() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("GetRAM() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getSwap 获取Swap内存
func getSwap() (float64, error) {
	memInfo, err := mem.SwapMemory()
	if err != nil {
		fmt.Println("GetSwap() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getDisk 获取磁盘容量
func getDisk() ([]*Disk, error) {
	res := make([]*Disk, 0)
	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("GetDisk() Error: ", err.Error())
	}
	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		res = append(res, &Disk{
			FileSystem:  partition.Device,
			MaxMemory:   float64(usage.Total),
			UsedMemory:  float64(usage.Used),
			MountedPath: partition.Mountpoint,
		})
	}

	return res, err
}