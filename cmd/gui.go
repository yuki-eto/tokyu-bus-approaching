package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/PuerkitoBio/goquery"
	tokyu_bus_approaching "github.com/yuki-eto/tokyu-bus-approaching"
)

func main() {
	a := app.New()
	defer a.Quit()

	t := tokyu_bus_approaching.NewCustomTheme()
	a.Settings().SetTheme(t)
	if err := t.LoadFont(); err != nil {
		panic(err)
	}

	w := a.NewWindow("Tokyu Bus Approach Info")
	defer w.Close()

	vbox := widget.NewVBox(widget.NewLabel("now loading..."))
	w.SetContent(vbox)

	ticker := time.NewTicker(1 * time.Minute)
	textStyle := fyne.TextStyle{
		Monospace: true,
	}
	go func() {
		for {
			vbox.Children = []fyne.CanvasObject{}

			log.Print("fetch approaching information...")
			infos, err := fetchInfos()
			if err != nil {
				log.Printf("error: %+v", err)
				continue
			}
			infos.SortByRemainTime()

			infos.Each(func(info *tokyu_bus_approaching.ApproachInformation) {
				label := fmt.Sprintf("あと%d分 >> %s系統/%s行", info.RemainMinutes, info.Routes, info.Destination)
				labelWidget := widget.NewLabelWithStyle(label, fyne.TextAlignCenter, textStyle)
				vbox.Append(labelWidget)
			})
			if len(vbox.Children) == 0 {
				vbox.Append(widget.NewLabelWithStyle("接近情報はありません", fyne.TextAlignCenter, textStyle))
			}
			<-ticker.C
		}
	}()

	w.Show()
	a.Run()
}

func fetchInfos() (*tokyu_bus_approaching.ApproachInformationList, error) {
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
		doc, err := makeDoc(url)
		if err != nil {
			return nil, err
		}

		var selections []*goquery.Selection
		doc.Find(".businfo").Each(func(i int, selection *goquery.Selection) {
			selections = append(selections, selection)
		})

		if err := infos.AppendBySelections(selections); err != nil {
			return nil, err
		}
	}

	return infos, nil
}

func makeDoc(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("error code: %d", res.StatusCode))
	}
	return goquery.NewDocumentFromReader(res.Body)
}
