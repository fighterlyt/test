package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()
	engine.UseRawPath = true
	engine.GET("/test", func(ctx *gin.Context) {
		spew.Dump(*ctx.Request)
		ctx.String(http.StatusOK, ctx.Request.URL.User.String())
	})
	engine.Run(":1235")
}
