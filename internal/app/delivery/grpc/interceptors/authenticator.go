package interceptors

type Authenticator interface {
	VerifyAccessToken(token string) (userID int64, err error)
}
