package http

import (
	"context"
	"errors"
	"medods-test-task/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const RequestTimeout = 10 * time.Second

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and generates access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user_id query string true "User GUID"
// @Success      200 {object} TokenResponse "access_token & refresh_token"
// @Failure      400 {object} ErrorResponse "User id is empty"
// @Failure      500 {object} ErrorResponse "Failed to create new session"
// @Router /auth/login [post]
func (c *AppController) Login(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	IPAddress := ctx.ClientIP()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	accessToken, refreshToken, err := c.serv.NewSession(ctxWithTimeout, userID, IPAddress)
	if err != nil {
		if errors.Is(err, models.ErrEmptyUserID) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "User id is empty."})

			return
		}

		c.logger.Error(ctx, "Failed to create new session", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create new session."})

		return
	}

	ctx.JSON(http.StatusOK, TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
