package user

import (
	"context"
	"socialmedia/user/app/user/usecase"

	"github.com/gofiber/fiber/v2"
)

type UpdateBannerHandler struct {
	usecase usecase.UpdateBannerUseCase
}
type UpdateBannerRequest struct {

	BannerURL string `json:"banner_url"`

}

type UpdateBannerResponse struct {
	Message string
}

func NewUpdateBannerHandler(usecase usecase.UpdateBannerUseCase) *UpdateBannerHandler {
	return &UpdateBannerHandler{usecase: usecase}
}

func (h *UpdateBannerHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *UpdateBannerRequest) (*UpdateBannerResponse, error) {

	err := h.usecase.Execute(fbrCtx, ctx, req.BannerURL)
	if err != nil {
		return nil, err
	}
	return &UpdateBannerResponse{Message: "user updated"}, nil
}
