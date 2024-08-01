package service

import (
	"time"

	"stool/jms-sdk-go/httplib"
	"stool/jms-sdk-go/model"
)

type option struct {
	// default http://127.0.0.1:8080
	CoreHost  string
	TimeOut   time.Duration
	sign      httplib.AuthSign
	accessKey model.AccessKey
	Insecure  bool
}

type Option func(*option)

func JMSCoreHost(coreHost string) Option {
	return func(o *option) {
		o.CoreHost = coreHost
	}
}

func JMSTimeOut(t time.Duration) Option {
	return func(o *option) {
		o.TimeOut = t
	}
}

func JMSAccessKey(keyID, secretID string) Option {
	return func(o *option) {
		o.sign = &httplib.SigAuth{
			KeyID:    keyID,
			SecretID: secretID,
		}
		o.accessKey = model.AccessKey{
			ID:     keyID,
			Secret: secretID,
		}
	}
}

func JMSAuthSign(sign httplib.AuthSign) Option {
	return func(o *option) {
		o.sign = sign
	}
}

func JMSInsecure() Option {
	return func(o *option) {
		o.Insecure = true
	}
}
