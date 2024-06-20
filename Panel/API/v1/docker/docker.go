/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：docker 管理 API
 */

package docker

import (
	"LoongPanel/Panel/Service/Docker"
	"github.com/gin-gonic/gin"
)

// GetContainerList 获取容器列表 API
func GetContainerList(c *gin.Context) {
	list, err := Docker.ContainerList()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})

}

// GetImageList 获取镜像列表 API
func GetImageList(c *gin.Context) {
	list, err := Docker.GetImageList()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"data": list,
	})
}
