package Status

import (
	"LoongPanel/Panel/Service/Cron"
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"time"
)

var (
	StepTime time.Duration = 5 * time.Second // 默认五秒
	SaveTime time.Duration = 1 * time.Hour   // 保存一小时
	CronID   cron.EntryID
)

type LoadAverage [3]float32

func (la *LoadAverage) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, la)
}

func (la LoadAverage) Value() (driver.Value, error) {
	return json.Marshal(la)
}

type RAM [2]uint64

func (r *RAM) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, r)
}

func (r RAM) Value() (driver.Value, error) {
	return json.Marshal(r)
}

type DiskIO []map[string]uint64

func (d *DiskIO) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, d)
}

func (d DiskIO) Value() (driver.Value, error) {
	return json.Marshal(d)
}

type NetworkIO [4]uint64

func (n *NetworkIO) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, n)
}

func (n NetworkIO) Value() (driver.Value, error) {
	return json.Marshal(n)
}

type Status struct {
	Time        uint64      `json:"time"`
	LoadAverage LoadAverage `json:"load_average"`
	CPU         float32     `json:"cpu"`
	RAM         RAM         `json:"ram"`
	NetworkIO   NetworkIO   `json:"network_io"`
	DiskIO      DiskIO      `json:"disk_io"`
}

func Job() {
	//PanelLog.DEBUG("[状态监控]", "保存服务器状态...")

	// 删除过期状态
	err := Database.DB.Where("time < ?", time.Now().Unix()-int64(SaveTime.Seconds())).Delete(&Status{}).Error
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
	}

	status := Status{}
	// 负载
	average, err := System.LoadAverage()
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
	}
	status.LoadAverage = average
	// CPU
	cpu := System.GetCpuUsage()
	status.CPU = cpu

	// RAM
	ramFree, ramUsed := System.GetRAMUsedAndFree()
	status.RAM = [2]uint64{ramFree, ramUsed}

	// Disk
	diskIO := []map[string]uint64{System.DiskWriteIO, System.DiskReadIO}
	status.DiskIO = diskIO

	// Network                          发送                 接收                   发送包                           接收包
	networkIO := [4]uint64{System.NetworkIOSend, System.NetworkIORecv, System.NetworkIOPacketsSent, System.NetworkIOPacketsRecv}
	status.NetworkIO = networkIO
	System.NetworkIOSend, System.NetworkIORecv, System.NetworkIOPacketsSent, System.NetworkIOPacketsRecv = 0, 0, 0, 0

	// 时间
	status.Time = uint64(time.Now().Unix())

	// 保存状态
	err = Database.DB.Create(&status).Error
	if err != nil {
		PanelLog.ERROR("[状态监控]", err.Error())
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
