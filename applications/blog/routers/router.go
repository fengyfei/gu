package routers

import (
	"github.com/fengyfei/gu/applications/blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
