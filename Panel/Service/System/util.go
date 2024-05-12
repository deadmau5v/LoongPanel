/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：系统信息相关工具类
 */

package System

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// getLocalIP 获取本地IP 数组
func getLocalIP() ([]string, error) {
	res := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		slog.Error("GetOSData() Error: ", err.Error())
	}
	for _, iface := range ifaces {
		adders := iface.Addrs
		if err != nil {
			slog.Error("GetOSData() Error: ", err.Error())
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
		slog.Error("GetRAM() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getSwap 获取Swap内存
func getSwap() (float64, error) {
	memInfo, err := mem.SwapMemory()
	if err != nil {
		slog.Error("GetSwap() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getDisk 获取磁盘容量
func getDisk() ([]*Disk, error) {
	res := make([]*Disk, 0)
	partitions, err := disk.Partitions(true)
	if err != nil {
		slog.Error("GetDisk() Error: ", err.Error())
	}
	for _, partition := range partitions {
		// 筛选不必要的空磁盘
		usage, _ := disk.Usage(partition.Mountpoint)
		if usage.Total != 0 {
			res = append(res, &Disk{
				FileSystem:  partition.Device,
				MaxMemory:   float64(usage.Total),
				UsedMemory:  float64(usage.Used),
				MountedPath: partition.Mountpoint,
			})
		}
	}

	return res, err
}

// getCPUPercent 获取CPU使用率
func getCPUPercent() float64 {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		slog.Error("Error: ", err)
		return 0
	}
	return percentages[0]
}

// SkipWindows 跳过Windows
func SkipWindows() bool {
	if Data.OSName == "windows" {
		return true
	} else {
		return false
	}
}

// getRAMMHz 获取内存频率
func getRAMMHz() (int, error) {
	if SkipWindows() {
		return 0, nil
	} else {
		out, err := exec.Command("dmidecode", "-t", "memory").Output()
		if err != nil {
			slog.Error("GetRAMMHz() Error: ", err.Error())
			return 0, err
		}
		res := string(out)
		res = strings.Split(res, "Speed: ")[1]
		res = strings.Split(res, " ")[0]
		resInt, err := strconv.Atoi(res)
		if err != nil {
			slog.Error("GetRAMMHz() Error: ", err.Error())
			return 0, err
		}
		return resInt, nil
	}
}
