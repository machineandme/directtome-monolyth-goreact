package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machineandme/directtome-monolyth-goreact/pkg/repository"
	"github.com/machineandme/directtome-monolyth-goreact/pkg/server"
)

func main() {
	ginServer := server.GetServer()
	storage, quit := repository.NewKVStorage(15 * time.Minute)
	defer close(quit)
	ginServer.GET("/new", func(ctx *gin.Context) {
		data := map[string]string{
			"fromURI":          ctx.Query("from"),
			"toURL":            ctx.Query("to"),
			"redirectAfter":    ctx.Query("after"),
			"forceContentType": ctx.Query("content_type"),
		}
		storage.Set(data["fromURI"], data)
		ctx.JSON(200, gin.H{
			"status":   "ok",
			"accepted": ctx.Query("from"),
		})
	})
	ginServer.GET("/list", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
			"keys":   storage.List(),
		})
	})
	ginServer.Any("/f/:key", func(ctx *gin.Context) {
		databaseRecord := storage.Get(ctx.Param("key"))
		data := make(map[string]interface{})
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 1 {
				data[k] = v[0]
			} else {
				data[k] = v
			}
		}
		var bodyData map[string]interface{}
		err := ctx.ShouldBind(&bodyData)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "cannot read body. " + err.Error(),
			})
			return
		}
		for k, v := range bodyData {
			data[k] = v
		}
		for k, v := range ctx.Request.PostForm {
			if len(v) == 1 {
				data[k] = v[0]
			} else {
				data[k] = v
			}
		}
		r, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{
				"status": "failed to process content",
			})
			return
		}
		go http.Post(databaseRecord["toURL"], "application/json", bytes.NewBuffer(r))
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
	ginServer.POST("/dev/print", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(string(data))
	})
	server.RunServer(ginServer)
}
