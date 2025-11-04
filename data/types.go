package data

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	Details          string `json:"details,omitempty"`
}

type AddToCartRequest struct {
	UserID    uint `json:"user_id"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type RealmAccess struct {
	Roles []string `json:"roles,omitempty"`
}

// ResourceRoles contains roles for a specific resource
type ResourceRoles map[string]struct {
	Roles []string `json:"roles,omitempty"`
}

// Identity Provider Audience type
type Audience []string

type IntrospectResponse struct {
	Exp               int64         `json:"exp"`
	Iat               int64         `json:"iat"`
	AuthTime          int64         `json:"auth_time"`
	Jti               string        `json:"jti"`
	Iss               string        `json:"iss"`
	Aud               Audience      `json:"aud"`
	Sub               string        `json:"sub"`
	Typ               string        `json:"typ"`
	Azp               string        `json:"azp"`
	Sid               string        `json:"sid"`
	Acr               string        `json:"acr"`
	AllowedOrigins    []string      `json:"allowed-origins"`
	RealmAccess       RealmAccess   `json:"realm_access"`
	ResourceAccess    ResourceRoles `json:"resource_access"`
	Scope             string        `json:"scope"`
	EmailVerified     bool          `json:"email_verified"`
	Name              string        `json:"name"`
	PreferredUsername string        `json:"preferred_username"`
	GivenName         string        `json:"given_name"`
	FamilyName        string        `json:"family_name"`
	Email             string        `json:"email"`
	ClientID          string        `json:"client_id"`
	Username          string        `json:"username"`
	TokenType         string        `json:"token_type"`
	Active            bool          `json:"active"`
}

// Payment Request Struct
type ProcessPaymentRequest struct {
	PaymentMethod string `json:"payment_method"`
}
