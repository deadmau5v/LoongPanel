package API

import (
	"github.com/gin-gonic/gin"
	"os"
)

var App *gin.Engine
var WORKDIR string

func init() {
	var err error
	WORKDIR, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	App = gin.Default()
	initRoute(App)
}
