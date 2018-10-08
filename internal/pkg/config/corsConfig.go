package config

const corsConfigName = "cors"

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	MaxAge           int
	AllowCredentials bool
}
