package controllers

import (
	"dksv-v2/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Post() {
	beego.Info("固定路由的get类型的方法 ")
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		//Data:   "存活",
	}
	c.Data["json"] = data
	c.ServeJSON()
}
