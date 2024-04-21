package API

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Error404(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR + "/Web/404.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = html.Execute(ctx.Writer, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func Home(ctx *gin.Context) {
	html, err := template.ParseFiles(WORKDIR + "/Web/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = html.Execute(ctx.Writer, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
