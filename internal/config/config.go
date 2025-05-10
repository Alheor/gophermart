package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

// DefaultLSignatureKey signature key for user authentication
const DefaultLSignatureKey = `40d40c8d1b5fff17e7edcabc6b2fa4ab`

type Options struct {
	RunAddr      string `env:"RUN_ADDRESS"`
	DatabaseUri  string `env:"DATABASE_URI"`
	AccrualAddr  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SignatureKey string `env:"SIGNATURE_KEY"`
}

var options Options

func init() {
	flag.StringVar(&options.RunAddr, `a`, `localhost:8080`, "listen host/ip:port")
	flag.StringVar(&options.DatabaseUri, `d`, ``, "database dsn")
	flag.StringVar(&options.AccrualAddr, `b`, `http://localhost:8080`, "accrual system address")
	flag.StringVar(&options.SignatureKey, `k`, DefaultLSignatureKey, "signature key")
}

func Load() Options {

	flag.Parse()

	err := env.Parse(&options)
	if err != nil {
		log.Fatal(err)
	}

	println(`--- Loaded configuration ---`)

	println(`listen: ` + options.RunAddr)
	println(`database uri: ` + options.DatabaseUri)
	println(`accrual system address: ` + options.AccrualAddr)

	if options.SignatureKey == DefaultLSignatureKey {
		println(`signature key status: used default key`)
	} else {
		println(`signature key status: key specified by parameter`)
	}

	return options
}
