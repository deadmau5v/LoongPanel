package System

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LoadAverage 负荷监控
func LoadAverage() ([3]float32, error) {
	// 负荷监控
	out, err := exec.Command("uptime").Output()
	out = out[:len(out)-1] // 去除 \n 换行符

	if err != nil {
		fmt.Println("loadAverage1m() 1 Error: ", err.Error())
		// 出现问题就返回 0, 0, 0 默认值 不至于崩溃
		return [3]float32{0, 0, 0}, err
	}

	// 分割出“命令输出”的 Load Average 结果
	splits := strings.Split(string(out), ": ")
	out = []byte(splits[len(splits)-1])
	numbers := strings.Split(string(out), ", ")

	// 转换为 float32 返回
	res := [3]float32{}
	for idx := range numbers {
		number, err := strconv.ParseFloat(numbers[idx], 32)
		if err != nil {
			fmt.Println("loadAverage1m() 2 Error: ", err.Error())
			return [3]float32{0, 0, 0}, err
		}
		res[idx] = float32(number)
	}
	return res, nil

}
func LoadAverage1m() (float32, error) {
	// 平均负荷 最近1分钟
	res, err := LoadAverage()
	return res[0], err
}
func LoadAverage5m() (float32, error) {
	// 平均负荷 最近5分钟
	res, err := LoadAverage()
	return res[1], err
}
func LoadAverage15m() (float32, error) {
	// 平均负荷 最近15分钟
	res, err := LoadAverage()
	return res[2], err
}

// CPU CPU负荷监控
func CPU() (float64, error) {
	// CPU负荷
	res, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("CPU() Error: ", err.Error())
		return 0, nil
	}
	return res[0], nil
}

func Memory() (uint64, error) {
	res, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Memory() Error: ", err.Error())
		return 0, err
	}
	return res.Used, nil
}

// DiskReadIO 磁盘IO监控
func DiskReadIO() (map[string]uint64, error) {
	res := make(map[string]uint64)
	io, _ := disk.IOCounters()
	for key := range io {
		res[key] = io[key].ReadCount
	}
	time.Sleep(1*time.Second + 100*time.Millisecond)
	io, _ = disk.IOCounters()
	for key := range io {
		res[key] = io[key].ReadCount - res[key]
	}
	return res, nil
}

// DiskWriteIO 磁盘IO监控
func DiskWriteIO() (map[string]uint64, error) {
	res := make(map[string]uint64)
	io, _ := disk.IOCounters()
	for key := range io {
		res[key] = io[key].WriteCount
	}
	time.Sleep(1*time.Second + 100*time.Millisecond)
	io, _ = disk.IOCounters()
	for key := range io {
		res[key] = io[key].WriteCount - res[key]
	}
	return res, nil
}

// NetworkIO 网络IO监控
func NetworkIO() ([]NetworkIOStat, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
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

	counters, err = net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	for idx, counter := range counters {
		stats[idx] = NetworkIOStat{
			InterfaceName: counter.Name,
			BytesSent:     counter.BytesSent - stats[idx].BytesSent,
			BytesRecv:     counter.BytesRecv - stats[idx].BytesRecv,
			PacketsSent:   counter.PacketsSent - stats[idx].PacketsSent,
			PacketsRecv:   counter.PacketsRecv - stats[idx].PacketsRecv,
		}
	}

	return stats, nil
}

// 重要
// Todo 持久化记录
// 		Todo 持久化记录 开关
// 		Todo 持久化记录 时长
// 		Todo 持久化记录 位置
// 		Todo 持久化记录 清理

// MonitorCPUPerCore 监控每个CPU核心的使用率
func MonitorCPUPerCore() ([]float64, error) {
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		fmt.Println("Error: ", err)
		return []float64{}, err
	}
	return percentages, nil
}

// GetOSData 获取系统信息
func GetOSData() (*OSData, error) {
	Data = &OSData{}
	cpuData, err := cpu.Info()

	// CPU相关信息
	Data.CPUNumber = len(cpuData)
	Data.CPUCores = int(cpuData[0].Cores)
	Data.CPUName = cpuData[0].ModelName
	Data.CPUMHz = cpuData[0].Mhz
	// 系统相关
	Data.OSArch = runtime.GOARCH
	Data.OSName = runtime.GOOS
	Data.HostName, err = os.Hostname()

	//Linux 内核版本
	if Data.OSName == "linux" {
		Data.LinuxVersion = runtime.GOOS + " " + runtime.GOARCH
	} else if Data.OSName == "windows" {
		Data.LinuxVersion = "windows无法获取"
	}
	// 网络相关
	Data.HostIP, err = getLocalIP()
	Data.RAM, err = getRAM()
	Data.Swap, err = getSwap()
	// Disk
	Data.Disks, err = getDisk()

	return Data, err
}

// GetRunTime 获取系统运行时间
func GetRunTime() string {
	if Data.OSName == "windows" {
		return "Windows无法获取"
	}
	out, err := exec.Command("uptime").Output()
	if err != nil {
		fmt.Println("GetRunTime() Error: ", err.Error())
		return ""
	}
	return string(out)
}
