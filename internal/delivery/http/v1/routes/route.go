package routes

import (
	"my-clean-architecture-template/config"
	v1 "my-clean-architecture-template/internal/delivery/http/v1"
	"my-clean-architecture-template/internal/usecase"
	"my-clean-architecture-template/pkg/logger"
	"my-clean-architecture-template/pkg/validator"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, handler *gin.Engine, t usecase.Auth, l logger.Interface) {

	h := handler.Group("/v1")
	{
		v1.NewTranslationRoutes(cfg, h, t, l, validator.NewValidator())
	}
}
