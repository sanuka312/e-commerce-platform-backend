package auth

import (
	"bytes"
	"e-commerce-platform-backend/config"
	"e-commerce-platform-backend/data"
	"e-commerce-platform-backend/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Key for storing the claims in the context
type contextKey string

const RolesContextKey = contextKey("roles")
const UserNameContextKey = contextKey("user_name")

// AuthMiddleware is a middleware that checks the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			logger.ActError("Missing or invalid Authorization header", zap.String("endpoint", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, data.ErrorResponse{
				Error:            "invalid_authorization",
				ErrorDescription: "Missing or invalid authorization",
			})
			return
		}

		tokenString := authHeader[7:]

		// Introspect Token Validity before token parse
		result, err := introspectToken(tokenString)
		if err != nil {
			logger.ActError("Token Introspect Error", zap.String("endpoint", c.Request.URL.Path), zap.Error(err))
			panic(err)
		}

		// Check result of token introspect
		if !result.Active {
			logger.ActError("Token is invalid or revoked!", zap.String("endpoint", c.Request.URL.Path), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, data.ErrorResponse{
				Error:            "invalid_token_revoked",
				ErrorDescription: "Token is invalid or revoked",
			})
			return
		}
		// Store roles in Gin context
		c.Set((string(RolesContextKey)), GetResourceRoles(result.ResourceAccess))

		c.Set((string(UserNameContextKey)), result.PreferredUsername)

		// Continue to the next checkpoint
		c.Next()
	}
}

// IntrospectToken performs an introspection of a given token on the IdP server to determine its validity.
// The response will contain the result of the introspection, which includes whether the token is active or not.
// If the token is invalid or revoked, an error will be returned with the appropriate status code.
func introspectToken(token string) (data.IntrospectResponse, error) {
	var baseUrl string = config.LoadConfig().IdpBaseUrl
	var realm string = config.LoadConfig().IdpRealm
	var clientID string = config.LoadConfig().IdpClientId
	var clientSecret string = config.LoadConfig().IdpClientSecret

	var endpoint string = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", baseUrl, realm)

	logger.ActDebug("Introspecting Token", zap.String("endpoint", endpoint))

	form := url.Values{}
	form.Set("token", token)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return data.IntrospectResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return data.IntrospectResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return data.IntrospectResponse{}, fmt.Errorf("introspection failed: %s", string(bodyBytes))
	}

	logger.ActDebug("bodyBytes", zap.String("bodyBytes", string(bodyBytes)))
	var introspectResp data.IntrospectResponse
	if err := json.Unmarshal(bodyBytes, &introspectResp); err != nil {
		return data.IntrospectResponse{}, err
	}

	logger.ActDebug("introspectResp", zap.Any("introspectResp", introspectResp))

	return introspectResp, nil
}
