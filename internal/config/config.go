package config

type Config struct {
	Environment    string
	HTTPServerAddr string
	DSN            string 
}

// Load creates a new Config struct with the provided environment, HTTP server address and data source name.
func Load(environment, httpServerAddr, dsn string) Config {
	return Config{
		Environment:    environment,
		HTTPServerAddr: httpServerAddr,
		DSN:            dsn, 
	}
}
