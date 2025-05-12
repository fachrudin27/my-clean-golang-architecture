package v1

import (
	"fmt"
	"my-clean-architecture-template/config"
	"my-clean-architecture-template/internal/delivery/http/v1/middleware"
	"my-clean-architecture-template/internal/model"
	"my-clean-architecture-template/internal/usecase"
	"my-clean-architecture-template/pkg/helper"
	"my-clean-architecture-template/pkg/logger"
	validators "my-clean-architecture-template/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authRoutes struct {
	t   usecase.Auth
	cfg *config.Config
	l   logger.Interface
	v   *validator.Validate
}

func NewTranslationRoutes(cfg *config.Config, handler *gin.RouterGroup, t usecase.Auth, l logger.Interface, v *validator.Validate) {
	r := &authRoutes{t, cfg, l, v}

	h := handler.Group("/auth")
	{
		h.POST("/login", r.login)
	}

	w := handler.Group("/users")
	{
		w.Use(middleware.VerifyJwtToken(cfg.JWT.SecretKey))
		// w.GET("/", r.Users)
	}
}

func (r *authRoutes) login(c *gin.Context) {

	// go func() {
	// 	// Cek apakah ada goroutine yang tidak berhenti
	// 	log.Println("Goroutine handler jalan...")

	// 	// Simulasi proses lama
	// 	time.Sleep(10 * time.Second)
	// 	log.Println("Goroutine handler selesai...")
	// }()
	// <-c.Done()

	req := model.LoginUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		customError := NewErrors([]error{err})
		c.JSON(http.StatusBadRequest, model.WebResponse[*model.UserResponse]{Errors: *customError})
		return
	}

	helper.AttemptMu.Lock()
	if attempt, exist := helper.LoginAttempt[req.Username]; exist && attempt >= 5 {
		helper.AttemptMu.Unlock()
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many login attempts,please try again later",
		})
		return
	}
	helper.LoginAttempt[req.Username]++
	helper.AttemptMu.Unlock()

	err := r.v.Struct(req)
	if err != nil {
		customError := validators.CustomError(err)
		c.JSON(http.StatusBadRequest, model.WebResponse[*model.UserResponse]{Errors: customError})
		return
	}

	data, code, errs := r.t.Login(c.Request.Context(), req)
	if errs != nil {
		r.l.Error(errs, "http - v1 - login")

		customError := NewErrors(errs)

		c.JSON(code, model.WebResponse[*model.UserResponse]{Errors: *customError})

		return
	}

	c.JSON(code, model.WebResponse[*model.UserResponse]{Data: &data})
}

func (r *authRoutes) Users(c *gin.Context) {
	token, err := r.t.Users(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - users")

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	fmt.Println(token)

	c.JSON(http.StatusOK, token)
}
