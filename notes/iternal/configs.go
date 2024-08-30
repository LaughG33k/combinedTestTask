package iternal

type ConfigApp struct {
	Addr              string   `yaml:"Addr"`
	ReadTimeoutInSec  int      `yaml:"ReadTimeoutInSec"`
	WriteTimeoutInSec int      `yaml:"WriteTimeoutInSec"`
	IdleTimeoutInSec  int      `yaml:"IdleTimeoutInSec"`
	NoteDB            DBConfig `yaml:"NoteDB"`
	TrustedAddrs      []string `yaml:"TrustedAddrs"`
}

type DBConfig struct {
	Host                string `yaml:"Host"`
	Port                uint16 `yaml:"Port"`
	User                string `yaml:"User"`
	Password            string `yaml:"Password"`
	DB                  string `yaml:"DB"`
	PoolSize            int    `yaml:"PoolSize"`
	TryAttempts         int    `yaml:"TryAttempts"`
	NewConnTimeoutInSec int    `yaml:"NewConnTimeoutInSec"`
}
