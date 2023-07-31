package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

const refreshAccessCookieName = "genesis_refresh_token"

func Login(c *gin.Context) {
	var body LoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// quick path if username is invalid
	if !validateUserName(body.User) {
		c.Status(http.StatusUnauthorized)
		return
	}

	if _, err := core.GetUser(body.User); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Error("failed to check for user", zap.Error(err))
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
		core.Error("failed to create auth token", zap.Error(err))
	} else if refreshToken, err := core.CreateRefreshToken(user); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Error("failed to create auth token", zap.Error(err))
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

func validateUserName(name string) bool {
	users := core.Config.AppUsersToCreate

	for _, user := range users {
		if user.Name == name {
			return true
		}
	}

	return false
}
