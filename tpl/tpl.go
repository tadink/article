package tpl

import (
	"GinTest/db"
	"fmt"
	"math/rand/v2"
	"path"
	"reflect"
	"strconv"

	"github.com/CloudyKit/jet"
)

var templateSet *jet.Set = jet.NewHTMLSet("./templates")

func init() {
	templateSet.AddGlobalFunc("randomUrl", randomUrl)
	templateSet.AddGlobalFunc("getArticles", getArticles)
	templateSet.AddGlobalFunc("getArticle", getArticle)
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
func getArticles(args jet.Arguments) reflect.Value {
	args.RequireNumOfArguments("getArticles", 1, 1)
	size := args.Get(0)
	articles, err := db.GetArticleList(int(size.Float()))
	if err != nil {
		panic(err)
	}
	return reflect.ValueOf(articles)
}
func getArticle(args jet.Arguments) reflect.Value {
	args.RequireNumOfArguments("getArticle", 0, 1)
	p := args.Get(0)
	id := ""
	if p.IsValid() && p.CanInt() {
		id = strconv.FormatInt(p.Int(), 10)
	}
	article, err := db.GetArticle(id)
	if err != nil {
		panic(err)
	}
	return reflect.ValueOf(article)
}
