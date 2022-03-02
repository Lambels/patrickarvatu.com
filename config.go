package pa

type Config struct {
	Github struct {
		ClientID       string `mapstructure:"client-id"`
		ClientSecret   string `mapstructure:"client-secret"`
		AdminUserEmail string `mapstructure:"admin-user-email"`
	} `mapstructure:"github"`

	HTTP struct {
		Addr        string `mapstructure:"addr"`
		Domain      string `mapstructure:"domain"`
		BlockKey    string `mapstructure:"block-key"`
		HashKey     string `mapstructure:"hash-key"`
		FrontendURL string `mapstructure:"frontend-url"`
	} `mapstructure:"http"`

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

	FileStructure struct {
		ProjectImagesDir string `mapstructure:"project-images-dir"`
		BlogImagesDir    string `mapstructure:"blog-images-dir"`
	} `mapstructure:"file-structure"`
}
