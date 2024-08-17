package m

import "time"

var (
	expiresIn = time.Duration(60)
	Issuer    = "testIssuer"
	secretKey = "testSecretKey"
	Id        = "testId"

	requestHeaderAuthKey = "Authorization"
	whiteList            = []string{"GET:/api/public", "POST:*"}
)
