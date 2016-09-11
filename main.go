package main

import (

	//  "strconv"
	//	"strings"
	//	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/panyingyun/wath/backend"
	loralog "github.com/panyingyun/wath/loralog"
)

var app = "application/+/node/+/#"
var gw = "gateway/+/#"

func run(c *cli.Context) error {
	//Connect to MQTT(for example:"tcp://127.0.0.1:1883", "pub", "pub")
	mqtt, err := backend.NewBackend(c.String("mqtt-server"), c.String("mqtt-username"), c.String("mqtt-password"))
	if err != nil {
		loralog.Error("can not connect mqtt server")
		return err
	}
	defer mqtt.Close()

	//Application Topic
	if err := mqtt.SubscribeTopic(app); err != nil {
		loralog.Errorf("SubscribeTopic %v Error", app)
		return err
	}
	defer mqtt.UnSubscribeTopic(app)

	if err := mqtt.SubscribeTopic(gw); err != nil {
		loralog.Errorf("SubscribeTopic %v Error", gw)
		return err
	}
	defer mqtt.UnSubscribeTopic(gw)

	//When receive rxdata, then write to database
	go func() {
		for rxData := range mqtt.RxDataChan() {
			loralog.Debugf("rxData = %v", rxData)
		}
	}()

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	loralog.Infof("signal received signal %v", <-sigChan)
	loralog.Warn("shutting down server")
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
	}
	app.Run(os.Args)
}
