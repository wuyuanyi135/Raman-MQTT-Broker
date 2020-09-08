package main

import (
	"context"
	"flag"
	"github.com/reactivex/rxgo/observer"
	"github.com/ryanuber/go-glob"
	"github.com/sirupsen/logrus"
	"github.com/wuyuanyi135/mqttraman/fileobserver"
	"github.com/wuyuanyi135/mqttraman/mqtt"
	"github.com/wuyuanyi135/mqttraman/process"
	"os"
	"sync"
	"time"
)

func main() {
	filter := flag.String("filter", "*.csv", "New file filter")
	serverUrl := flag.String("server", "tcp://192.168.43.1:1883", "mqtt server url")
	waitWrite := flag.Duration("wait", 3 * time.Second, "Time delay to read the file after it is created")
	flag.Parse()
	if flag.NArg() == 0 {
		os.Exit(0)
	}
	args := flag.Args()

	client := mqtt.NewMqttClient(*serverUrl)

	wg := sync.WaitGroup{}
	for _, value := range args {
		wg.Add(1)

		err, ob := fileobserver.StartMonitor(value, context.Background())
		if err != nil {
			panic(err)
		}
		obs := observer.Observer{

			// Register a handler function for every next available item.
			NextHandler: func(item interface{}) {
				s, ok := item.(string)
				if !ok {
					panic("file path is not a string")
				}
				if filter != nil {
					if glob.Glob(*filter, s) {
						logrus.Info("Detected new file " + s)
						<-time.After(*waitWrite)
						intensity, err := process.ProcessFile(s)
						if err != nil {
							logrus.Error(err)
						}

						err = client.PublishSpectrumData(intensity)
						if err != nil {
							logrus.Error(err)
						}
					} else {
						logrus.Info("Ignored new file " + s)
					}
				} else {
					logrus.Error("Filter not set")
					return
				}
			},

			// Register a handler for any emitted error.
			ErrHandler: func(err error) {
				logrus.Error(err)
			},

			// Register a handler when a stream is completed.
			DoneHandler: func() {
				logrus.Info("Done")
				wg.Done()
			},
		}
		ob.Subscribe(obs)
	}
	wg.Wait()
}
