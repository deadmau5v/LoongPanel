/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：系统详细的全局变量
 */

package System

var WORKDIR string
var PublicIP string
var Data *OSData
var CPUPercent float64
var DiskReadIO map[string]uint64
var DiskWriteIO map[string]uint64
var NetworkIOSend uint64        // 网络发送
var NetworkIORecv uint64        // 网络接收
var NetworkIOPacketsSent uint64 // 网络发送数据包
var NetworkIOPacketsRecv uint64 // 网络接收数据包
