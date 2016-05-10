package syslog

import (
	"encoding/json"
  "log/syslog"
	"strings"
	"strconv"
	"bytes"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
)

func init() {
	outputs.RegisterOutputPlugin("syslog", New)
}

type syslogOutput struct {
  Protocol string
  LogLevel syslog.Priority
  ProgName string
  Prefix []byte
	Suffix []byte
  log *syslog.Writer
}

// New instantiates a new file output instance.
func New(cfg *common.Config, _ int) (outputs.Outputer, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, err
	}

	// disable bulk support in publisher pipeline
	cfg.SetInt("flush_interval", -1, -1)
	cfg.SetInt("bulk_max_size", -1, -1)

	output := &syslogOutput{}
	if err := output.init(config); err != nil {
		return nil, err
	}
	return output, nil
}

func (out *syslogOutput) init(config config) error {
	dest := ""
	if config.Host != "" {
		dest = strings.Join([]string{config.Host,strconv.Itoa(config.Port)}, ":")
		if config.Protocol != "" {
	    out.Protocol = config.Protocol
	  }
	}

  switch config.LogLevel {
    case "debug": out.LogLevel = syslog.LOG_DEBUG
    case "info": out.LogLevel = syslog.LOG_INFO
    case "notice": out.LogLevel = syslog.LOG_NOTICE
    case "warning": out.LogLevel = syslog.LOG_WARNING
    case "err": out.LogLevel = syslog.LOG_ERR
    case "crit": out.LogLevel = syslog.LOG_CRIT
    case "alert": out.LogLevel = syslog.LOG_ALERT
    case "emerg": out.LogLevel = syslog.LOG_EMERG
  default: out.LogLevel = syslog.LOG_INFO
  }

  out.Prefix = []byte(config.Prefix)
	out.Suffix = []byte(config.Suffix)
  out.ProgName = config.ProgName

  logger, err := syslog.Dial(out.Protocol, dest, out.LogLevel, out.ProgName)
	out.log = logger
  if err != nil {
    logp.Err("Error opening syslog: %v", err)
    return err
  }

	return nil
}

func (out *syslogOutput) Close() error {
  out.log.Close()
	return nil
}

func (out *syslogOutput) PublishEvent(
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

  _, err = out.log.Write( bytes.Join([][]byte{out.Prefix, jsonEvent, out.Suffix}, []byte("") ))
	if err != nil {
		if opts.Guaranteed {
			logp.Critical("Unable to write events to syslog: %s", err)
		} else {
			logp.Err("Error when writing line to syslog: %s", err)
		}
	}
	op.Sig(sig, err)
	return err
}
