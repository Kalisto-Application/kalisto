//go:build production
// +build production

package config

func init() {
	C = Conf{
		SentryDsn: "https://3b1581aea95a4abfb4437f38b1f35b66@o4505605913444352.ingest.sentry.io/4505605917573120",
	}
}
