/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：系统信息对象基本属性的定义
 */

package System

// OSData 系统信息
type OSData struct {
	OSName   string `json:"OSName"`   // 系统名称
	OSArch   string `json:"OSArch"`   // 系统架构
	HostName string `json:"HostName"` // 主机名称

	HostIP []string `json:"HostIP"` // 本地IP
	//PublicIP string   `json:"PublicIP"` // 公网IP 影响速度 取消

	RAM    float64 `json:"RAM"`    // 运行内存
	Swap   float64 `json:"Swap"`   // 交换空间内存
	RAMMHz any     `json:"RAMMHz"` // 内存频率

	CPUNumber int     `json:"CPUNumber"` // CPU 数量
	CPUCores  int     `json:"CPUCores"`  // CPU 核心数
	CPUName   string  `json:"CPUName"`   // CPU 名称
	CPUMHz    float64 `json:"CPUMHz"`    // CPU 频率

	LinuxVersion string `json:"LinuxVersion"` // Linux 版本

	Disks     []*Disk `json:"Disks"`     // 盘符
	DiskTotal float64 `json:"DiskTotal"` // 硬盘总容量

	PkgManager string `json:"PkgManager"` // 包管理器
}

// Disk 磁盘信息
type Disk struct {
	FileSystem  string  `json:"FileSystem"`  // 盘符名称
	MaxMemory   float64 `json:"MaxMemory"`   // 容量
	UsedMemory  float64 `json:"UsedMemory"`  // 已使用
	MountedPath string  `json:"MountedPath"` // 挂载位置
}

// NetworkIOStat 网络IO监控
type NetworkIOStat struct {
	InterfaceName string `json:"interfaceName"` // 网卡名称
	BytesSent     uint64 `json:"bytesSent"`     // 发送字节数
	BytesRecv     uint64 `json:"bytesRecv"`     // 接收字节数
	PacketsSent   uint64 `json:"packetsSent"`   // 发送数据包数
	PacketsRecv   uint64 `json:"packetsRecv"`   // 接收数据包数
}
