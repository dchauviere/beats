package syslog

type config struct {
	Host          string `config:"host"`
  Port          int    `config:"port"`
	Protocol      string `config:"protocol"`
  LogLevel      string `config:"loglevel"`
  ProgName      string `config:"progname"`
  Prefix        string `config:"prefix"`
	Suffix        string `config:"suffix"`
}

var (
	defaultConfig = config{
		Host: "",
    Port: 514,
    Protocol: "tcp",
    LogLevel: "info",
    ProgName: "",
    Prefix: "",
		Suffix: "",
	}
)
