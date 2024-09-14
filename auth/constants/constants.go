package constants

import "time"

const (
	DATABASE   string        = "olt_blueprint"
	USER_TABLE string        = "users"
	EXPIRE_JWT time.Duration = 30 * time.Minute
	SALT       int           = 14
)

const (
	FORBIDDEN_RESPONSE string = "you do not have permission to perform the operation"
)
