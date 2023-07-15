package handlers

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"demo/hd/dto"
	"demo/hd/errcode"
	"demo/hd/service"
)

func (*HdHandlerService) GrabOrder(c *gin.Context) {
	req := dto.GradOrderReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, errcode.ParseParameterFailed)
		return
	}
	if req.ProductId == 0 {
		service.HdSrv.QueryProductInfoByName(&req)
	}

	if req.ProductId == 0 {
		c.JSON(http.StatusOK, errcode.ProductNotFound)
		return
	}

	// 判断是否在任务中
	if _, ok := service.TaskIdToTaskJobMap.Load(req.ProductId); ok { // 在任务中，直接返回
		c.JSON(http.StatusOK, errcode.Success)
		return
	}

	// 创建抢购任务
	ctx, _ := context.WithCancel(context.Background())
	service.TaskJobCh <- service.TaskJob{
		Ctx:              ctx,
		QuitCh:           make(chan struct{}, 1),
		TaskId:           atomic.AddUint64(&service.TaskId, 1),
		Num:              req.Num,
		ProductId:        req.ProductId,
		NftProductSizeId: req.NftProductSizeId,
		ProductName:      req.ProductName,
		Price:            req.Price,
		AutoPay:          req.AutoPay,
	}

	c.JSON(http.StatusOK, errcode.Success)
}

func (*HdHandlerService) StopGrabOrder(c *gin.Context) {
	req := dto.StopTaskReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, errcode.ParseParameterFailed)
		return
	}

	task, ok := service.TaskIdToTaskJobMap.Load(req.ProductId)
	if !ok {
		c.JSON(http.StatusOK, errcode.TaskNotExist)
		return
	}

	task.(service.TaskJob).QuitCh <- struct{}{}
	c.JSON(http.StatusOK, errcode.Success)
}

func (*HdHandlerService) GrabOrderTaskList(c *gin.Context) {
	list := make([]dto.GrabOrderTaskInfo, 0)
	service.TaskIdToTaskJobMap.Range(func(key, value any) bool {
		task := value.(*service.TaskJob)
		list = append(list, dto.GrabOrderTaskInfo{
			ProductId:        task.ProductId,
			NftProductSizeId: task.NftProductSizeId,
			ProductName:      task.ProductName,
			Price:            task.Price,
			Num:              task.Num,
		})
		return true
	})

	c.JSON(http.StatusOK, dto.Resp{
		Code: errcode.Success.Code,
		Msg:  errcode.Success.Msg,
		Data: list,
	})
}
