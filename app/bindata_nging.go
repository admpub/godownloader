// +build bindata

package main

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	bindataLib "github.com/webx-top/echo/middleware/bindata"
)

func init() {
	bindata = true
	staticMW = bindataLib.Static("/public/", &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "",
	})
	renderMgr = bindataLib.NewTmplManager(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "template",
	})
	langConf.SetFSFunc(func(dir string) http.FileSystem {
		return &assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    dir,
		}
	})
}
