package main

import (
	_ "classNews/models"
	_ "classNews/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.Run()
}

func HandlePrePage(data int) int{
	pageIndex:=data-1
	//pageIndex1:=strconv.Itoa(pageIndex)
	return pageIndex
}
func HandleNextPage(data int) int  {
	pageIndex:=data+1
	return pageIndex
}