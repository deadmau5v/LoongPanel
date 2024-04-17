package Status

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type OSData struct {
	OSName   string // 系统名称
	OSArch   string // 系统架构
	HostName string // 主机名称

	HostIP   []string // 本地IP
	PublicIP string   // 公网IP

	RAM  float64 // 运行内存
	Swap float64 // 交换空间内存

	CPUCores   int      // CPU 核心数
	CPUThreads int      // CPU 线程数
	CPUName    int      // CPU 名称
	CPUArch    int      // CPU 架构
	CPUMaxMHz  float64  // CPU 最大赫兹
	CPUMinMHz  float64  // CPU 最小赫兹
	CPUMods    []string // CPU 指令集

	Disks []*Disk // 盘符
}

type Disk struct {
	FileSystem  string  // 盘符名称
	MaxMemory   float64 // 容量
	UsedMemory  float64 // 已使用
	MountedPath string  // 挂载位置
}

func LoadAverage() [3]float32 {
	// 负荷监控
	out, err := exec.Command("uptime").Output()
	out = out[:len(out)-1] // 去除 \n 换行符

	if err != nil {
		fmt.Println("loadAverage1m() 1 Error: ", err.Error())
		// 出现问题就返回 0, 0, 0 默认值 不至于崩溃
		return [3]float32{0, 0, 0}
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
			return [3]float32{0, 0, 0}
		}
		res[idx] = float32(number)
	}
	return res

}
func LoadAverage1m() float32 {
	// 平均负荷 最近1分钟
	return LoadAverage()[0]
}
func LoadAverage5m() float32 {
	// 平均负荷 最近5分钟

	return LoadAverage()[1]
}
func LoadAverage15m() float32 {
	// 平均负荷 最近15分钟
	return LoadAverage()[2]
}

func CPUPercent() float64 {
	// CPU负荷
	res, _ := cpu.Percent(time.Second, false)
	return res[0]
}

func MemroyPercent() float64 {
	res, _ := mem.VirtualMemory()
	return res.UsedPercent
}

// 重要
// Todo 磁盘IO监控
// Todo 网络IO监控
// Todo 持久化记录
// 		Todo 持久化记录 开关
// 		Todo 持久化记录 时长
// 		Todo 持久化记录 位置
// 		Todo 持久化记录 清理

// 创新
// Todo CPU核心实时监控
