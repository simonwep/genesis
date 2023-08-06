package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type loginBody struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

const refreshAccessCookieName = "genesis_refresh_token"

var validate = validator.New()

func Login(c *gin.Context) {
	var body loginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	} else if err := validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if user, err := core.GetUser(body.User); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to check for user", zap.Error(err))
		return
	} else if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := core.AuthenticateUser(body.User, body.Password)
	if user == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	issueTokens(c, user)
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshAccessCookieName)

	if err != nil || len(refreshToken) == 0 {
		c.Status(http.StatusUnauthorized)
	} else if parsed, err := core.ParseAuthToken(refreshToken); err != nil || parsed == nil {
		c.Status(http.StatusUnauthorized)
	} else if err := core.StoreInvalidatedToken(parsed.ID, parsed.ExpiresAt.Sub(time.Now())); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     refreshAccessCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Now(),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		c.Status(http.StatusOK)
	}
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshAccessCookieName)

	if err != nil || len(refreshToken) == 0 {
		c.Status(http.StatusUnauthorized)
	} else if parsed, err := core.ParseAuthToken(refreshToken); err != nil || parsed == nil {
		c.Status(http.StatusUnauthorized)
	} else if user, err := core.GetUser(parsed.User); err != nil {
		c.Status(http.StatusUnauthorized)
	} else if err := core.StoreInvalidatedToken(parsed.ID, parsed.ExpiresAt.Sub(time.Now())); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		issueTokens(c, user)
	}
}

func issueTokens(c *gin.Context, user *core.User) {
	if accessToken, err := core.CreateAccessToken(user); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to create auth token", zap.Error(err))
	} else if refreshToken, err := core.CreateRefreshToken(user); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to create auth token", zap.Error(err))
	} else {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     refreshAccessCookieName,
			Value:    refreshToken,
			Path:     "/",
			Expires:  time.Now().Add(core.Config.JWTRefreshExpiration),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		c.JSON(http.StatusOK, LoginResponse{
			Token:     accessToken,
			ExpiresAt: time.Now().UnixMilli() + core.Config.JWTAccessExpiration.Milliseconds(),
		})
	}
}
