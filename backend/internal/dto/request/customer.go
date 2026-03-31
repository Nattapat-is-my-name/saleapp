package request

type CreateCustomerRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"max=20"`
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	Address   string `json:"address"`
	Notes     string `json:"notes"`
}

type UpdateCustomerRequest struct {
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone" binding:"omitempty,max=20"`
	FirstName *string `json:"first_name" binding:"omitempty,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,max=100"`
	Address   *string `json:"address"`
	Notes     *string `json:"notes"`
}

type ListCustomersRequest struct {
	Page   int    `form:"page" binding:"min=0"`
	Limit  int    `form:"limit" binding:"min=0,max=100"`
	Search string `form:"search"`
}
