package sentrylog

import (
	"log"


	"github.com/getsentry/sentry-go"
)

type Sentry struct {

}



func SentryLog() {
	err := sentry.Init(sentry.ClientOptions{

		Dsn: "https://1d5d87aa7d7c4cd6afdefcacda20833b@o1377776.ingest.sentry.io/6688875",

		TracesSampleRate: 1.0,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	
	log.Panicln("initialized sentry")
}