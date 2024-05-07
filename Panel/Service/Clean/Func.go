/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：清理系统垃圾（apt、yum）包的缓存
 */

package Clean

import "os/exec"

// apt

func AptAutoClean() ([]byte, error) {
	output, err := exec.Command("apt", "autoclean").Output()
	if err != nil {
		return nil, err
	}
	return output, err
}

func AptAutoRemove() ([]byte, error) {
	output, err := exec.Command("apt", "autoremove").Output()
	if err != nil {
		return nil, err
	}
	return output, err
}

// yum

func YumAutoClean() ([]byte, error) {
	output, err := exec.Command("yum", "clean", "all").Output()
	if err != nil {
		return nil, err
	}
	return output, err
}
