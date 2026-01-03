package auth

// Claims represents the standard JWT claim set used by this service.
type Claims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}
