/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-25
 * 文件作用：用户操作
 */

package Database

import "LoongPanel/Panel/Service/PanelLog"

func (user *User) Save() {
	// 保存到数据库
	DB.Create(&user)

}

func (user *User) Delete() {
	DB.Delete(&user)
}

func (user *User) Update() {
	// 将空字段设置为NULL
	updateData := map[string]interface{}{
		"name":  user.Name,
		"mail":  user.Mail,
		"phone": user.Phone,
		"role":  user.Role,
	}

	for key, value := range updateData {
		if value == "" {
			updateData[key] = nil
		}
	}

	DB.Model(&User{}).Where("id = ?", user.ID).Updates(updateData)
}

func UserFind() []User {
	var Users []User
	DB.Find(&Users)

	return Users
}

func init() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		PanelLog.ERROR("[数据库] 初始化数据库失败")
		return
	}
}
