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
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 0,
		"data":   list,
	})

}

// GetImageList 获取镜像列表 API
func GetImageList(c *gin.Context) {
	list, err := Docker.GetImageList()
	if err != nil {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 0,
		"data":   list,
	})
}

// DeleteContainer 删除容器 API
func DeleteContainer(c *gin.Context) {
	id := c.Query("id")
	err := Docker.DeleteContainer(id)
	if err != nil {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 0,
		"msg":    "删除成功",
	})
}

// DeleteImage 删除镜像 API
func DeleteImage(c *gin.Context) {
	id := c.Query("id")
	err := Docker.DeleteImage(id)
	if err != nil {
		c.JSON(200, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 0,
		"msg":    "删除成功",
	})
}
