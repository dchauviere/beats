package rabbitmq

type config struct {
	Host          string `config:"host"`
  Port          int    `config:"port"`
	Username			string `config:"username"`
	Password      string `config:"password"`
	Exchange      string `config:"exchange"`
	ExchangeType  string `config:"exchange_type"`
	Durable				bool `config:"durable"`
	AutoDelete		bool   `config:"auto_delete"`
	Internal			bool	 `config:"internal"`
	NoWait				bool	 `config:"no_wait"`
	Arguments			string `config:"arguments"`
	Key						string `config:"key"`
	Mandatory			bool   `config:"mandatory"`
	Immediate			bool   `config:"immediate"`
	DeliveryMode  uint8  `config:"delivery_mode"`
	Priority			uint8  `config:"priority"`
}

var (
	defaultConfig = config{
		Host: "",
    Port: 5672,
		Username: "guest",
		Password: "guest",
		Exchange: "logs",
		ExchangeType: "direct",
		Durable: true,
		AutoDelete: false,
		Internal: false,
		NoWait: false,
		Arguments: "",
		Key: "logs",
		Mandatory: false,
		Immediate: false,
		DeliveryMode: 1,
		Priority: 0,
	}
)
