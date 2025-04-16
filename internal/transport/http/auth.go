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
// @Param user_id query string true "User GUID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"
// @Success      200 {object} TokenResponse "access_token & refresh_token"
// @Failure      400 {object} ErrorResponse "User id is empty"
// @Failure      400 {object} ErrorResponse "User id must be a valid UUID"
// @Failure      500 {object} ErrorResponse "An unexpected error occurred"
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

		if errors.Is(err, models.ErrInvalidUserID) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "User id must be a valid UUID."})

			return
		}

		c.logger.Error(ctx, "Failed to create new session", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "An unexpected error occurred."})

		return
	}

	ctx.JSON(http.StatusOK, TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

// RefreshToken godoc
// @Summary      RefreshToken
// @Description  Refreshes token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param token body RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} TokenResponse "access_token & refresh_token"
// @Failure      400 {object} ErrorResponse "Invalid or missing refresh token"
// @Failure      400 {object} ErrorResponse "Token is invalid"
// @Failure      401 {object} ErrorResponse ""
// @Failure      403 {object} ErrorResponse "Session is invalid"
// @Failure      500 {object} ErrorResponse "An unexpected error occurred"
// @Router /auth/refresh [post]
func (c *AppController) RefreshToken(ctx *gin.Context) {
	var refreshTokenRequest RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil || refreshTokenRequest.RefreshToken == "" {
		c.logger.Error(ctx, "Refresh token is missing or invalid in the request body.", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid or missing refresh token."})

		return
	}
	IPAddress := ctx.ClientIP()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	accessToken, refreshToken, err := c.serv.RefreshToken(ctxWithTimeout, refreshTokenRequest.RefreshToken, IPAddress)
	if err != nil {
		if errors.Is(err, models.ErrTokenExpired) ||
			errors.Is(err, models.ErrMismatchedHashAndToken) || errors.Is(err, models.ErrSessionNotFound) {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})

			return
		}

		if errors.Is(err, models.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Token is invalid."})

			return
		}

		if errors.Is(err, models.ErrInvalidSession) {
			ctx.JSON(http.StatusForbidden, ErrorResponse{Error: "Session is invalid."})

			return
		}

		c.logger.Error(ctx, "Failed to refresh token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "An unexpected error occurred."})

		return
	}

	ctx.JSON(http.StatusOK, TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
