package usecase

type Authenticator interface {
	CreateAccessToken(userID int64) (accessToken string, err error)
	CreateRefreshToken(userID int64) (refreshToken string, err error)
	VerifyAccessToken(token string) (userID int64, err error)
	VerifyRefreshToken(token string) (userID int64, err error)
}
