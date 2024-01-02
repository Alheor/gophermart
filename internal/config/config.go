package config

import (
	"flag"
	"os"
)

const defaultAddr = `localhost:8080`
const defaultDatabaseURL = `host=localhost port=5432 user=app password=pass dbname=gophermart sslmode=disable`
const DefaultLSignatureKey = `ceecb67c69b7415c162dbcd83fddddf3`
const defaultLogLevel = `debug`
const defaultAccrualSystemAddress = `http://localhost:8081`

const envRunAddress = `RUN_ADDRESS`
const envDatabaseURL = `DATABASE_URI`
const envAccrualSystemAddress = `ACCRUAL_SYSTEM_ADDRESS`
const envSignatureKey = `SIGNATURE_KEY`
const envLogLevel = `LOG_LEVEl`

var Options struct {
	Addr                 string
	DatabaseURI          string
	LogLevel             string
	AccrualSystemAddress string
	SignatureKey         string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, defaultAddr, "listening host:port")
	flag.StringVar(&Options.LogLevel, `l`, defaultLogLevel, "log handler level")
	flag.StringVar(&Options.DatabaseURI, `d`, defaultDatabaseURL, "database uri")
	flag.StringVar(&Options.AccrualSystemAddress, `r`, defaultAccrualSystemAddress, "accrual system address")
	flag.StringVar(&Options.SignatureKey, `k`, DefaultLSignatureKey, "signature secret key")
}

func Load() {
	flag.Parse()

	addr, exist := os.LookupEnv(envRunAddress)
	if exist && len(addr) > 0 {
		Options.Addr = addr
	}

	DatabaseURI, exist := os.LookupEnv(envDatabaseURL)
	if exist && len(DatabaseURI) > 0 {
		Options.DatabaseURI = DatabaseURI
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
