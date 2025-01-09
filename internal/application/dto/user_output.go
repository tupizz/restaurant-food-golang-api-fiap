package dto

type UserOutput struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
