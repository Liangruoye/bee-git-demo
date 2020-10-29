package controllers

import (
	"beeDemo/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path"
	"time"
	"math"
)

type ArticleController struct {
	beego.Controller
}

func(this *ArticleController)ShowArticleList(){
	//session判断登录状态
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//_,err := qs.All(&articles) //获取表中所有数据
	//if err != nil {
	//	beego.Error("查询数据错误！")
	//}

	typeName := this.GetString("select")
	var count int64
	pageSize := 2
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}
	start := (pageIndex-1)*pageSize
	if typeName == ""{
		count,_ = qs.Count()
	} else {
		count,_ = qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		//注意此处是两个下划线，字段__关联表__关联表的字段
	}
	pageCount := int(math.Ceil(float64(count)/float64(pageSize)))

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	if typeName == ""{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles) //此处不加.RelatedSel("ArticleType")是为了让没有分类的文章也显示
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}

	this.Data["userName"] = userName
	this.Data["typeName"] = typeName
	this.Data["types"] = types
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

func(this *ArticleController)ShowAddArticle(){
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["types"] = types

	this.Layout = "layout.html"
	this.TplName = "add.html"
}

func(this *ArticleController)HandleAddArticle(){
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	if articleName == "" || content == ""{
		this.Data["errmsg"] = "添加数据不完整！"
		this.TplName = "add.html"
		return
	}
	//beego.Info(articleName,content)

	filePath := UploadFile(&this.Controller,"uploadname")

	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = filePath
	//获取文章类型数据
	typeName := this.GetString("select")
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")
	article.ArticleType = &articleType

	o.Insert(&article)

	this.Redirect("/article/showArticleList",302)
}

func(this *ArticleController)ShowArticleDetail(){
	id,err := this.GetInt("articleId")
	if err != nil{
		beego.Error("连接错误！")
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	//o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).One(&article)
	//注意这里的filter和上面的filter的使用区别，因为这里是article表的id，所以filter不需要指定表

	article.Acount += 1
	o.Update(&article)

	//多对多插入浏览记录
	m2m := o.QueryM2M(&article,"Users")
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}
	var user models.User
	user.Name = userName.(string) //类型断言
	o.Read(&user,"Name")
	m2m.Add(user)

	//多对多查询（两种方法）
	//o.LoadRelated(&article,"Users") //方法一，前提是前面已经绑定了这种多对多的关系：o.QueryM2M(&article,"Users")，直接加载
	//方法二：
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users) //注意此处是两个下划线

	this.Data["users"] = users
	this.Data["article"] = article

	this.Layout = "layout.html"
	this.TplName = "content.html"
}

func(this *ArticleController)ShowUpdateArticle(){
	id,err := this.GetInt("articleId")
	if err != nil{
		beego.Error("请求文章错误！")
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	this.Data["article"] = article

	this.Layout = "layout.html"
	this.TplName = "update.html"
}

//封装上传文件函数
func UploadFile(this *beego.Controller,filePath string) string {
	file,head,err := this.GetFile(filePath)
	if head.Filename == "" {
		return "NoImg"
	}
	if err != nil{
		this.Data["errmsg"] = "文件获取失败！"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close() //file是文件流，这一句要放在判断err之后，因为有err，就没有file了。
	if head.Size > 5000000{
		this.Data["errmsg"] = "文件太大！"
		this.TplName = "add.html"
		return ""
	}
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误！"
		this.TplName = "add.html"
		return ""
	}
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext //传说是go语言诞生时间，固定数字，不能填错
	err2 := this.SaveToFile(filePath,"./static/upload/"+fileName)
	if err2 != nil{
		beego.Error(err2)
		this.Data["errmsg"] = "文件上传失败！"
		this.TplName = "add.html"
		return ""
	}
	return "/static/upload/"+fileName
} //这里this的类型填的是父类

func(this *ArticleController)HandleUpdateArticle(){
	id,err := this.GetInt("articleId")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filePath := UploadFile(&this.Controller,"uploadname")
	if err != nil || filePath == "" || articleName == "" || content == ""{
		beego.Error("请求错误！")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil{
		beego.Error("更新的文章不存在！")
		return
	}

	article.ArtiName = articleName
	article.Acontent = content
	if filePath != "NoImg"{
		article.Aimg = filePath
	}
	o.Update(&article)

	this.Redirect("/article/showArticleList",302)
}

func(this *ArticleController)DeleteArticle(){
	id,err := this.GetInt("articleId")
	if err != nil{
		beego.Error("删除文章请求错误！")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)

	this.Redirect("/article/showArticleList",302)
}

func(this *ArticleController)ShowAddType(){
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["types"] = types

	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

func(this *ArticleController)HandleAddType(){
	typename := this.GetString("typeName")
	if typename == ""{
		beego.Error("信息不完整！")
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typename
	o.Insert(&articleType)

	this.Redirect("/article/addType",302)
}

func(this *ArticleController)DeleteType(){
	id,err := this.GetInt("id")
	if err != nil{
		beego.Error("删除类型错误！",err)
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	//beego默认会级联删除，即删除类型之后，与之外键关联的文章也都会删除。若要防止级联删除，则建表时，要使用on_delete属性

	this.Redirect("/article?addType",302)
}
