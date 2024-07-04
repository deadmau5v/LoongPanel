/*
 * 创建人： deadmau5v
 * 创建时间： 2024-7-4
 * 文件作用： 面板定时任务管理
 */

package Cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

var (
	Cron *cron.Cron
)

func DurationToCron(d time.Duration) string {
	return fmt.Sprintf("@every %s", d.String())
}

func init() {
	Cron = cron.New(cron.WithSeconds())
	Cron.Start()
}
