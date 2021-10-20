package main

import (
	"clipOne/clipboard"
	"clipOne/deviceid"
	"clipOne/mq"
	"context"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	sec := cfg.Section(ini.DefaultSection)
	user, err := sec.GetKey("user")
	if err != nil {
		log.Fatal(err)
	}
	pass, err := sec.GetKey("password")
	if err != nil {
		log.Fatal(err)
	}
	url, err := sec.GetKey("url")
	if err != nil {
		log.Fatal(err)
	}

	clipboard.UseCompress()
	clipboard.UseEncryptor(deviceid.MD5([]byte(pass.Value())))

	clipboardCtx, clipboardCancel := context.WithCancel(context.Background())
	clipboardM := clipboard.New()

	go clipboardM.Watching(clipboardCtx)
	defer clipboardCancel()

	reconnectCh := make(chan struct{}, 1)
	reconnectCh <- struct{}{}
	errCh := make(chan error, 1)

	var (
		msgManager *mq.MsgManager
		cancel     func()
		ctx        context.Context
	)

	for {
		select {
		case <-reconnectCh:
			msgManager = mq.NewMsgManager(user.Value(), url.Value())
			cancel = nil
			err := msgManager.Init()
			if err != nil {
				log.Println("init msg fail")
				msgManager = nil
				errCh <- err
				continue
			}

			ctx, cancel = context.WithCancel(context.Background())
			err = msgManager.Receive(ctx)
			if err != nil {
				log.Println("init msg receive err")
				errCh <- err
				continue
			}

			log.Println("rabbitMQZ init success")
		Deal:
			for {
				select {
				case cell := <-clipboardM.CellChan:
					data, err := cell.Encode()
					if err != nil {
						log.Println("encode fail: ", err)
						continue
					}
					err = msgManager.Send(data)
					if err != nil {
						log.Println("send fail: ", err)
						continue
					}

					log.Println("send: ", cell.Time)
				case data := <-msgManager.ReceiveCh:
					if len(data) == 0 {
						log.Println("empty payload, reconnect")
						errCh <- fmt.Errorf("empty payload")
						break Deal
					}
					c, err := clipboard.Decode(data)
					if err != nil {
						log.Println("decode fail: ", err)
						continue
					}
					//log.Println("delay: ", time.Now().Sub(c.Time))
					err = clipboardM.Write(c)
					if err != nil {
						log.Println("write fail: ", err)
						continue
					}
				}
			}
		case err := <-errCh:
			log.Println("err: ", err)
			if cancel != nil {
				cancel()
			}
			if msgManager != nil {
				msgManager.Close()
			}
			clipboardM.Clean()
			reconnectCh <- struct{}{}
			<-time.After(time.Second * 3)
		}
	}
}
