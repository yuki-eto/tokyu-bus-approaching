package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/martinlindhe/notify"
	tokyu_bus_approaching "github.com/yuki-eto/tokyu-bus-approaching"
)

const AppName = "TokyuBusNotification"
const TitleFormat = "あと%d分で%s行きのバスが来ます！"

func main() {
	ticker := time.NewTicker(time.Minute * 2)

	if err := run(); err != nil {
		panic(err)
	}

	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}
		if err := run(); err != nil {
			log.Printf("error: %+v", err)
		}
	}
}

func run() error {
	infos := tokyu_bus_approaching.NewApproachInformationList()

	urls := []string{
		// 渋谷駅（中目黒方面）行き
		"http://tokyu.bus-location.jp/blsys/navi?VID=rsc&EID=rd&FID=&SID=&DSMK=2536&ASMK=2351&ART=30&ARC=0&PGL=3&SCT=2",
		// 渋谷駅東口（恵比寿方面）行き
		"http://tokyu.bus-location.jp/blsys/navi?VID=lsc&EID=nt&FID=&SID=&PRM=&SCT=2&FDSN=0&FASN=0&DSMK=2536&ASMK=2440",
		// 東京駅南口（虎ノ門方面）行き
		"http://tokyu.bus-location.jp/blsys/navi?VID=lsc&EID=nt&FID=&SID=&PRM=&SCT=2&FDSN=0&FASN=0&DSMK=2536&ASMK=2481",
	}
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return errors.New(fmt.Sprintf("error code: %d", res.StatusCode))
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}
		res.Body.Close()

		var selections []*goquery.Selection
		doc.Find(".businfo").Each(func(i int, selection *goquery.Selection) {
			selections = append(selections, selection)
		})

		if err := infos.AppendBySelections(selections); err != nil {
			return err
		}
	}

	infos.Each(func(info *tokyu_bus_approaching.ApproachInformation) {
		log.Print(info)
		if info.RemainMinutes < 3 || info.RemainMinutes > 5 {
			return
		}
		notify.Notify(AppName, fmt.Sprintf(TitleFormat, info.RemainMinutes, info.Destination), "", "")
		time.Sleep(time.Second * 5)
	})

	return nil
}
