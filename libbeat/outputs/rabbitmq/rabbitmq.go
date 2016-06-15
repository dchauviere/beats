package rabbitmq

import (
	"encoding/json"
	"strings"
	"strconv"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/streadway/amqp"
)

func init() {
	outputs.RegisterOutputPlugin("rabbitmq", New)
}

type rabbitmqOutput struct {
  Url string
	Exchange string
	Key string
	Mandatory bool
	Immediate bool
	DeliveryMode uint8
	Priority uint8
  Channel *amqp.Channel
}

// New instantiates a new file output instance.
func New(cfg *common.Config, _ int) (outputs.Outputer, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, err
	}

	output := &rabbitmqOutput{}
	if err := output.init(config); err != nil {
		return nil, err
	}
	return output, nil
}

func (out *rabbitmqOutput) init(config config) error {
	dest := ""
	if config.Host != "" {
		dest = strings.Join([]string{"amqp://",config.Username,":",config.Password,"@",config.Host,":",strconv.Itoa(config.Port)}, "")
	}

	conn, err := amqp.Dial(dest)
	if err != nil {
    logp.Err("Failed to connect to RabbitMQ")
    return err
  }
	//defer conn.Close()

	out.Channel, err = conn.Channel()
	if err != nil {
    logp.Err("Failed to open a channel")
    return err
  }
	//defer ch.Close()

	err = out.Channel.ExchangeDeclare(
		config.Exchange, // name
		config.ExchangeType,      // type
		config.Durable,         // durable
		config.AutoDelete,        // auto-deleted
		config.Internal,        // internal
		config.NoWait,        // no-wait
		nil,          // arguments
	)
	if err != nil {
    logp.Err("Failed to declare an exchange")
    return err
  }
	out.Exchange = config.Exchange
	out.Key = config.Key
	out.Mandatory = config.Mandatory
	out.Immediate = config.Immediate
	out.DeliveryMode = config.DeliveryMode
	out.Priority = config.Priority
	return nil
}

func (out *rabbitmqOutput) Close() error {
  out.Channel.Close()
	return nil
}

func (out *rabbitmqOutput) PublishEvent(
	sig op.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		// mark as success so event is not sent again.
		op.SigCompleted(sig)

		logp.Err("Fail to json encode event(%v): %#v", err, event)
		return err
	}

	err = out.Channel.Publish(
			out.Exchange,          // exchange
			out.Key, // routing key
			out.Mandatory, // mandatory
			out.Immediate, // immediate
			amqp.Publishing{
				Headers: amqp.Table{},
				ContentType: "application/json",
				Body: []byte(jsonEvent),
				DeliveryMode:    out.DeliveryMode, // 1=non-persistent, 2=persistent
				Priority: out.Priority, // 0-9
	})

	if err != nil {
		if opts.Guaranteed {
			logp.Critical("Unable to write events to rabbitmq: %s", err)
		} else {
			logp.Err("Error when writing line to rabbitmq: %s", err)
		}
	}
	op.Sig(sig, err)
	return err
}
