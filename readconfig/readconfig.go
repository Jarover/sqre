package readconfig

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// Config - структура для считывания конфигурационного файла
type Config struct {
	//Db_url       string `yaml:"db_url" json:"db_url"`
	Port uint `yaml:"port" json:"port" `
	//Host  string `yaml:"host" json:"host"`
	//Jaeger_url   string `yaml:"jaeger_url" json:"jaeger_url"`
	//Sentry_url   string `yaml:"sentry_url" json:"sentry_url"`
	//Kafka_broker string `yaml:"kafka_broker" json:"kafka_broker"`
	//Some_app_id  string `yaml:"some_app_id" json:"some_app_id"`
	//Some_app_key string `yaml:"some_app_key" json:"some_app_key"`
}

func (e *Config) Validate() error {
	var err error
	/*
		err = e.CheckUrl(e.Db_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Jaeger_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Sentry_url)
		if err != nil {
			return err
		}

		err = e.CheckUrl(e.Kafka_broker)
		if err != nil {
			return err
		}
	*/
	return err
}

func (e *Config) SetPort(p uint) error {

	e.Port = p
	return nil

}

/*
func (e *Config) SetDb(p string) error {

	err := e.CheckUrl(p)
	if err == nil {
		e.Db_url = p
		return nil
	} else {
		return err

	}

}
*/

func (e *Config) CheckUrl(path string) error {

	_, err := url.ParseRequestURI(path)

	if err != nil {
		return err
	}
	return nil

}

func ReadConfig(ConfigName string) (x *Config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile(ConfigName); err != nil {
		return nil, err
	}
	x = new(Config)
	switch strings.ToLower(path.Ext(ConfigName)) {

	case ".yaml", ".yml":
		err = yaml.Unmarshal(file, &x)
	case ".json":
		err = json.Unmarshal(file, &x)

	case ".ini":
		cfg, err := ini.Load(ConfigName)
		if err == nil {
			x.Port = cfg.Section("").Key("port").MustUint()
			//x.Port2 = cfg.Section("").Key("port2").MustUint()
			//x.Db_url = cfg.Section("").Key("db_url").String()
			//x.Jaeger_url = cfg.Section("").Key("jaeger_url").String()
			//x.Sentry_url = cfg.Section("").Key("sentry_url").String()
			//x.Kafka_broker = cfg.Section("").Key("kafka_broker").String()
			//x.Some_app_id = cfg.Section("").Key("some_app_id").String()
			//x.Some_app_key = cfg.Section("").Key("some_app_key").String()
		}

	}

	if err != nil {
		return nil, err
	}

	if err = x.Validate(); err != nil {
		return nil, err
	}

	return x, nil
}
