package common

import (
	"github.com/kelseyhightower/envconfig"
	"hamster/lib/clog"
	"time"
)

type envType struct {
	CoincheckApiKey    string `split_words:"true"`
	CoincheckApiSecret string `split_words:"true"`
}

var Env envType

const location = "Asia/Tokyo"

func Init() {
	// 環境変数
	err := envconfig.Process("hamster", &Env)
	if err != nil {
		clog.Logger.Fatal(err)
	}

	// 時刻
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}
