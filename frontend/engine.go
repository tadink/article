package frontend

import (
	"GinTest/db"
	"GinTest/tpl"
	"net"
	"sync"

	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

func initEngine(configs *sync.Map) *gin.Engine {
	engine := gin.Default()
	engine.Use(SiteConfigMiddleware(configs))
	engine.GET("/", indexHandle)
	engine.GET("/list:suffix/:id", listHandle)
	engine.GET("/detail:suffix/:id", detailHandle)
	engine.NoRoute(noRoute)
	return engine
}

func indexHandle(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	template, err := tpl.GetIndexTemplate(templateName)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
		return
	}
	vars := make(jet.VarMap)
	vars.Set("siteConfig", s)
	err = template.Execute(c.Writer, vars, nil)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
	}

}

func listHandle(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	template, err := tpl.GetListTemplate(templateName)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
		return
	}
	articles, err := db.GetArticleList(20)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("获取文章列表错误"+err.Error()))
		return
	}
	vars := make(jet.VarMap)
	vars.Set("siteConfig", s)
	vars.Set("articles", articles)
	err = template.Execute(c.Writer, vars, nil)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
	}

}

func detailHandle(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	template, err := tpl.GetDetailTemplate(templateName)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
		return
	}
	id := c.Param("id")
	article, err := db.GetArticle(id)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte(err.Error()))
		return
	}
	vars := make(jet.VarMap)
	vars.Set("article", article)
	vars.Set("siteConfig", s)
	err = template.Execute(c.Writer, vars, nil)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
	}

}

func noRoute(c *gin.Context) {
	s := c.MustGet("siteconfig").(*db.SiteConfig)
	templateName := s.GetTemplateName()
	template, err := tpl.GetDetailTemplate(templateName)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
		return
	}
	article, err := db.GetArticle("")
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte(err.Error()))
		return
	}
	c.Status(200)
	vars := make(jet.VarMap)
	vars.Set("article", article)
	vars.Set("siteConfig", s)
	err = template.Execute(c.Writer, vars, nil)
	if err != nil {
		c.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
	}
}

func SiteConfigMiddleware(configs *sync.Map) func(*gin.Context) {
	return func(c *gin.Context) {
		host, _, _ := net.SplitHostPort(c.Request.Host)
		config, ok := configs.Load(host)
		if !ok {
			c.Data(404, "text/html;charset=utf-8", []byte("域名错误:"+host))
			c.Abort()
			return
		}
		c.Set("siteconfig", config)
		c.Next()
	}

}
