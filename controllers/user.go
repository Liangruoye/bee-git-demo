package controllers

import(
	"beeDemo/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func(c *UserController)ShowRegister(){
	c.TplName = "register.html"
}

func(c *UserController)HandlePost(){
	//1.获取数据
	userName := c.GetString("userName")
	pwd := c.GetString("passWord")
	//beego.Info(userName,pwd)

	//2.校验数据
	if userName == "" || pwd == ""{
		beego.Info("数据不完整！")
		c.Data["errmsg"] = "数据不完整！请重新注册。"
		c.TplName = "register.html"
		return
	}

	//3.处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	beego.Info(user)
	count,err := o.Insert(&user)
	if err != nil{
		beego.Error("注册失败！")
	}
	beego.Info(count)

	//4.返回结果
	//c.Ctx.WriteString("注册成功！")
	//c.TplName = "register.html"
	c.Redirect("/login",302)
}

func(c *UserController)ShowLogin(){
	userName := c.Ctx.GetCookie("userName")
	if userName == ""{
		c.Data["userName"] = ""
		c.Data["checked"] = ""
	}else{
		c.Data["userName"] = userName
		c.Data["checked"] = "checked"
	}
	c.TplName = "login.html"
}

func(c *UserController)HandleLogin(){
	userName := c.GetString("userName")
	pwd := c.GetString("passWord")

	if userName == "" || pwd == ""{
		c.Data["errmsg"] = "登录数据不完整！"
		c.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user,"Name")
	if err != nil{
		c.Data["errmsg"] = "用户不存在！"
		c.TplName = "login.html"
		return
	}
	if user.PassWord != pwd{
		c.Data["errmsg"] = "密码错误！"
		c.TplName = "login.html"
		return
	}

	//设置cookie，实现记住用户名功能
	data := c.GetString("remember")
	if data == "on" {
		c.Ctx.SetCookie("userName",userName,100) //100秒
	} else {
		c.Ctx.SetCookie("userName",userName,-1) //删除cookie
	}

	//设置session，实现记住登录状态功能
	c.SetSession("userName",userName)

	//c.Ctx.WriteString("登录成功！")
	c.Redirect("/article/showArticleList",302)
}

func(c *UserController)Logout(){
	//删除登录状态
	c.DelSession("userName")

	c.Redirect("/login",302)
}
