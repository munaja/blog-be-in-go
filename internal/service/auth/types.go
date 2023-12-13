package auth

type TokenType string

const AccessToken = "Access"
const RefreshToken = "Refresh"

type AuthInfo struct {
	Uuid       string
	User_Id    int
	User_Name  string
	User_Email string
}
