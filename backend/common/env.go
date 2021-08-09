package common

import (
	"github.com/kelseyhightower/envconfig"
	"hamster/lib/clog"
)

type envType struct {
	CoincheckApiKey    string `split_words:"true"`
	CoincheckApiSecret string `split_words:"true"`
}

var Env envType

func Init() {
	err := envconfig.Process("hamster", &Env)
	if err != nil {
		clog.Logger.Fatal(err)
	}
}
