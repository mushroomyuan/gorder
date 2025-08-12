package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	_ = viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (relPath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, here, _, _ := runtime.Caller(0)
	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	logrus.Infof("caller from:%s,here:%s,relpath:%s", callerPwd, here, relPath)
	return
}
