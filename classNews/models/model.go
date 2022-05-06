package models

import (
	"github.com/beego/beego/v2/client/orm"
     _"github.com/go-sql-driver/mysql"
	"time"
)
type User struct {
	Id int
	UserName string
	Passwd   string
	Articles []*Article `orm:"rel(m2m)"`            //多表间的外键关联/
}
//文章表和文章类型表，是一对多
type Article struct {
	Id      int         `orm:"pk;auto";description:"自动增长"`       //pk主键+auto自增
	Title   string      `orm:"size(20)"`                           //文章标题 长度20
	Content string      `orm:"size(500)"`                          //文章内容
	Img     string      `orm:"size(50);null"`                     //图片路径 大小50，可空
	Time    time.Time   `orm:"type(datetime);auto_now_add"`   //发布时间 时间类型 auto_now_add第一次保存时设置时间, auto_now 每次model保存时都会对时间自动更新
	Count   int         `orm:"default(0)"`                    //阅读量
	ArticleType  *ArticleType `orm:"rel(fk)"`                 //外键 类型
	Users       []*User       `orm:"reverse(many)"`          //多对多查询字段
}
type ArticleType struct {
	Id int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`          //一对多
}
func init()  {
	orm.RegisterDataBase("default","mysql","root:Sxzy029@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	// 第二个参数 是否强制更新
	orm.RunSyncdb("default",false,true)
}
