package API

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Error404(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Web/pages/404.html", WORKDIR+"/Web/layout.html", WORKDIR+"/Web/import.html", WORKDIR+"/Web/aside.html")
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

func Home(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR+"/Web/pages/index.html", WORKDIR+"/Web/layout.html", WORKDIR+"/Web/import.html", WORKDIR+"/Web/aside.html")
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
