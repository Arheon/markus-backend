package helpers

func IsDevEnv(env string) bool {
	return env == "dev"
}
func IsProdEnv(env string) bool {
	return env == "prod"
}
