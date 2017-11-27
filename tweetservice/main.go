package main

import (
        "github.com/djdnl13/twitter/tweetservice/service"
        "github.com/djdnl13/twitter/tweetservice/dbclient"
        "github.com/djdnl13/twitter/common/config"
        "github.com/djdnl13/twitter/common/messaging"
        "github.com/Sirupsen/logrus"
	"flag"
        "github.com/spf13/viper"
	"os/signal"
	"os"
	"syscall"
)

var appName = "tweetservice"

func init() {
        profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	if *profile == "dev" {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp: true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
        configServerUrl := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")
        configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")

        flag.Parse()

        viper.Set("profile", *profile)
        viper.Set("configServerUrl", *configServerUrl)
        viper.Set("configBranch", *configBranch)
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
        logrus.Infof("Starting %v\n", appName)

        config.LoadConfigurationFromBranch(
                viper.GetString("configServerUrl"),
                appName,
                viper.GetString("profile"),
                viper.GetString("configBranch"))
        initializeBoltClient()
        initializeMessaging()
        handleSigterm(func() {
                service.MessagingClient.Close()
        })
        service.StartWebServer(viper.GetString("server_port"))
}

func initializeMessaging() {
        if !viper.IsSet("amqp_server_url") {
                panic("No 'amqp_server_url' set in configuration, cannot start")
        }

        service.MessagingClient = &messaging.MessagingClient{}
        service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
        service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
}

func initializeBoltClient() {
        service.DBClient = &dbclient.BoltClient{}
        service.DBClient.OpenBoltDb()
        service.DBClient.Seed()
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        signal.Notify(c, syscall.SIGTERM)
        go func() {
                <-c
                handleExit()
                os.Exit(1)
        }()
}

