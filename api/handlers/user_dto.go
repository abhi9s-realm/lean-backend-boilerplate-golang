package handlers

// CreateUserRequest defines the structure for creating a new user.
// Validation tags are used by Gin to automatically validate the request body.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	// Add other fields like Password if required for creation
}

// UpdateUserRequest defines the structure for updating an existing user.
// Fields are optional (e.g. using pointers or relying on PATCH semantics if implemented).
// Here, we assume all fields are optional for PUT, and the service layer handles partial updates.
// Alternatively, use binding tags like "omitempty" or make fields pointers.
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
	// Add other updatable fields
}
