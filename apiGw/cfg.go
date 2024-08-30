package apigw

type Cfg struct {
	ReadTimeoutInSec     int        `yaml:"ReadTimeoutInSec"`
	WriteTimeoutInSec    int        `yaml:"WriteTimeoutInSec"`
	Addr                 string     `yaml:"Addr"`
	TimeoutInSec         int        `yaml:"TimeoutInSec"`
	IdleConnTimeoutInSec int        `yaml:"IdleConnTimeoutInSec"`
	MaxIdleConns         int        `yaml:"MaxIdleConns"`
	MaxConnsPerEndPoint  int        `yaml:"MaxConnsPerEndPoint"`
	EndPoint             []EndPoint `yaml:"EndPoint"`
	AuthApi              string     `yaml:"AuthApi"`
}
