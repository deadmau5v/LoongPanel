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
	PanelLog.INFO("[Docker管理]", "获取容器列表")
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

// ContainerCount 获取容器数量
func ContainerCount() int {
	list, err := ContainerList()
	if err != nil {
		return 0
	}
	return len(list)
}

// GetImageList 获取镜像列表
func GetImageList() ([]image.Summary, error) {
	PanelLog.INFO("[Docker管理]", "获取镜像列表")
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

// ImageCount 获取镜像数量
func ImageCount() int {
	list, err := GetImageList()
	if err != nil {
		return 0
	}
	return len(list)
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

// DeleteContainer 删除容器
func DeleteContainer(id string) error {
	app := AppStore.FindApp("Docker")
	if !app.IsRunning() || !app.IsInstall() {
		return errors.New("请先安装Docker并启动")
	}

	docker := GetClient()
	err := docker.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	PanelLog.INFO("[Docker管理]", "删除容器")
	return nil
}

// DeleteImage 删除镜像
func DeleteImage(id string) error {
	PanelLog.INFO("[Docker管理]", "删除镜像")
	app := AppStore.FindApp("Docker")
	if !app.IsRunning() || !app.IsInstall() {
		return errors.New("请先安装Docker并启动")
	}

	docker := GetClient()
	_, err := docker.ImageRemove(context.Background(), id, image.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}
