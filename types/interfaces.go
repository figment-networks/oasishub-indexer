package types

import (
	"context"
	"github.com/gin-gonic/gin"
)

type HttpHandler interface {
	Handle(*gin.Context)
}

type WorkerHandler interface {
	Handle()
}

type CmdHandler interface {
	Handle(context.Context)
}