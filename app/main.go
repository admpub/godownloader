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
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/webx-top/echo/defaults"
	mw "github.com/webx-top/echo/middleware"
	bindataLib "github.com/webx-top/echo/middleware/bindata"
	"github.com/webx-top/echo/middleware/render"
)

var bindata bool

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
	defaults.Use(mw.Log(), mw.Recover())

	// 注册静态资源文件
	if bindata {
		defaults.Use(bindataLib.Static("/public/", &assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "",
		}))
	} else {
		defaults.Use(mw.Static(&mw.StaticOptions{
			Path: "/public/",
			Root: "./public/",
		}))
	}

	renderOptions := &render.Config{
		TmplDir: `./template`,
		Engine:  `standard`,
	}
	renderOptions.ApplyTo(defaults.Default)
	// 注册模板引擎
	if bindata {
		manager := bindataLib.NewTmplManager(&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "template",
		})
		renderOptions.Renderer().SetManager(manager)
	}

	log.Println(gdownsrv.Start(port))
}
