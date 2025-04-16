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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and generates access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user_id query string true "User GUID"
// @Success      200 {object} map[string]string "access_token & refresh_token"
// @Failure      400 {object} map[string]string "User id is empty"
// @Failure      500 {object} map[string]string "Failed to create new session"
// @Router /auth/login [post]
func (c *AppController) Login(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	IPAddress := ctx.ClientIP()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	accessToken, refreshToken, err := c.serv.NewSession(ctxWithTimeout, userID, IPAddress)
	if err != nil {
		if errors.Is(err, models.ErrEmptyUserID) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User id is empty."})

			return
		}

		c.logger.Error(ctx, "Failed to create new session", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new session."})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
