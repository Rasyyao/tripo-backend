package dto

// MessageResponse is the universal body for actions that return no resource
// (logout, delete, ...). Build one with NewMessageResponse.
type MessageResponse struct {
	Message string `json:"message"`
}

func NewMessageResponse(message string) *MessageResponse {
	return &MessageResponse{Message: message}
}
