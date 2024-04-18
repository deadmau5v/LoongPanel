package Files

type File struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	User     string `json:"user"`
	Group    string `json:"group"`
	Mod      string `json:"mod"`
	Time     string `json:"time"`
	IsHidden bool   `json:"isHidden"`
	IsDir    bool   `json:"isDir"`
}
