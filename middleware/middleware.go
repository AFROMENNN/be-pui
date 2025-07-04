package middleware

import (
	"be-pui/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil) *AuthMiddleware {
	return &AuthMiddleware{jwtUtil: jwtUtil}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token otorisasi tidak ditemukan."})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Format token tidak valid. Harusnya 'Bearer <token>'."})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims, err := m.jwtUtil.ParseJWTToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token tidak valid atau telah kedaluwarsa."})
			c.Abort()
			return
		}

		utils.SetUserClaimsToContext(c, claims)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserClaims, ok := utils.GetCurrentUserClaims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range allowedRoles {
			if currentUserClaims.Role == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Anda tidak memiliki izin untuk mengakses sumber daya ini."})
			c.Abort()
			return
		}

		c.Next()
	}
}
