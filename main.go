package main

import (
	"GinTest/tpl"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

func main() {

	err := initDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	siteConfigs, err := loadConfigs()
	if err != nil {
		log.Fatalln(err.Error())
	}

	engine := gin.Default()
	engine.Use(func(ctx *gin.Context) {
		host, _, _ := net.SplitHostPort(ctx.Request.Host)
		config, ok := siteConfigs.Load(host)
		if !ok {
			ctx.Data(404, "text/html;charset=utf-8", []byte("域名错误:"+host))
			ctx.Abort()
			return
		}
		ctx.Set("siteconfig", config)
		ctx.Next()
	})

	engine.GET("/", func(ctx *gin.Context) {
		siteConfig := ctx.MustGet("siteconfig").(*SiteConfig)
		templateName := GetTemplateName(siteConfig)
		template, err := tpl.GetIndexTemplate(templateName)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
			return
		}
		vars := make(jet.VarMap)
		vars.Set("test", "test")
		err = template.Execute(ctx.Writer, vars, nil)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
		}

	})
	engine.GET("/detail:suffix/:id", func(ctx *gin.Context) {
		siteConfig := ctx.MustGet("siteconfig").(*SiteConfig)
		templateName := GetTemplateName(siteConfig)
		template, err := tpl.GetDetailTemplate(templateName)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
			return
		}
		id := ctx.Param("id")
		article, err := GetArticle(id)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte(err.Error()))
			return
		}
		vars := make(jet.VarMap)
		vars.Set("article", article)
		vars.Set("siteConfig", siteConfig)
		err = template.Execute(ctx.Writer, vars, nil)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
		}

	})
	engine.NoRoute(func(ctx *gin.Context) {
		siteConfig := ctx.MustGet("siteconfig").(*SiteConfig)
		templateName := GetTemplateName(siteConfig)
		template, err := tpl.GetDetailTemplate(templateName)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("读取模板文件错误"+err.Error()))
			return
		}
		article, err := GetArticle("")
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte(err.Error()))
			return
		}
		ctx.Status(200)
		vars := make(jet.VarMap)
		vars.Set("article", article)
		vars.Set("siteConfig", siteConfig)
		err = template.Execute(ctx.Writer, vars, nil)
		if err != nil {
			ctx.Data(200, "text/html;charset=utf-8", []byte("渲染模板错误"+err.Error()))
		}
	})

	serverStart(engine)

}
func serverStart(engine *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
func GetTemplateName(c *SiteConfig) string {
	t := "default"
	if c.TemplateName != "" {
		return c.TemplateName
	}
	return t
}
