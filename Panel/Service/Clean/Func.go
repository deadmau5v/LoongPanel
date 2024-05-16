/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：清理系统垃圾（apt、yum）包的缓存
 */

package Clean

import (
	"LoongPanel/Panel/Service/Log"
	"fmt"
	"os/exec"
	"regexp"
)

// apt

func AptAutoClean() (string, error) {
	_, err := exec.Command("apt", "autoclean", "-y").Output()
	if err != nil {
		return "autoClean失败", err
	}
	return "完成", err
}

func AptAutoRemove() (string, error) {
	output, err := exec.Command("apt", "autoremove", "-y").Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\d+ to remove`)
	res := re.FindAllString(string(output), -1)

	return fmt.Sprintf("，已删除 %s 个无用包。", res[0]), err
}

// yum

func YumAutoClean() (string, error) {
	output, err := exec.Command("yum", "clean", "all", "-y").Output()
	if err != nil {
		return "clean失败，", err
	}

	re := regexp.MustCompile(`\d+ files removed`)
	res := re.FindAllString(string(output), -1)

	return fmt.Sprintf("清理了%s个包", res), err
}

func YumAutoRemove() (string, error) {
	_, err := exec.Command("yum", "autoremove", "-y").Output()
	if err != nil {
		return "，autoRemove失败。", err
	}
	return "，完成!", err
}

// RemoveTmpDir 清理临时目录
func RemoveTmpDir() {
	TmpDirs := []string{
		"/tmp/*",
		"/var/tmp/*",
		"/var/cache/*",
		"/var/log/*",
		"/root/.cache/*",
		"/root/.local/share/Trash/*",
	}

	for _, dir := range TmpDirs {
		_, err := exec.Command("rm", "-rf", dir).Output()
		if err != nil {
			Log.ERROR("RemoveTmpDir() Error: ", err.Error())
		}
	}
}
