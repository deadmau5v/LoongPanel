/*
 * 创建人： deadmau5v
 * 创建时间： 2024-7-4
 * 文件作用： 面板定时任务管理
 */

package Cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	Cron *cron.Cron
)

func DurationToCron(d time.Duration) string {
	return fmt.Sprintf("@every %s", d.String())
}

// 统计所有任务数量
func AllJobCount() int {
	return len(Cron.Entries())
}

func init() {
	Cron = cron.New(cron.WithSeconds())
	Cron.Start()
}
