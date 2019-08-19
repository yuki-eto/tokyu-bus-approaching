package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/PuerkitoBio/goquery"
	tokyuBusApproaching "github.com/yuki-eto/tokyu-bus-approaching"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()
	if configPath != "" {
		if err := tokyuBusApproaching.LoadConfig(configPath); err != nil {
			panic(err)
		}
	}

	a := app.New()
	defer a.Quit()

	t := tokyuBusApproaching.NewCustomTheme()
	a.Settings().SetTheme(t)
	if err := t.LoadFont(); err != nil {
		panic(err)
	}

	w := a.NewWindow("Tokyu Bus Approach Info")
	defer w.Close()

	textStyle := fyne.TextStyle{}
	loadingLabel := widget.NewLabelWithStyle("Now loading...", fyne.TextAlignCenter, textStyle)
	vbox := widget.NewVBox(loadingLabel)
	w.SetContent(vbox)

	refreshTime := time.Duration(tokyuBusApproaching.GetConfig().RefreshTime) * time.Second
	ticker := time.NewTicker(refreshTime)
	go func() {
		for {
			vbox.Children = []fyne.CanvasObject{loadingLabel}
			widget.Refresh(vbox)

			log.Print("fetch approaching information...")
			infos, err := fetchInfos()
			if err != nil {
				log.Printf("error: %+v", err)
				continue
			}
			infos.SortByRemainTime()

			vbox.Children = []fyne.CanvasObject{}
			infos.Each(func(info *tokyuBusApproaching.ApproachInformation) {
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

	w.ShowAndRun()
}

func fetchInfos() (*tokyuBusApproaching.ApproachInformationList, error) {
	infos := tokyuBusApproaching.NewApproachInformationList()

	for _, url := range tokyuBusApproaching.GetConfig().URLs {
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
