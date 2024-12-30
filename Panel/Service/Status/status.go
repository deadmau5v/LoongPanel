/*
 * 创建人： deadmau5v
 * 创建时间： 2024-7-4
 * 文件作用：定义状态监控服务的基本结构和功能
 */

package Status

import (
	"LoongPanel/Panel/Service/Cron"
	"LoongPanel/Panel/Service/Database"
	notice "LoongPanel/Panel/Service/Notice"
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	StepTime time.Duration = 5 * time.Second // 默认5秒
	SaveTime time.Duration = 5 * time.Minute // 5分钟
	CronID   cron.EntryID
	// 初次运行时 磁盘IO归零
	DiskIOReadInit  uint64
	DiskIOWriteInit uint64
)

type LoadAverage [3]float32

type RAM [2]uint64

type DiskIO [2]uint64

type NetworkIO [4]uint64

type Status struct {
	Time        uint64      `json:"time"`
	LoadAverage LoadAverage `json:"load_average"`
	CPU         float32     `json:"cpu"`
	RAM         RAM         `json:"ram"`
	NetworkIO   NetworkIO   `json:"network_io"`
	DiskIO      DiskIO      `json:"disk_io"`
}

// 通用的 Scan 和 Value 方法，用于处理 JSON 序列化和反序列化
func Scan(value interface{}, target interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, target)
}

func Value(source interface{}) (driver.Value, error) {
	return json.Marshal(source)
}

// 实现各类型的 Scan 和 Value 方法
func (la *LoadAverage) Scan(value interface{}) error {
	return Scan(value, la)
}

func (la LoadAverage) Value() (driver.Value, error) {
	return Value(la)
}

func (r *RAM) Scan(value interface{}) error {
	return Scan(value, r)
}

func (r RAM) Value() (driver.Value, error) {
	return Value(r)
}

func (d *DiskIO) Scan(value interface{}) error {
	return Scan(value, d)
}

func (d DiskIO) Value() (driver.Value, error) {
	return Value(d)
}

func (n *NetworkIO) Scan(value interface{}) error {
	return Scan(value, n)
}

func (n NetworkIO) Value() (driver.Value, error) {
	return Value(n)
}

func sumNetworkIO(n map[string]uint64) uint64 {
	var sum uint64 = 0
	for _, v := range n {
		sum += v
	}
	return sum
}

func Job() {
	//PanelLog.DEBUG("[状态监控]", "保存服务器状态...")

	// 删除过期状态
	err := Database.DB.Where("time < ?", time.Now().Unix()-int64(SaveTime.Seconds())).Delete(&Status{}).Error
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
		return // 添加返回以防止后续代码在错误情况下执行
	}

	status := Status{}
	// 负载
	average, err := System.LoadAverage()
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
		return // 添加返回以防止后续代码在错误情况下执行
	}
	status.LoadAverage = average

	// CPU
	status.CPU = System.GetCpuUsage()

	// RAM
	ramFree, ramUsed := System.GetRAMUsedAndFree()
	status.RAM = [2]uint64{ramFree, ramUsed}

	// Disk
	if DiskIOReadInit == 0 && DiskIOWriteInit == 0 {
		r := sumNetworkIO(System.DiskReadIO)
		w := sumNetworkIO(System.DiskWriteIO)
		DiskIOReadInit = r
		DiskIOWriteInit = w
		status.DiskIO = [2]uint64{0, 0}
	} else {
		// 计算差值
		r := sumNetworkIO(System.DiskReadIO) - DiskIOReadInit
		w := sumNetworkIO(System.DiskWriteIO) - DiskIOWriteInit
		status.DiskIO = [2]uint64{r, w}
		DiskIOReadInit += r
		DiskIOWriteInit += w
	}

	// Network                             收                   发                     收包                         发包
	status.NetworkIO = [4]uint64{System.NetworkIORecv, System.NetworkIOSend, System.NetworkIOPacketsRecv, System.NetworkIOPacketsSent}
	System.NetworkIORecv, System.NetworkIOSend, System.NetworkIOPacketsRecv, System.NetworkIOPacketsSent = 0, 0, 0, 0

	// 时间
	status.Time = uint64(time.Now().Unix())

	// 保存状态
	err = Database.DB.Create(&status).Error
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
	}

	var settings []notice.UserNotificationSetting
	Database.DB.Preload("User").Find(&settings)
	for _, v := range settings {
		if v.NotifyOnCPU && status.CPU > v.MaxCPU {
			if v.NotifyIntervalLatestCPU == nil || time.Since(*v.NotifyIntervalLatestCPU) > time.Duration(v.NotifyInterval)*time.Minute {
				notice.SendMail(v.User.Mail, "CPU告警", fmt.Sprintf("警告CPU使用率过高: %f", status.CPU))
				now := time.Now()
				v.NotifyIntervalLatestCPU = &now
				Database.DB.Save(&v)
			}
		}
		if v.NotifyOnRAM && float64(status.RAM[0]+status.RAM[1])/float64(status.RAM[0])*100 > float64(v.MaxRAM) {
			if v.NotifyIntervalLatestRAM == nil || time.Since(*v.NotifyIntervalLatestRAM) > time.Duration(v.NotifyInterval)*time.Minute {
				notice.SendMail(v.User.Mail, "内存告警", fmt.Sprintf("警告内存使用率过高: %.2f%%", float64(status.RAM[0]+status.RAM[1])/float64(status.RAM[0])*100))
				now := time.Now()
				v.NotifyIntervalLatestRAM = &now
				Database.DB.Save(&v)
			}
		}
	}
}

func SetStepTime(t time.Duration) {
	StepTime = t
	Cron.Cron.Remove(CronID)
	if t != 0 {
		CronID, _ = Cron.Cron.AddFunc(Cron.DurationToCron(StepTime), Job)
	}
}

func SetSaveTime(t time.Duration) {
	SaveTime = t
}

func GetStatus(start uint64) []Status {
	status := []Status{}
	Database.DB.Where("time >= ?", start).Find(&status)
	return status
}

// 初始化
func init() {
	err := Database.DB.AutoMigrate(&Status{})
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
		return
	}

	CronID, err = Cron.Cron.AddFunc(Cron.DurationToCron(StepTime), Job)
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
		return
	}
}
