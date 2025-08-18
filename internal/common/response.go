package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/common/tracing"
)

type BaseResponse struct{}

type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(ctx *gin.Context, err error, data any) {
	if err != nil {
		base.error(ctx, err)
	} else {
		base.success(ctx, data)
	}
}

func (base *BaseResponse) success(ctx *gin.Context, data any) {
	r := response{
		Errno:   0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	}

	resp, _ := json.Marshal(r)
	ctx.Set("response", string(resp))
	ctx.JSON(http.StatusOK, r)
}

func (base *BaseResponse) error(ctx *gin.Context, err error) {
	r := response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	}

	resp, _ := json.Marshal(r)
	ctx.Set("response", string(resp))
	ctx.JSON(http.StatusOK, r)
}
