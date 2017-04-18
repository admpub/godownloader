package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/admpub/godownloader/service"
	"github.com/webx-top/echo/defaults"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
)

var (
	bindata bool

	staticMW  interface{}
	renderMgr driver.Manager
	langConf  = &language.Config{
		Default:      `en`,
		Fallback:     `zh-cn`,
		AllList:      []string{`zh-cn`, `en`},
		RulesPath:    []string{`data/i18n/rules`},
		MessagesPath: []string{`data/i18n/messages`},
		Reload:       false,
	}
)

func getSetPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".godownload")
}

func main() {
	var port int
	flag.IntVar(&port, `p`, 9981, `port`)

	gdownsrv := new(DownloadService.DServ)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		func() {
			gdownsrv.StopAllTask()
			log.Println("info: save setting ", gdownsrv.SaveSettings(getSetPath()))
		}()
		os.Exit(1)
	}()
	gdownsrv.LoadSettings(getSetPath())
	log.Printf("GUI located add http://localhost:%d/\n", port)

	defaults.SetDebug(true)
	//defaults.Use(mw.Log())
	defaults.Use(mw.Recover())
	// 注册静态资源文件
	defaults.Use(staticMW)

	// 注册模板引擎
	renderOptions := &render.Config{
		TmplDir: `./template`,
		Engine:  `standard`,
	}
	renderOptions.ApplyTo(defaults.Default)
	if renderMgr != nil {
		renderOptions.Renderer().SetManager(renderMgr)
	}
	defaults.Use(mw.FuncMap(map[string]interface{}{
		"Languages": func() []string {
			return langConf.AllList
		},
	}))
	defaults.Use(language.New(langConf).Middleware())
	log.Println(gdownsrv.Start(port))
}
