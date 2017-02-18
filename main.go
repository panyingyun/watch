package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codegangsta/cli"
	log "github.com/fatih/color"
	"github.com/panyingyun/watch/backend"
)

func run(c *cli.Context) error {
	//Connect to MQTT(for example:"tcp://127.0.0.1:1883", "pub", "pub")
	mqtt, err := backend.NewBackend(c.String("mqtt-server"), c.String("mqtt-username"), c.String("mqtt-password"))
	if err != nil {
		log.Red("can not connect mqtt server")
		return err
	}
	defer mqtt.Close()

	//Gateway Topic
	gw := c.String("gw-topic")
	if len(gw) > 0 {
		if err := mqtt.SubscribeTopic(gw); err != nil {
			log.Red("SubscribeTopic %v Error", gw)
			return err
		}
		defer mqtt.UnSubscribeTopic(gw)
	}

	//Application Topic
	app := c.String("app-topic")
	if len(app) > 0 {
		if err := mqtt.SubscribeTopic(app); err != nil {
			log.Red("SubscribeTopic %v Error", app)
			return err
		}
		defer mqtt.UnSubscribeTopic(app)
	}

	//When receive rxdata, then write to database
	go func() {
		for rxData := range mqtt.RxDataChan() {
			topic := rxData.Topic
			if strings.HasPrefix(topic, "gateway") {
				log.Green("[GW] topic = %v, msg = %v", topic, rxData.Msg)
			} else if strings.Contains(topic, "mac/rx") ||
				strings.Contains(topic, "mac/tx") ||
				strings.Contains(topic, "mac/error") ||
				strings.Contains(topic, "rxinfo") {
				log.Cyan("[MAC] topic = %v, msg = %v", topic, rxData.Msg)
			} else {
				log.Magenta("[APP] topic = %v, msg = %v", topic, rxData.Msg)
			}
		}
	}()

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.White("signal received signal %v", <-sigChan)
	log.White("shutting down server")
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
