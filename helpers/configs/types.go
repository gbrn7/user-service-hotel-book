package configs

type (
	Config struct {
		Service    Service
		Database   Database
		GrpcConfig GrpcConfig
		JwtConfig  JwtConfig
	}

	Service struct {
		AppName               string  `json:"appName"`
		AppEnv                string  `json:"appEnv"`
		Port                  int     `json:"port"`
		SignatureKey          string  `json:"signatureKey"`
		RateLimiterMaxRequest float64 `json:"rateLimiterMaxRequest"`
		RateLimiterMaxSecond  int     `json:"rateLimiterMaxSecond"`
	}

	Database struct {
		Host                  string `json:"host"`
		Port                  int    `json:"port"`
		Name                  string `json:"name"`
		Username              string `json:"username"`
		Password              string `json:"password"`
		MaxOpenConnection     int    `json:"maxOpenConnection"`
		MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
		MaxIdleConnection     int    `json:"maxIdleConnection"`
		MaxIdleTime           int    `json:"maxIdleTime"`
	}

	JwtConfig struct {
		JwtSecretKey      string `json:"jwtSecretKey"`
		JwtExpirationTime int    `json:"jwtExpirationTime"`
	}

	GrpcConfig struct {
		Host string `json:"host"`
	}

	ElasticConfig struct {
		ApMServerUrl string `json:"apmServerUrl"`
	}
)
