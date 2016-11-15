/*
@author: foolbread
@time: 2016-11-15
@file:tohno/server/server.go
*/
package server

import (
	"github.com/foolbread/fbcommon/golog"
	"github.com/julienschmidt/httprouter"

	"net/http"

	"github.com/foolbread/tohno/config"
)

func InitServer() {
	golog.Info("tohno server initing......")

	g_router.POST("/file", handlerFile)
	g_router.POST("/dir", handlerDir)

	go startServer()
}

var g_router *httprouter.Router = httprouter.New()

func startServer() {
	http.ListenAndServe(config.GetConfig().GetServerAddr(), g_router)
}

func handlerFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func handlerDir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
