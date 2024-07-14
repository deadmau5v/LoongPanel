package notice

import (
	"LoongPanel/Panel/Service/Database"
	"time"

	"gorm.io/gorm"
)

type UserNotificationSetting struct {
	gorm.Model
	UserID         uint          `json:"user_id"`
	User           Database.User `json:"-" gorm:"foreignKey:UserID"`
	NotifyInterval uint          `json:"notify_interval" gorm:"default:5"` // 预警间隔 默认五分钟 单位分钟

	NotifyIntervalLatestCPU *time.Time `json:"notify_interval_latest_cpu"`         // 上一次预警
	NotifyOnCPU             bool       `json:"notify_on_cpu" gorm:"default:false"` // CPU告警
	MaxCPU                  float32    `json:"max_cpu" gorm:"default:80"`          // CPU警告值

	NotifyIntervalLatestRAM *time.Time `json:"notify_interval_latest_ram"`         // 上一次预警
	NotifyOnRAM             bool       `json:"notify_on_ram" gorm:"default:false"` // 内存告警
	MaxRAM                  float32    `json:"max_ram" gorm:"default:80"`          // 内存警告值

	ClamAVScanNotify bool `json:"clamav_scan_notify" gorm:"default:false"` // 病毒扫描告警

	InspectionNotify bool `json:"inspection_notify" gorm:"default:false"` // 巡检告警
}

func GetAllSettings() []UserNotificationSetting {
	var settings []UserNotificationSetting
	Database.DB.Find(&settings)
	return settings
}

func AddNotice(userID uint) {
	var setting UserNotificationSetting
	setting.UserID = userID
	Database.DB.Create(&setting)
}

func UpdateNotice(notice UserNotificationSetting) {
	Database.DB.Where("ID = ?", notice.ID).Updates(&notice)
}

func DeleteNotice(notice UserNotificationSetting) {
	Database.DB.Delete(&notice)
}

func init() {
	Database.DB.AutoMigrate(&UserNotificationSetting{})
}
