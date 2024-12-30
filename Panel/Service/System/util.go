/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：系统信息相关工具类
 */

package System

import (
	"LoongPanel/Panel/Service/PanelLog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// getLocalIP 获取本地IP 数组
func getLocalIP() ([]string, error) {
	res := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		PanelLog.ERROR("[系统监控]", "GetOSData() Error: ", err.Error())
	}
	for _, iface := range ifaces {
		adders := iface.Addrs
		if err != nil {
			PanelLog.ERROR("[系统监控]", "GetOSData() Error: ", err.Error())
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
		PanelLog.ERROR("[系统监控]", "GetRAM() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getSwap 获取Swap内存
func getSwap() (float64, error) {
	memInfo, err := mem.SwapMemory()
	if err != nil {
		PanelLog.ERROR("[系统监控]", "GetSwap() Error: ", err.Error())
	}
	return float64(memInfo.Total), err
}

// getDisk 获取磁盘容量
func getDisk() ([]*Disk, error) {
	res := make([]*Disk, 0)
	partitions, err := disk.Partitions(true)
	if err != nil {
		PanelLog.ERROR("[系统监控]", "GetDisk() Error: ", err.Error())
	}
	for _, partition := range partitions {
		// 筛选不必要的空磁盘
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			if err.Error() == "operation not permitted" {
				PanelLog.WARN("GetDisk() Error: 权限不足警告")
			} else {
				PanelLog.ERROR("GetDisk() Error: ", err.Error())
			}
			continue
		}
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
		PanelLog.ERROR("[系统监控]", "Error: ", err)
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
			// executable file not found in
			if strings.Contains(err.Error(), "executable file not found in") {
				PanelLog.WARN("[系统监控] GetRAMMHz() Error: dmidecode 命令不存在")
			} else {
				PanelLog.ERROR("[系统监控] GetRAMMHz() Error: ", err.Error())
			}
			return 0, err
		}
		res := string(out)
		if strings.Contains(res, "Speed: Unknown") {
			// 如果频率未知
			PanelLog.DEBUG("[系统监控] GetRAMMHz() 频率未知")
			return 0, nil
		}
		if strings.Contains(res, "Permission denied") {
			// 如果权限不足
			PanelLog.WARN("[系统监控] GetRAMMHz() Error: 权限不足警告")
			return 0, err
		}

		// 正常情况
		if strings.Contains(res, "Speed: ") {
			res = strings.Split(res, "Speed: ")[1]
			res = strings.Split(res, " ")[0]
			resInt, err := strconv.Atoi(res)
			if err != nil {
				PanelLog.ERROR("[系统监控]", "GetRAMMHz() Error: ", err.Error())
				return 0, err
			}
			return resInt, nil
		} else {
			return 0, nil
		}
	}
}
