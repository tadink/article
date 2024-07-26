package frontend

import (
	"GinTest/db"
	"GinTest/tpl"
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net"
	"os"
	"path"
	"sync"

	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

const HTML_CONTENT_TYPE = "text/html;charset=utf-8"

func initEngine(configs *sync.Map) *gin.Engine {
	engine := gin.Default()
	engine.Use(CacheMiddleware())
	engine.Use(SiteConfigMiddleware(configs))
	engine.GET("/favicon.ico", faviconHandle)
	engine.GET("/", indexHandle)
	engine.GET("/list:suffix/:id", listHandle)
	engine.GET("/detail:suffix/:id", detailHandle)
	engine.NoRoute(noRoute)
	return engine
}

func indexHandle(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	vars := c.MustGet("vars").(jet.VarMap)
	respond(c, templateName, "index.html", vars)
}

func listHandle(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	articles, err := db.GetArticleList(20)
	if err != nil {
		c.Data(200, HTML_CONTENT_TYPE, []byte("获取文章列表错误"+err.Error()))
		return
	}
	vars := c.MustGet("vars").(jet.VarMap)
	vars.Set("articles", articles)
	respond(c, templateName, "list.html", vars)
}

func detailHandle(c *gin.Context) {

	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	id := c.Param("id")
	article, err := db.GetArticle(id)
	if err != nil {
		c.Data(200, HTML_CONTENT_TYPE, []byte(err.Error()))
		return
	}
	vars := c.MustGet("vars").(jet.VarMap)
	vars.Set("article", article)
	respond(c, templateName, "detail.html", vars)

}

func faviconHandle(c *gin.Context) {
	data, _ := os.ReadFile("luxisr.ttf")
	fon, _ := freetype.ParseFont(data)
	fg, bg := image.Black, image.Transparent
	fontSize := 25.0
	rgba := image.NewRGBA(image.Rect(0, 0, 50, 50))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	f := freetype.NewContext()
	//f.SetDPI(320.0)
	f.SetFont(fon)
	f.SetFontSize(fontSize)
	f.SetClip(rgba.Bounds())
	f.SetDst(rgba)
	f.SetSrc(fg)
	f.SetHinting(font.HintingNone)
	// Draw the text.
	pt := freetype.Pt(int(f.PointToFixed(fontSize)>>6)-8, int(f.PointToFixed(fontSize)>>6)+8)
	_, err := f.DrawString("H", pt)
	if err != nil {
		panic(err)
	}
	err = png.Encode(c.Writer, rgba)
	if err != nil {
		panic(err)
	}
}

func noRoute(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	article, err := db.GetArticle("")
	if err != nil {
		c.Data(200, HTML_CONTENT_TYPE, []byte(err.Error()))
		return
	}
	vars := make(jet.VarMap)
	vars.Set("article", article)
	vars.Set("siteConfig", s)
	respond(c, templateName, "detail.html", vars)
}

func SiteConfigMiddleware(configs *sync.Map) func(*gin.Context) {
	return func(c *gin.Context) {
		host, _, _ := net.SplitHostPort(c.Request.Host)
		config, ok := configs.Load(host)
		if !ok {
			c.Data(404, HTML_CONTENT_TYPE, []byte("域名错误:"+host))
			c.Abort()
			return
		}
		vars := make(jet.VarMap)
		vars.Set("siteConfig", config)
		c.Set("siteconfig", config)
		c.Set("vars", vars)
		c.Next()
	}

}

func CacheMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		data, _ := getCache(c.Request.Host, c.Request.URL.String())

		if data != nil {
			c.Data(200, HTML_CONTENT_TYPE, data)
			c.Abort()
			return
		}
		c.Next()
	}

}

func setCache(host, url string, data []byte) {
	name := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	cachePath := fmt.Sprintf("./cache/%s/%s/%s", host, name[:2], name)
	dir := path.Dir(cachePath)
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err = os.WriteFile(cachePath, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func getCache(host, url string) ([]byte, error) {
	name := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	cachePath := fmt.Sprintf("./cache/%s/%s/%s", host, name[:2], name)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func handleTemplate(templateDir string, templateName string, vars jet.VarMap) ([]byte, error) {
	template, err := tpl.GetTemplate(templateDir, templateName)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = template.Execute(&buf, vars, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func respond(c *gin.Context, templateDir string, templateName string, vars jet.VarMap) {
	data, err := handleTemplate(templateDir, templateName, vars)
	if err != nil {
		c.Data(200, HTML_CONTENT_TYPE, []byte("handleTemplate error"+err.Error()))
		return
	}
	setCache(c.Request.Host, c.Request.URL.String(), data)
	c.Data(200, HTML_CONTENT_TYPE, data)
}
