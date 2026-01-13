package models

// User roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Valid roles list
var ValidRoles = []string{RoleAdmin, RoleUser}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {

			return true
		}
	}
	return false
}
