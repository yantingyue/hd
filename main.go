package main

//
//import (
//	"crypto/md5"
//	"encoding/hex"
//	"fmt"
//	"github.com/spf13/cast"
//	"github.com/valyala/fasthttp"
//	"sync"
//	"time"
//)

//import (
//	"demo/hd/route"
//	"demo/hd/service"
//	"github.com/gin-gonic/gin"
//)
//
//func main() {
//	defer close(service.TaskJobCh)
//	go func() {
//		service.InitTask()
//	}()
//	g := gin.Default()
//
//	route.Router(g)
//
//	if err := g.Run(":9988"); err != nil {
//		panic(err)
//	}
//}
