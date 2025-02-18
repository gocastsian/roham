package postgresql

type Config struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	DBName          string `koanf:"db_name"`
	SSLMode         string `koanf:"ssl_mode"`
	MaxIdleConns    int    `koanf:"max_idle_conns"`
	MaxOpenConns    int    `koanf:"max_open_conns"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime"`
	PathOfMigration string `koanf:"path_of_migration"`
}
