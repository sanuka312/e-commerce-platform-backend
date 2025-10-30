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
