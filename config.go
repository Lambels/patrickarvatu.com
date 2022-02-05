package pa

type Config struct {
	Github struct {
		ClientID     string `mapstructure:"client-id"`
		ClientSecret string `mapstructure:"client-secret"`
	} `mapstructure:"github"`

	HTTP struct {
		Addr        string `mapstructure:"addr"`
		Domain      string `mapstructure:"domain"`
		BlockKey    string `mapstructure:"block-key"`
		HashKey     string `mapstructure:"hash-key"`
		FrontendURL string `mapstructure:"frontend-url"`
	} `mapstructure:"http"`

	User struct {
		AdminUserEmail string `mapstructure:"admin-user-email"`
	} `mapstructure:"user"`

	Database struct {
		SqliteDSN string `mapstructure:"sqlite-dsn"`
		RedisDSN  string `mapstructure:"redis-dsn"`
	} `mapstructure:"database"`

	Smtp struct {
		Addr     string `mapstructure:"addr"`
		Identity string `mapstructure:"identity"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
	} `mapstructure:"smtp"`
}
