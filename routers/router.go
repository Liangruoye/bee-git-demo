package routers

import (
	"beeDemo/controllers"
    "github.com/astaxie/beego/context" //这里要将go语言的"context"包，改成现在这种beego的context包
    "github.com/astaxie/beego"
)

func init() {
    //过滤器
    beego.InsertFilter("/article/*",beego.BeforeExec,Filter) //参数1为要过滤的路由，参数2为过滤器执行的时间点，参数三为自定义的过滤函数

    beego.Router("/", &controllers.MainController{}) //参数1为请求路径，参数2为对应控制器，参数3为对应方法
    //beego.Router("/", &controllers.MainController{},"get:ShowGet")

    beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandlePost")
    beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/logout",&controllers.UserController{},"get:Logout")

    beego.Router("/article/showArticleList", &controllers.ArticleController{},"get:ShowArticleList")
    beego.Router("/article/addArticle", &controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    beego.Router("/article/showArticleDetail", &controllers.ArticleController{},"get:ShowArticleDetail")
    beego.Router("/article/updateArticle", &controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
    beego.Router("/article/deleteArticle", &controllers.ArticleController{},"get:DeleteArticle")
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")

    //自定义请求方法
	//beego.Router("/login",&controllers.LoginController{},"get:ShowLogin;post:PostFunc")
	//给多种请求指定一个方法
	//beego.Router("/index",&controllers.IndexController{},"get,post:HandleFunc")
	//给所有种类请求指定一个方法
	//beego.Router("/index",&controllers.IndexController{},"*:HandleFunc")
	//范围越小，优先级越高
	//beego.Router("/index",&controllers.IndexController{},"*:HandleFunc,post:PostFunc")
}

var Filter = func(ctx * context.Context){
    userName := ctx.Input.Session("userName")
    if userName == nil{
        ctx.Redirect(302,"/login")
        return
    }
}
