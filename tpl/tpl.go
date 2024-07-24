package tpl

import (
	"fmt"
	"math/rand"
	"path"
	"reflect"

	"github.com/CloudyKit/jet"
)

var templateSet *jet.Set = jet.NewHTMLSet("./templates")

func init() {
	templateSet.AddGlobalFunc("randomUrl", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("randomUrl", 1, 2)
		urlType := args.Get(0)
		id := args.Get(1)
		var articleId int64
		if id.IsNil() || id.Int() == 0 {
			articleId = rand.Int63n(9999999999) + 10000
		} else {
			articleId = id.Int()
		}
		var u string
		switch urlType.String() {
		case "list":
			u = fmt.Sprintf("/list/%d", articleId)
		case "detail":
		default:
			u = fmt.Sprintf("/detail/%d", articleId)
		}

		return reflect.ValueOf(u)
	})

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
