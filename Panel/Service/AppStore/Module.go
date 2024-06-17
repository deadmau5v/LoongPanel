/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

package AppStore

type App struct {
	Name    string                 // 名称
	Version func() (string, error) // 版本
	Icon    string                 // 图标
	Tags    []string               // 标签

	IsInstall func() bool          // 是否安装
	IsRunning func() bool          // 是否运行
	Install   func() (bool, error) // 安装
	Uninstall func() (bool, error) // 卸载

	Path   string               // 安装路径
	Start  func() (bool, error) // 启动
	Stop   func() (bool, error) // 停止
	Status func() (bool, error) // 状态

}
