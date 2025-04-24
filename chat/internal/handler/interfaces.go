package handler

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Request any
type Response any

// BasicHandler - fiber.Ctx'ye ihtiyaç duymayan handler'lar için
type BasicHandler[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, error)
}

// FiberHandler - fiber.Ctx gerektiren handler'lar için
type FiberHandler[R Request, Res Response] interface {
	Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *R) (*Res, error)
}

type FiberWSHandler[R Request] interface {
	Handle(c *websocket.Conn, ctx context.Context, req *R)
}
