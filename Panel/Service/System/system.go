/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供系统信息获取
 */

package System

import (
	"LoongPanel/Panel/Service/PanelLog"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// LoadAverage 负荷监控
func LoadAverage() ([3]float32, error) {
	// 负荷监控
	if Data.OSName == "windows" {
		// Windows 平台暂不支持
		return [3]float32{0, 0, 0}, nil
	}

	avg, err := load.Avg()
	if err != nil {
		PanelLog.ERROR("[系统监控] loadAverage1m() Error: ", err.Error())
		return [3]float32{0, 0, 0}, err
	}

	res := [3]float32{
		float32(avg.Load1),
		float32(avg.Load5),
		float32(avg.Load15),
	}
	return res, nil
}

func LoadAverage1m() (float32, error) {
	// 平均负荷 最近1分钟
	res, err := LoadAverage()
	return res[0], err
}

// MemoryUsage 内存使用率
func MemoryUsage() (float32, error) {
	res, err := mem.VirtualMemory()
	if err != nil {
		PanelLog.ERROR("[系统监控]", "MemoryUsage() Error: ", err.Error())
		return 0, err
	}
	return float32(res.UsedPercent), nil
}

// GetCpuUsage 获取CPU使用率
func GetCpuUsage() float32 {
	return float32(CPUPercent)
}

// diskReadIO 磁盘IO监控
func diskReadIO() {
	res := make(map[string]uint64)
	io, _ := disk.IOCounters()
	for key := range io {
		res[key] = io[key].ReadCount
	}
	DiskReadIO = res
	time.Sleep(1 * time.Second)
}

// diskWriteIO 磁盘IO监控
func diskWriteIO() {
	res := make(map[string]uint64)
	io, _ := disk.IOCounters()
	for key := range io {
		res[key] = io[key].WriteCount
	}
	DiskWriteIO = res
	time.Sleep(1 * time.Second)
}

// networkIO 网络IO监控
func networkIO() error {
	counters, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	var stats []NetworkIOStat
	for _, counter := range counters {
		stats = append(stats, NetworkIOStat{
			InterfaceName: counter.Name,
			BytesSent:     counter.BytesSent,
			BytesRecv:     counter.BytesRecv,
			PacketsSent:   counter.PacketsSent,
			PacketsRecv:   counter.PacketsRecv,
		})
	}
	time.Sleep(1 * time.Second)
	counters, err = net.IOCounters(true)
	if err != nil {
		return err
	}

	for idx, counter := range counters {
		if idx >= len(stats) {
			break
		}
		stats[idx] = NetworkIOStat{
			InterfaceName: counter.Name,
			BytesSent:     counter.BytesSent - stats[idx].BytesSent,
			BytesRecv:     counter.BytesRecv - stats[idx].BytesRecv,
			PacketsSent:   counter.PacketsSent - stats[idx].PacketsSent,
			PacketsRecv:   counter.PacketsRecv - stats[idx].PacketsRecv,
		}
	}

	for _, counter := range stats {
		NetworkIOSend += counter.BytesSent
		NetworkIORecv += counter.BytesRecv
		NetworkIOPacketsSent += counter.PacketsSent
		NetworkIOPacketsRecv += counter.PacketsRecv
	}
	return nil
}

// MonitorCPUPerCore 监控每个CPU核心的使用率
func MonitorCPUPerCore() ([]float64, error) {
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		PanelLog.ERROR("[系统监控]", "Error: ", err)
		return []float64{}, err
	}
	return percentages, nil
}

// GetOSData 获取系统信息
func GetOSData() (*OSData, error) {
	Data = &OSData{}
	cpuData, _ := cpu.Info()

	// CPU相关信息
	Data.CPUNumber = len(cpuData)
	Data.CPUCores = int(cpuData[0].Cores)
	Data.CPUName = cpuData[0].ModelName
	Data.CPUMHz = cpuData[0].Mhz
	// 系统相关
	Data.OSArch = runtime.GOARCH
	Data.OSName = runtime.GOOS
	Data.HostName, _ = os.Hostname()

	//Linux 内核版本
	if SkipWindows() {
		Data.LinuxVersion = "Windows平台"
	} else {
		Data.LinuxVersion = GetLinuxVersion()
	}
	// 网络相关
	Data.HostIP, _ = getLocalIP()
	// 内存
	Data.RAM, _ = getRAM()
	Data.Swap, _ = getSwap()
	Data.RAMMHz, _ = getRAMMHz()
	// Disk
	Data.Disks, _ = getDisk()
	for _, d := range Data.Disks {
		Data.DiskTotal += d.MaxMemory
	}

	// 包管理器
	Data.PkgManager = getPkgManager()

	return Data, nil
}

// GetRunTime 获取系统运行时间
func GetRunTime() string {
	if SkipWindows() {
		return "上次关机还是在上次"
	}
	out, err := exec.Command("uptime").Output()
	if err != nil {
		PanelLog.ERROR("[系统监控]", "GetRunTime() Error: ", err.Error())
		return ""
	}

	res := string(out)
	res = strings.Split(res, "up")[1]
	res = strings.Split(res, ",")[0]
	res = strings.Replace(res, "hour", "时", -1)
	res = strings.Replace(res, ":", "时", -1)
	res = strings.Replace(res, "min", "分", -1)
	res = strings.Replace(res, "days", "天", -1)
	res = strings.Replace(res, "day", "天", -1)
	res = strings.Replace(res, " ", "", -1)
	res = strings.Replace(res, ",", "", -1)
	res = strings.Replace(res, "\n", "", -1)
	res = strings.Replace(res, "\t", "", -1)
	return res
}

// GetLinuxVersion 获取Linux版本
func GetLinuxVersion() string {
	if SkipWindows() {
		return "Windows"
	}
	out, err := exec.Command("uname", "-sr").Output()
	if err != nil {
		PanelLog.ERROR("[系统监控]", "GetLinuxVersion() Error: ", err.Error())
		return ""
	}
	res := string(out)
	res = strings.Replace(res, "\n", "", -1)
	res = strings.Replace(res, "\t", "", -1)
	return res
}

// GetDiskUsage 获取磁盘使用
func GetDiskUsage() float32 {
	usage, _ := disk.Usage("/")
	return float32(usage.UsedPercent)
}

// Shutdown 关机
func Shutdown() {
	if SkipWindows() {
		return
	}
	err := exec.Command("shutdown", "-h", "now").Run()
	if err != nil {
		PanelLog.ERROR("[电源管理]", "Shutdown Error: ", err.Error())
		return
	}
}

// Reboot 重启
func Reboot() {
	if SkipWindows() {
		return
	}
	err := exec.Command("reboot").Run()
	if err != nil {
		PanelLog.ERROR("[电源管理]", "Shutdown Error: ", err.Error())
		return
	}
}

// getPkgManager 获取包管理器
func getPkgManager() string {
	if SkipWindows() {
		return "Windows"
	}
	_, err := exec.LookPath("apt")
	if err == nil {
		return "apt"
	}
	_, err = exec.LookPath("yum")
	if err == nil {
		return "yum"
	}
	return ""
}

// GetRAMUsedAndFree 获取内存剩余
func GetRAMUsedAndFree() (uint64, uint64) {
	memory, _ := mem.VirtualMemory()
	return memory.Free, memory.Used
}

// SetDNS 设置DNS
func SetDNS(dns string) {
	if SkipWindows() {
		return
	}
	err := exec.Command("echo", "nameserver "+dns, ">", "/etc/resolv.conf").Run()
	if err != nil {
		PanelLog.ERROR("[系统管理]", "SetDNS Error: ", err.Error())
		return
	}
}

// SetHOSTS 设置HOSTS
func SetHOSTS(hosts string) {
	if SkipWindows() {
		return
	}
	err := exec.Command("echo", hosts, ">", "/etc/hosts").Run()
	if err != nil {
		PanelLog.ERROR("[系统管理]", "SetHOSTS Error: ", err.Error())
		return
	}
}

// SetHOSTNAME 设置主机名
func SetHOSTNAME(hostname string) {
	if SkipWindows() {
		return
	}
	err := exec.Command("hostnamectl", "set-hostname", hostname).Run()
	if err != nil {
		PanelLog.ERROR("[系统管理]", "SetHOSTNAME Error: ", err.Error())
		return
	} else {
		Data.HostName = hostname
	}
}

// SetTimeZone 设置时区
func SetTimeZone(timezone string) {
	if SkipWindows() {
		return
	}
	err := exec.Command("timedatectl", "set-timezone", timezone).Run()
	if err != nil {
		PanelLog.ERROR("[系统管理]", "SetTimeZone Error: ", err.Error())
		return
	}
}

// TimeSync 时间同步
func TimeSync() {
	if SkipWindows() {
		return
	}
	err := exec.Command("ntpdate", "cn.pool.ntp.org").Run()
	if err != nil {
		PanelLog.ERROR("[系统管理]", "TimeSync Error: ", err.Error())
		return
	}
}
