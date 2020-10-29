package main

import (
	_ "beeDemo/routers" //下划线是指导入这个包时，执行这个包里面所有的init()函数
	_ "beeDemo/models" //自动执行models包里的所有init函数，连接数据库
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("nextpage",ShowNextPage)
	beego.Run()
}

func ShowPrePage(pageIndex int)int{
	if pageIndex == 1 {
		return pageIndex
	}
	return pageIndex-1
}
func ShowNextPage(pageIndex,pageCount int)int{
	if pageIndex == pageCount {
		return pageCount
	}
	return pageIndex+1
}

