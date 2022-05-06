package routers

import (
	"classNews/controllers"
	"github.com/beego/beego/v2/server/web/context"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	beego.InsertFilter("/Article/*",beego.BeforeRouter,filterFunc)
    //beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{},"get:ShowReg;post:HandleReg")
	beego.Router("/", &controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle",&controllers.ArticleController{},"get:ShowArticleList;post:HandleTypeSelect")
	beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/ArticleContent",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:HandleDeleteArticle")
    beego.Router("/Article/UpdateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowArticleType;post:HandleAddArticletType")
	//退出登录
	beego.Router("/Article/logout",&controllers.ArticleController{},"get:Logout")
}
var filterFunc= func(ctx *context.Context) {
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/")
	}
}
