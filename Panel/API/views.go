package API

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Error404(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Panel/Front/pages/404.html", WORKDIR+"/Panel/Front/layout.html", WORKDIR+"/Panel/Front/import.html", WORKDIR+"/Panel/Front/aside.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data := map[string]interface{}{
		"title": AppName + " - 404",
	}
	err = html.Execute(ctx.Writer, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func home(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Panel/Front/pages/index.html", WORKDIR+"/Panel/Front/layout.html", WORKDIR+"/Panel/Front/import.html", WORKDIR+"/Panel/Front/aside.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data := map[string]interface{}{
		"title":    AppName + " - 首页",
		"pageHome": true,
	}
	err = html.Execute(ctx.Writer, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func files(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Panel/Front/pages/files.html", WORKDIR+"/Panel/Front/layout.html", WORKDIR+"/Panel/Front/import.html", WORKDIR+"/Panel/Front/aside.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path := ctx.Query("path")
	if path == "" {
		path = "/"
	}
	data := map[string]interface{}{
		"title":     AppName + " - 文件管理",
		"pageFiles": true,
		"path":      path,
	}
	err = html.Execute(ctx.Writer, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func terminal(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Panel/Front/pages/terminal.html", WORKDIR+"/Panel/Front/layout.html", WORKDIR+"/Panel/Front/import.html", WORKDIR+"/Panel/Front/aside.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data := map[string]interface{}{
		"title":        AppName + " - 首页",
		"pageTerminal": true,
	}
	err = html.Execute(ctx.Writer, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func appInstall(ctx *gin.Context) {

}
