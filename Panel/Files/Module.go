package Files

type File struct {
	Name     string `json:"name"`     // 文件名
	Size     int64  `json:"size"`     // 文件大小
	Path     string `json:"path"`     // 文件路径
	User     uint32 `json:"user"`     // 所有者
	Group    uint32 `json:"group"`    // 所属组
	Mod      string `json:"mod"`      // 权限
	Time     string `json:"time"`     // 修改时间
	IsHidden bool   `json:"isHidden"` // 是否是隐藏文件
	IsDir    bool   `json:"isDir"`    // 是否是目录
	Ext      string `json:"ext"`      // 扩展名
}
