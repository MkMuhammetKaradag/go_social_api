package handler

import (
	"context"
	"errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func HandleBasic[R Request, Res Response](handler BasicHandler[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		ctx := c.UserContext()
		res, err := handler.Handle(ctx, &req)

		if err != nil {
			zap.L().Error("Failed to handle request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(res)
	}
}

func HandleWithFiber[R Request, Res Response](handler FiberHandler[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		ctx := c.UserContext()
		res, err := handler.Handle(c, ctx, &req)

		if err != nil {
			zap.L().Error("Failed to handle request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(res)
	}
}
func HandleWithFiberWS[R Request](handler FiberWSHandler[R]) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		var req R

		ctx := context.Background()

		handler.Handle(c, ctx, &req)
	})
}

func parseRequest[R any](c *fiber.Ctx, req *R) error {
	if err := c.BodyParser(req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return err
	}

	if err := c.ParamsParser(req); err != nil {
		return err
	}

	if err := c.QueryParser(req); err != nil {
		return err
	}

	if err := c.ReqHeaderParser(req); err != nil {
		return err
	}

	return nil
}
