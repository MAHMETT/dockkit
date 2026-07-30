package config

// DefaultConfig returns a fresh config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version: VersionCurrent,
		General: GeneralConfig{
			Timezone:        "UTC",
			DefaultNetwork:  "dockkit-network",
			AutoRefresh:     true,
			RefreshInterval: "30s",
		},
		Services: map[string]Service{},
	}
}

// DefaultServiceVersions returns default version configs for built-in services.
func DefaultServiceVersions() map[string]Service {
	return map[string]Service{
		"postgresql": {
			Prefix: "PG",
			Versions: map[string]ServiceVersion{
				"16": {
					Enabled:       false,
					Port:          5433,
					ContainerName: "dockkit-postgresql-16",
					Image:         "postgres:16-alpine",
					User:          "postgres",
					Password:      "postgres",
					Database:      "postgres",
				},
			},
		},
		"mysql": {
			Prefix: "MYSQL",
			Versions: map[string]ServiceVersion{
				"8": {
					Enabled:       false,
					Port:          3306,
					ContainerName: "dockkit-mysql-8",
					Image:         "mysql:8.0",
					User:          "root",
					Password:      "mysql",
					Database:      "mysql",
				},
			},
		},
		"redis": {
			Prefix: "REDIS",
			Versions: map[string]ServiceVersion{
				"7": {
					Enabled:       false,
					Port:          6379,
					ContainerName: "dockkit-redis-7",
					Image:         "redis:7-alpine",
				},
			},
		},
	}
}
