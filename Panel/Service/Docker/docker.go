/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-17
 * 文件作用：管理docker
 */

package Docker

import (
	"LoongPanel/Panel/Service/AppStore"
	"LoongPanel/Panel/Service/PanelLog"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ContainerList 获取容器列表
func ContainerList() ([]types.Container, error) {
	app := AppStore.FindApp("Docker")
	if !app.IsRunning() || !app.IsInstall() {
		return nil, errors.New("请先安装Docker并启动")
	}
	docker := GetClient()
	if docker == nil {
		return nil, errors.New("无法连接到Docker")
	}

	list, err := docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetImageList 获取镜像列表
func GetImageList() ([]image.Summary, error) {
	app := AppStore.FindApp("Docker")
	if !app.IsRunning() || !app.IsInstall() {
		return nil, errors.New("请先安装Docker并启动")
	}
	docker := GetClient()
	if docker == nil {
		return nil, errors.New("无法连接到Docker")
	}

	list, err := docker.ImageList(context.Background(), image.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetClient 获取docker客户端
func GetClient() *client.Client {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		PanelLog.ERROR("[Docker管理]", err.Error())
		return nil
	}
	return docker
}
