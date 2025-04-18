package initializer

import (
	jwthelper "socialmedia/auth/infra/jwthelper"
	"socialmedia/auth/pkg/config"
)

func InitJwtHelper(appConfig config.Config) *jwthelper.JwtHelperService {
	jwtService := jwthelper.NewJwtHelperService(appConfig.JWT.Secret)
	return jwtService
}
