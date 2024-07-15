/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供文件管理相关的API 主要实现在 Service/Files 中
 */

package files

import (
	FileService "LoongPanel/Panel/Service/Files"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func FileDir(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		path = "/"
	}

	f, err := os.Stat(path)
	if err != nil {
		ctx.JSON(200, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
		return
	} else if !f.IsDir() {
		ctx.JSON(200, gin.H{
			"msg":    "不是一个目录",
			"status": 1,
		})
		return
	}

	files, err := FileService.Dir(path)
	if err != nil {
		ctx.JSON(200, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"files":  files,
			"status": 0,
			"msg":    "ok",
		})
	}
}

func FileRead(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":    "path is empty",
			"status": 1,
		})
		return
	}

	content, err := FileService.Content(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data":   content,
			"status": 0,
			"msg":    "ok",
		})
	}
}

// Upload 上传文件
func Upload(ctx *gin.Context) {
	// 从请求中获取文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "status": 1})
		return
	}

	// 从请求中获取路径
	path := ctx.Query("path")
	if path == "" {
		path = "/"
	}

	// 定义文件保存的路径
	dst := filepath.Join(path, file.Filename)

	// 将文件保存到服务器
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error(), "status": 1})
		return
	}

	// 返回响应给客户端
	ctx.JSON(http.StatusOK, gin.H{"msg": "文件上传成功", "status": 0})
}

// Download 下载文件
func Download(ctx *gin.Context) {
	// 从请求中获取文件路径
	path := ctx.Query("path")

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "文件不存在", "status": 1})
		return
	}

	_, filename := filepath.Split(path)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 将文件发送给客户端
	ctx.File(path)
}

// Copy 复制文件
func Copy(ctx *gin.Context) {
	path := ctx.Query("path")
	dest := ctx.Query("dest")
	err := FileService.Copy(path, dest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "文件复制成功",
			"status": 0,
		})
	}

}

// Delete 删除文件
func Delete(ctx *gin.Context) {
	path := ctx.Query("path")
	err := FileService.Delete(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "文件删除成功",
			"status": 0,
		})
	}
}

// Move 移动文件
func Move(ctx *gin.Context) {
	path := ctx.Query("path")
	dest := ctx.Query("dest")
	err := FileService.Move(path, dest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "移动文件成功",
			"status": 0,
		})
	}
}

// Rename 重命名文件
func Rename(ctx *gin.Context) {
	path := ctx.Query("path")
	newName := ctx.Query("name")
	err := FileService.Rename(path, newName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "重命名成功",
			"status": 0,
		})
	}
}

// Decompress 解压
func Decompress(ctx *gin.Context) {
	path := ctx.Query("path")
	err := FileService.Decompress(path)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "解压成功",
			"status": 0,
		})
	}
}

// Compress 压缩
func Compress(ctx *gin.Context) {
	path := ctx.Query("path")
	err := FileService.Compress(path)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    err.Error(),
			"status": 1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "压缩成功",
			"status": 0,
		})
	}
}
