package cmd

import (
	"time"

	"github.com/jgarrone/foaas-limiting/server"
	"github.com/jgarrone/foaas-limiting/services/foaasapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serveFlags struct {
	address     string
	logLevel    string
	limitCount  int
	limitWindow int
}

const (
	defaultAddress           = "localhost:8080"
	defaultLogLevel          = ""
	defaultRateLimitCount    = 5
	defaultRateLimitWindowMS = 10000
)

func ServeCommand() *cobra.Command {
	var flags serveFlags

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the FOAAS Limiting server",
		Long:  `Run the FOAAS Limiting server`,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			setLogLevel(&flags)
		},
		Run: func(_ *cobra.Command, _ []string) {
			runServer(&flags)
		},
	}

	cmd.Flags().StringVar(&flags.address, "listen_address", defaultAddress,
		"address for listening to clients")
	cmd.Flags().StringVar(&flags.logLevel, "log_level", defaultLogLevel,
		"log level (empty means default)")
	cmd.Flags().IntVar(&flags.limitCount, "rate_limit_count", defaultRateLimitCount,
		"maximum number of requests a user can make inside a window of time")
	cmd.Flags().IntVar(&flags.limitWindow, "rate_limit_window_ms", defaultRateLimitWindowMS,
		"period of time in milliseconds to limit the amount of requests a user can make")

	return cmd
}

func setLogLevel(flags *serveFlags) {
	level := flags.logLevel
	if level == "" {
		return
	}
	ll, err := log.ParseLevel(level)
	if err != nil {
		log.Warnf("error parsing log level %q: %v", level, err)
		return
	}
	log.Infof("Setting log level to %s", ll)
	log.SetLevel(ll)
}

func runServer(flags *serveFlags) {
	limitWindow := time.Duration(flags.limitWindow) * time.Millisecond
	limiter := server.NewTokenBucketLimiter(flags.limitCount, limitWindow)
	log.Infof("Limiter set to allow %d requests every %v", flags.limitCount, limitWindow)

	sv := server.New(flags.address, limiter, foaasapi.NewService())
	sv.Run()
}
