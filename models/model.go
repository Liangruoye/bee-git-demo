package models

import (
	//"github.com/astaxie/beego"
	//"database/sql" //传统方法操作数据库需要这个包和下面mysql的包
	_"github.com/go-sql-driver/mysql" //这个包设置下划线，要引进来就执行init，因为在"database/sql"包里面要用它
	"github.com/astaxie/beego/orm" //ORM方法操作数据库需要这个包
	"time"
)

//models里面放表的设计。数据操作一般还是放控制器。

/**
传统方法操作数据库
*/
//func init() {
//	//连接数据库(参数1是数据库驱动，参数2是连接数据库字符串）
//	conn,err := sql.Open("mysql","root:liang@tcp(127.0.0.1:3306)/beeTest?charset=utf8")
//	if err != nil{
//		beego.Info("连接数据库错误！",err)//beego打印消息的第一种方式，蓝色字显示
//		beego.Error("连接数据库错误！",err)//beego打印错误的第二种方式，红色字显示
//		return
//	}
//	defer conn.Close()
//
//	//创建表
//	_,err2:= conn.Exec("create table user(name varchar(40),password varchar(40));")
//	if err2 != nil{
//		beego.Info("创建表失败",err2)
//		beego.Error("创建表失败",err2)
//	}
//
//	//插入数据
//	conn.Exec("insert into user(name,password) values (?,?);","papa","80")
//
//	//查询数据
//	res,err3 := conn.Query("select name from user;")
//	if err3 != nil{
//		beego.Error("查询错误！",err3)
//		return
//	}
//	var name string
//	for res.Next(){
//		res.Scan(&name)
//		beego.Info(name)
//	}
//}

/**
ORM操作数据库
*/
type User struct{
	Id int //在表中会自动转为小写id
	Name string
	PassWord string //在表中字段转为：pass_word
	Articles []*Article `orm:"reverse(many)"` //reverse反向只可以使用one或many
	//在orm里面，__(双下划线)是有特殊含义的，所以要避免使用下划线命名字段，如Pass_Word。
}
type Article struct{
	Id int `orm:"pk;auto"`
	ArtiName string `orm:"size(20)"`
	Atime time.Time `orm:"auto_now"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"size(500)"`
	Aimg string `orm:"size(100)"`
	ArticleType *ArticleType `orm:"rel(fk)"` //fk外键
	//ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"` //on_delete的设置是为了关闭级联删除
	Users []*User `orm:"rel(m2m)"` //many to many，多对多
}
type ArticleType struct{
	Id int
	TypeName string `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` //类型表和文章表的关系是一对多(many)
}
func init() {
	//获取连接对象
	orm.RegisterDataBase("default","mysql","root:liang@tcp(127.0.0.1:3306)/beeTest?charset=utf8")
	//创建表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//生成表（参数1数据库别名，参数2是否强制更新，参数3是否可见这个过程）
	orm.RunSyncdb("default",false,true)
}

