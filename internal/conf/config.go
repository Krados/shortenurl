package conf

type Config struct {
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type Mysql struct {
	DSN string `yaml:"source"`
}

type Data struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}

type Redis struct {
	Addrs []string `yaml:"addrs"`
}
