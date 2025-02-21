package proxy

type Config struct {
	Host     string
	User     string
	Password string
}

func (c Config) IsEnabled() bool {
	return c.Host != ""
}
func (c Config) AuthEnabled() bool {
	return c.Password != "" && c.User != ""
}
