package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/panyingyun/watch/backend"
)

func run(c *cli.Context) error {

	//Connect to MQTT(for example:"tcp://127.0.0.1:1883", "pub", "pub")
	mqtt, err := backend.NewBackend(c.String("mqtt-server"), c.String("mqtt-username"), c.String("mqtt-password"))
	if err != nil {
		log.Error("can not connect mqtt server")
		return err
	}
	defer mqtt.Close()

	//Gateway Topic
	gw := c.String("gw-topic")
	if len(gw) > 0 {
		if err := mqtt.SubscribeTopic(gw); err != nil {
			log.Errorf("SubscribeTopic %v Error", gw)
			return err
		}
		defer mqtt.UnSubscribeTopic(gw)
	}

	//Application Topic
	app := c.String("app-topic")
	if len(app) > 0 {
		if err := mqtt.SubscribeTopic(app); err != nil {
			log.Errorf("SubscribeTopic %v Error", app)
			return err
		}
		defer mqtt.UnSubscribeTopic(app)
	}

	//When receive rxdata, then write to database
	go func() {
		for rxData := range mqtt.RxDataChan() {
			if strings.HasPrefix(rxData.Topic, "gateway") {
				log.WithFields(log.Fields{
					"topic": rxData.Topic,
				}).Info("[GateWay]")
				log.WithFields(log.Fields{
					"msg": rxData.Msg,
				}).Info("[GateWay]")
			} else {
				log.WithFields(log.Fields{
					"topic": rxData.Topic,
				}).Info("[App]")
				log.WithFields(log.Fields{
					"msg": rxData.Msg,
				}).Info("[App]")
			}
		}
	}()

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Info("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "wath lowan message transport "
	app.Usage = "watch message received from lowan server with gateway and application"
	app.Copyright = "panyingyun(at)gmail.com"
	app.Version = "0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mqtt-server",
			Usage:  "MQTT server",
			Value:  "tcp://127.0.0.1:1883",
			EnvVar: "MQTT_SERVER",
		},
		cli.StringFlag{
			Name:   "mqtt-username",
			Usage:  "MQTT username",
			Value:  "pub",
			EnvVar: "MQTT_USERNAME",
		},
		cli.StringFlag{
			Name:   "mqtt-password",
			Usage:  "MQTT password",
			Value:  "pub",
			EnvVar: "MQTT_PASSWORD",
		},
		cli.StringFlag{
			Name:   "gw-topic",
			Usage:  "GW topic",
			Value:  "",
			EnvVar: "GW_TOPIC",
		},
		cli.StringFlag{
			Name:   "app-topic",
			Usage:  "APP topic",
			Value:  "",
			EnvVar: "APP_TOPIC",
		},
	}
	app.Run(os.Args)
}
