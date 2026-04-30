package factory

import (
	"fmt"
	"sync"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"
	"github.com/asaskevich/govalidator"

	"github.com/free5gc/openapi/models"
)

const (
	NfDefaultConfigPath     = "./config/nfcfg.yaml"
	NfDefaultTLSKeyLogPath  = "./log/nfsslkey.log"
	NfDefaultCertPemPath    = "./cert/nf.pem"
	NfDefaultPrivateKeyPath = "./cert/nf.key"
)

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
	Logger        *Logger        `yaml:"logger" valid:"required"`
	sync.RWMutex
}

type Info struct {
	Version     string `yaml:"version" valid:"required,in(1.0.1)"`
	Description string `yaml:"description,omitempty" valid:"type(string)"`
}

type Configuration struct {
	NfName     string      `yaml:"nfName,omitempty"`
	Sbi        *Sbi        `yaml:"sbi"`
	NrfUri     string      `yaml:"nrfUri" valid:"url,required"`
	NrfCertPem string      `yaml:"nrfCertPem,omitempty" valid:"optional"`
	WebUri     string      `yaml:"webUri,omitempty" valid:"url,optional"`
	GateKeeper *GateKeeper `yaml:"gateKeeper,omitempty" valid:"optional"`
}

type GateKeeper struct {
	Enable  bool     `yaml:"enable" valid:"type(bool)"`
	Servers []Server `yaml:"servers" valid:"required"`
}

type Server struct {
	Name string `yaml:"name" valid:"required"`
	Addr string `yaml:"addr" valid:"host,required"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
}

type Sbi struct {
	Scheme       models.UriScheme `yaml:"scheme"`
	BindingIPv4  string           `yaml:"bindingIPv4,omitempty" valid:"host,required"`
	RegisterIPv4 string           `yaml:"registerIPv4,omitempty" valid:"host,optional"`
	Port         int              `yaml:"port"`
	Tls          *Tls             `yaml:"tls,omitempty" valid:"optional"`
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (c *Config) Validate() (bool, error) {
	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

func (c *Configuration) validate() (bool, error) {
	if sbi := c.Sbi; sbi != nil {
		if result, err := sbi.validate(); err != nil {
			return result, err
		}
	}
	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

func (s *Sbi) validate() (bool, error) {
	govalidator.TagMap["scheme"] = govalidator.Validator(func(str string) bool {
		return str == "https" || str == "http"
	})

	if tls := s.Tls; tls != nil {
		if result, err := tls.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(s)
	return result, appendInvalid(err)
}

func (t *Tls) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(t)
	return result, err
}

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	c.RLock()
	defer c.RUnlock()

	if c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) SetLogEnable(enable bool) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Enable: enable,
			Level:  "info",
		}
	} else {
		c.Logger.Enable = enable
	}
}

func (c *Config) SetLogLevel(level string) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level: level,
		}
	} else {
		c.Logger.Level = level
	}
}

func (c *Config) SetLogReportCaller(reportCaller bool) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level:        "info",
			ReportCaller: reportCaller,
		}
	} else {
		c.Logger.ReportCaller = reportCaller
	}
}

func (c *Config) GetLogEnable() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.Enable
}

func (c *Config) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return "info"
	}
	return c.Logger.Level
}

func (c *Config) GetLogReportCaller() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.ReportCaller
}
