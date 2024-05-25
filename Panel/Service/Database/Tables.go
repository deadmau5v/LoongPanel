/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-24
 * 文件作用：定义表结构
 */

package Database

type User struct {
	ID       int    `json:"id" comment:"用户ID"`
	Name     string `json:"name" comment:"用户名"`
	Password string `json:"password" comment:"密码"`
	Role     string `json:"role" comment:"角色"`
}
