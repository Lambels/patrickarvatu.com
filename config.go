package pa

type Config struct {
	Github struct {
		ClientID     string `toml:"client-id"`
		ClientSecret string `toml:"client-secret"`
	} `toml:"github"`

	HTTP struct {
		Addr        string `toml:"addr"`
		Domain      string `toml:"domain"`
		BlockKey    string `toml:"block-key"`
		HashKey     string `toml:"hash-key"`
		FrontendURL string `toml:"frontend-url"`
	} `toml:"http"`

	User struct {
		AdminUserEmail string `toml:"admin-user-email"`
	} `toml:"user"`
}
