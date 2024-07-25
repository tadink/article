package tpl

import (
	"GinTest/db"
	"fmt"
	"math/rand/v2"
	"path"
	"reflect"

	"github.com/CloudyKit/jet"
)

var templateSet *jet.Set = jet.NewHTMLSet("./templates")

func init() {
	templateSet.AddGlobalFunc("randomUrl", randomUrl)
	templateSet.AddGlobalFunc("articles", articles)
}

func GetTemplate(templateDir, templateName string) (*jet.Template, error) {
	p := path.Join(templateDir, templateName)
	return templateSet.GetTemplate(p)
}
func GetIndexTemplate(templateDir string) (*jet.Template, error) {
	return GetTemplate(templateDir, "index.html")
}

func GetListTemplate(templateDir string) (*jet.Template, error) {
	return GetTemplate(templateDir, "list.html")
}

func GetDetailTemplate(templateDir string) (*jet.Template, error) {
	return GetTemplate(templateDir, "detail.html")
}
func randomUrl(args jet.Arguments) reflect.Value {
	args.RequireNumOfArguments("randomUrl", 1, 2)
	urlType := args.Get(0)
	id := args.Get(1)
	var articleId int64 = rand.Int64N(9999999999) + 10000
	if id.IsValid() && !id.IsNil() && !id.IsZero() {
		articleId = id.Int()
	}
	siteConfig := args.Runtime().Resolve("siteConfig")
	var u string
	var t string = "detail"
	f := siteConfig.MethodByName("DetailSuffix")
	if urlType.String() == "list" {
		t = urlType.String()
		f = siteConfig.MethodByName("ListSuffix")
	}
	suffix := f.Call([]reflect.Value{})
	u = fmt.Sprintf("/%s%s/%d", t, suffix[0].String(), articleId)
	return reflect.ValueOf(u)
}
func articles(args jet.Arguments) reflect.Value {
	args.RequireNumOfArguments("articles", 5, 5)
	typeId := args.Get(0)
	page := args.Get(1)
	size := args.Get(2)
	order := args.Get(3)
	direction := args.Get(4)
	typeId.CanInt()
	articles, err := db.QueryArticleList(int(typeId.Float()), int(page.Float()), int(size.Float()), order.String(), direction.String())
	if err != nil {
		panic(err)
	}
	return reflect.ValueOf(articles)
}
