package controllers

import (
	//"beeDemo/models"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	//bee run 中，控制器改变时，会直接重新build，热更新。静态文件改变时，不会重build
	c.TplName = "index.tpl"
	//c.TplName = "hello.html"
}

func (c *MainController) Post()  {
	 c.Data["data"] = "ningning"
	 c.TplName = "hello.html"
}

func (c *MainController) ShowGet(){
	//获取orm对象
	//o := orm.NewOrm()

	//执行某个操作函数，增删改查
	/*var user models.User
	user.Name = "ning"
	user.PassWord = "80"
	count,err := o.Insert(&user)
	if err != nil{
		beego.Error("插入失败")
	}
	beego.Info(count)*/

	/*var user models.User //查询对象，用来存储数据
	user.Id = 1
	err := o.Read(&user,"Id") //参数2为查询条件，如果为主键id，也可以省略不写。
	if err != nil{
		beego.Error("查询失败")
	}
	beego.Info(user)*/

	/*var user models.User
	user.Id = 2
	err := o.Read(&user)
	if err != nil{
		beego.Error("数据不存在")
	}
	user.Name = "papa"
	count,err := o.Update(&user)
	if err != nil{
		beego.Error("更新失败")
	}
	beego.Info(count)*/

	/*var user models.User
	user.Id = 2
	count,err := o.Delete(&user) //删之前也可以先查询一下，最好先查询一下
	if err != nil{
		beego.Error("删除失败")
	}
	beego.Info(count)*/

	//显示视图
	c.TplName = "hello.html"
}