package errcode

import "demo/hd/dto"

var (
	Success = dto.Resp{
		Code: 0,
		Msg:  "success",
	}
	ParseParameterFailed = dto.Resp{
		Code: 400,
		Msg:  "接受请求参数失败",
	}
	TaskNotExist = dto.Resp{
		Code: 401,
		Msg:  "任务不存在",
	}
	ProductNotFound = dto.Resp{
		Code: 402,
		Msg:  "未查询到藏品",
	}
)
