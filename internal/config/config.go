package config

import (
	"flag"
	"os"
)

const defaultAddr = `localhost:8081`
const defaultDatabaseUri = `host=localhost port=5432 user=app password=pass dbname=gophermart sslmode=disable`
const DefaultLSignatureKey = `ceecb67c69b7415c162dbcd83fddddf3`
const defaultLogLevel = `debug`
const defaultAccrualSystemAddress = `http://localhost:8080`

const envRunAddress = `RUN_ADDRESS`
const envDatabaseUri = `DATABASE_URI`
const envAccrualSystemAddress = `ACCRUAL_SYSTEM_ADDRESS`
const envSignatureKey = `SIGNATURE_KEY`
const envLogLevel = `LOG_LEVEl`

var Options struct {
	Addr                 string
	DatabaseUri          string
	LogLevel             string
	AccrualSystemAddress string
	SignatureKey         string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, defaultAddr, "listening host:port")
	flag.StringVar(&Options.LogLevel, `l`, defaultLogLevel, "log handler level")
	flag.StringVar(&Options.DatabaseUri, `d`, defaultDatabaseUri, "database uri")
	flag.StringVar(&Options.AccrualSystemAddress, `r`, defaultAccrualSystemAddress, "accrual system address")
	flag.StringVar(&Options.SignatureKey, `k`, DefaultLSignatureKey, "signature secret key")
}

func Load() {
	flag.Parse()

	addr, exist := os.LookupEnv(envRunAddress)
	if exist && len(addr) > 0 {
		Options.Addr = addr
	}

	DatabaseUri, exist := os.LookupEnv(envDatabaseUri)
	if exist && len(DatabaseUri) > 0 {
		Options.DatabaseUri = DatabaseUri
	}

	AccrualSystemAddress, exist := os.LookupEnv(envAccrualSystemAddress)
	if exist && len(AccrualSystemAddress) > 0 {
		Options.AccrualSystemAddress = AccrualSystemAddress
	}

	SignatureKey, exist := os.LookupEnv(envSignatureKey)
	if exist && len(SignatureKey) > 0 {
		Options.SignatureKey = SignatureKey
	}

	logLevel, exist := os.LookupEnv(envLogLevel)
	if exist && len(logLevel) > 0 {
		Options.LogLevel = logLevel
	}
}
