package controllers

import (
	"bytes"
	"classNews/models"
	"encoding/gob"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
	"math"
	path "path/filepath"
	"strconv"
	"time"
)

type  ArticleController struct {
	beego.Controller
}
//处理下拉框请求
func (c *ArticleController)HandleTypeSelect()  {
	//接收数据
	typeName:=c.GetString("select")
	//logs.Info(typeName)
	//处理数据
	if typeName==""{
		logs.Info("下拉框传递数据失败")
		return
	}
	//3 查询数据
	o:=orm.NewOrmUsingDB("default")
	var articles []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	logs.Info("articles:",articles)
}
//文章列表
func (c *ArticleController) ShowArticleList() {
	////查询
	//userName:=c.GetSession("userName")
	//if userName==nil{
	//	c.Redirect("/",302)
	//	return
	//}
	//接收数据
	typeName:=c.GetString("select")
	logs.Info("下拉框数据：",typeName)
	//查询
	o:=orm.NewOrmUsingDB("default")
	qs:=o.QueryTable("Article")
	//var articles[] models.Article
	//qs.All(&articles)// select * from article
	var err error
	//pageIndex:=1
	var pageIndex int
	pageIndex,err=c.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}
	//logs.Info(articles)
	//1,显示记录数
	count,err:=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	if err!=nil{
		logs.Info("查询错误")
		return
	}
	//	２．共几页
	pageSize:=3
	start:=pageSize*(pageIndex-1)
	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)//1,一页显示多少行,2 start 起始位置 relatedsel关联表
	pageCount１:=float64(count)/float64(pageSize)
	pageCount:=math.Ceil(pageCount１)//向上取正
	// pageCount=math.Floor(pageCount) //向下取正
	//	３．首页和末页
	//	４，上一页，下一页
	//首页末页数据处理
	FirstPage:=false
	if pageIndex==1{
		FirstPage=true
	}
	EndPage:=false
	if pageIndex==int(pageCount){
		EndPage=true
	}
	var types []models.ArticleType
	connRedis,err:=redis.Dial("tcp",":6379")
	rel,err1:=redis.Bytes(connRedis.Do("get","types"))
	if err1!=nil{
		logs.Info("redis get",err)
		return
	}
	dec:=gob.NewDecoder(bytes.NewReader(rel))
	dec.Decode(&types)
	if len(types)==0{
		o.QueryTable("ArticleType").All(&types)
		logs.Info("mysql get",types)
		var buffer bytes.Buffer
		enc:=gob.NewEncoder(&buffer)
		enc.Encode(&types)
		_,err=connRedis.Do("set","types",buffer.Bytes())
		if err!=nil{
			logs.Info("redis do",err)
		}
	}
	c.Data["Types"]=types
	//-----------根据类型获取数据--------------
	//处理数据
	var articleswithtype [] models.Article
	if typeName==""{
		logs.Info("下拉框传递数据失败")
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articleswithtype)
		//return
	}else{
		 qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}
	//把数据传递给视图显示
	userName:=c.GetSession("userName")
	logs.Info(userName)
	c.Data["userName"]=userName
	c.Data["count"]=count
	c.Data["pageCount"]=pageCount
    c.Data["pageIndex"]=pageIndex
	c.Data["FirstPage"]=FirstPage
	c.Data["EndPage"]=EndPage
	c.Data["articles"]=articleswithtype   //articles
	c.Data["typeName"]=typeName
	c.Layout="layout.html"
	c.TplName="index.html"
}

func (c *ArticleController) ShowAddArticle() {
   //查询类型数据
	o:=orm.NewOrmUsingDB("default")
	var articleTypes []models.ArticleType
	_,err:=o.QueryTable("ArticleType").All(&articleTypes)
	if err!=nil{
		logs.Info("查询数据错误:",err)
		return
	}
	//传递到视图
	c.Data["Types"]=articleTypes
	c.TplName="add.html"
}

func (c *ArticleController) HandleAddArticle()  {
	articleName:=c.GetString("articleName")
	content:=c.GetString("content")
	f,h,err:=c.GetFile("uploadname")
	defer f.Close()
	//1,格式判断
	ext:=path.Ext(h.Filename)
	logs.Info(ext)
	if ext!=".jpeg"&&ext!=".png"&&ext!=".jpg"{
		logs.Info("上传文件格式不正确！")
		return
	}
	//2,文件大小
	if h.Size>500000{
		logs.Info("文件超大，不允许上传！",h.Size)
		return
	}
	//3,不能重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")
	if err!=nil{
		logs.Info("上传文件失败",err)
	}
	err=c.SaveToFile("uploadname","static/img/"+fileName+ext)
	if err!=nil{
      logs.Info("保存文件失败",err,"static/img/"+fileName+ext)
	}
	logs.Info(articleName,content,fileName+ext)
	//获取一个对象
	o:=orm.NewOrmUsingDB("default")
	//创建一赋值个插入对象
	article:=models.Article{}
	//赋值
	article.Title=articleName
	article.Content=content
	article.Img="./static/img/"+fileName+ext
	//插入
	//_,err=o.Insert(&article)
	//if err!=nil{
	//	logs.Info("插入失败")
	//	return
	//}
	//给article对象赋值
	//获取下拉框的值
	typeName:=c.GetString("select")
	//类型判断
	if typeName==""{
		logs.Info("下拉框数据句错误")
		return
	}
	//获取TYPE对象
	var articleType models.ArticleType
	articleType.TypeName=typeName
	err=o.Read(&articleType,"TypeName")
	if err!=nil{
		logs.Info("获取类型错误。")
		return
	}
	article.ArticleType=&articleType
	_,err=o.Insert(&article)
	if err!=nil{
		logs.Info("插入数据失败",err)
		return
	}
	c.Redirect("/Article/ShowArticle",302)
}

//显示文件详情
func (c *ArticleController)ShowContent(){
	//获取ID
    id:=c.GetString("id")
	logs.Info("id",id)
	//获取orm对象
	o:=orm.NewOrmUsingDB("default")
	//获取查询对象
	id2,_:=strconv.Atoi(id)
	article:=models.Article{Id:id2}
	//查询
	err:=o.Read(&article)
	if err!=nil{
		logs.Info("查询数据为空",err)
		return
	}
	article.Count++
	//多对多插入读者
	//获取多对多操作对象
	m2m:=o.QueryM2M(&article,"Users")
	//获取插入对象
	userName:=c.GetSession("userName")
	user:=models.User{UserName: userName.(string)}
	o.Read(&user,"UserName")
    //多对多插入
	_,err=m2m.Add(&user)
	if err!=nil{
		logs.Info("多对多插入失败",err)
		return
	}
	o.Update(&article)
	//o.LoadRelated(&article,"Users")
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id2).Distinct().All(&users)
	//logs.Info(users)
	//传递数据给视图
	c.Data["users"]=users
	c.Data["article"]=article
	c.LayoutSections=make(map[string]string)
	c.LayoutSections["contentHead"]="head.html"
	c.Layout="layout.html"
	c.TplName="content.html"
}
//删除文章
func (c *ArticleController) HandleDeleteArticle() {
	//获取ID
	id,_:=c.GetInt("id")
	logs.Info("id",id)
	//获取orm对象
	o:=orm.NewOrmUsingDB("default")
	//获取查询对象
	article:=models.Article{Id:id}
	//删除
	_,err:=o.Delete(&article)
	if err!=nil{
		logs.Info("删除数据错误:",err)
		return
	}
	c.Redirect("/Article/ShowArticle",302)
}
func (c *ArticleController)ShowUpdateArticle()  {
	//获取ID
	id,_:=c.GetInt("id")
	logs.Info("id",id)
	//获取orm对象
	o:=orm.NewOrmUsingDB("default")
	//获取查询对象
	article:=models.Article{Id:id}
	err:=o.Read(&article)
	if err!=nil{
		logs.Info("查询数据为空:",err)
		return
	}
	//传递数据给视图
	c.Data["article"]=article
	c.TplName="update.html"
}

func (c *ArticleController)HandleUpdateArticle() {
	id,_:=c.GetInt("id")
	logs.Info("id",id)
	articleName:=c.GetString("articleName")
	content:=c.GetString("content")
	f,h,err:=c.GetFile("uploadname")

	if err!=nil{
		logs.Info("获取文件失败",err)
		return
	}
	defer f.Close()
	//1,格式判断
	ext:=path.Ext(h.Filename)
	logs.Info("文件后缀:",ext)
	if ext!=".jpeg"&&ext!=".png"&&ext!=".jpg"{
		logs.Info("上传文件格式不正确！")
		return
	}
	//2,文件大小
	if h.Size>500000{
		logs.Info("文件超大，不允许上传！",h.Size)
		return
	}
	//3,不能重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")
	err=c.SaveToFile("uploadname","static/img/"+fileName+ext)
	if err!=nil{
		logs.Info("保存文件失败",err,"static/img/"+fileName+ext)
	}
	logs.Info(articleName,content,fileName+ext)
	//获取一个对象
	o:=orm.NewOrmUsingDB("default")
	//创建一赋值个插入对象
	logs.Info("id",id)
	article:=models.Article{Id: id}
	err=o.Read(&article)
	if err!=nil{
		logs.Info("查询数据为空:",article,err)
		return
	}
	//赋值
	article.Title=articleName
	article.Content=content
	article.Img="static/img/"+fileName+ext
	o.Update(&article,"Title","Content","Img")
	c.Redirect("/Article/ShowArticle",302)
}

func (c *ArticleController)ShowArticleType()  {
	//读取类型表
	o:=orm.NewOrmUsingDB("default")
	var articleTypes []models.ArticleType
	_,err:=o.QueryTable("ArticleType").All(&articleTypes)
	if err!=nil{
		logs.Info("查询数据错误:",err)
		return
	}
	c.Data["Types"]=articleTypes
	c.TplName="addType.html"
}
func (c *ArticleController)HandleAddArticletType()  {
	//获取数据
	typename:=c.GetString("typeName")
	//判断数据
	if typename==""{
		logs.Info("添加数据类型为空")
		return
	}
	//执行插入操作
	o:=orm.NewOrmUsingDB("default")
	var articleType models.ArticleType
	articleType.TypeName=typename
	_,err:=o.Insert(&articleType)
	if err!=nil{
		logs.Info("添加数据类型为空")
		return
	}
	//展示视图
	c.Redirect("/Article/AddArticleType",302)
}
/*退出登录*/
func (c *ArticleController) Logout() {
	//删除登录状态
	c.DelSession("userName")
	//跳转登录界面
	c.Redirect("/",302)
}