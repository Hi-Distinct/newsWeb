package controllers

import (
	"classNews/models"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type RegController struct {
	beego.Controller
}
func (c *RegController) ShowReg() {
	c.TplName="register.html"
}
func (c * RegController) HandleReg()  {
	//获取数据
 name:=c.GetString("userName")
 passwd:=c.GetString("password")
 //处理数据
 if name==""|| passwd ==""{
	 logs.Info("用户名和密码不能为空")
     c.TplName="register.html"
	 return
 }
	logs.Info(name,passwd)
 //插入数据
  o:=orm.NewOrmUsingDB("default")
  user:=models.User{UserName:name,Passwd:passwd}
  _,err:=o.Insert(&user)
  if err!=nil{
	  logs.Info("插入用户失败:",err)
	  c.TplName="register.html"   //(服务端的请求)
	  return
  }
  //c.TplName="login.html"
  //c.Ctx.WriteString("注册成功!")
  c.Redirect("/",302) //重定向（浏览器段的）
}

type  LoginController struct {
	beego.Controller
}

func (c *LoginController) ShowLogin() {
	name := c.Ctx.GetCookie("userName")
	if name!=""{
		c.Data["userName"]=name
		c.Data["check"]="checked"
     }
    c.TplName="login.html"
}
func (c *LoginController) HandleLogin()  {
	//获取数据
	name:=c.GetString("userName")
	passwd:=c.GetString("password")
	//var name string
	//c.Ctx.Input.Bind(&name,"userName")
	//处理数据
	if name==""|| passwd ==""{
		logs.Info("用户名和密码不能为空")
		c.TplName="login.html"
		return
	}
	logs.Info(name,passwd)
	//查询数据
	o:=orm.NewOrmUsingDB("default")
	//user:=models.User{UserName:name,Passwd:passwd}
    //err:=o.Read(&user,"UserName","Passwd")
	//if err!=nil{
	//	logs.Info("用户名或者密码失败",err)
	//	c.TplName="login.html"
	//	return
	//}
	user:=models.User{UserName:name}
	err:=o.Read(&user,"UserName")
	if err!=nil{
		logs.Info("用户名失败",err)
		c.TplName="login.html"
		return
	}
	if user.Passwd!=passwd{
		logs.Info("密码失败",err)
		c.TplName="login.html"
		return
	}
	//记住文件名
	check:=c.GetString("remember")
	if check=="on"{
		c.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		c.Ctx.SetCookie("userName",name,-1)
	}
	c.SetSession("userName",name)
	c.Redirect("/Article/ShowArticle",302)
	//c.Ctx.WriteString("登录成功!")
}