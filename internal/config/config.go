package config

//Config holds all configuration for this service
type Config struct {
	Version     string
	ServiceName string
	LogLevel    string
	DbPath      string
	HTTPPort    int
	JWTSecret   string
}
