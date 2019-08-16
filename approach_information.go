package tokyu_bus_approaching

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var busInfoRegexp = regexp.MustCompile(`([^\s]+)\x{00a0}(.+)行【([0-9]{2})分待ち】`)

type ApproachInformation struct {
	Routes        string
	Destination   string
	RemainMinutes uint32
}

type ApproachInformationList struct {
	values []*ApproachInformation
}

func NewApproachInformationList() *ApproachInformationList {
	return &ApproachInformationList{
		values: []*ApproachInformation{},
	}
}

func (a *ApproachInformationList) Append(info *ApproachInformation) {
	a.values = append(a.values, info)
}
func (a *ApproachInformationList) AppendBySelections(selections []*goquery.Selection) error {
	for _, selection := range selections {
		busInfo := strings.TrimSpace(selection.Text())
		if busInfo == "" {
			continue
		}
		matches := busInfoRegexp.FindStringSubmatch(busInfo)
		if len(matches) < 4 {
			continue
		}
		routes, destination := matches[1], matches[2]
		remainMin, err := strconv.ParseUint(matches[3], 10, 32)
		if err != nil {
			return err
		}
		a.Append(&ApproachInformation{
			Routes:        routes,
			Destination:   destination,
			RemainMinutes: uint32(remainMin),
		})
	}
	return nil
}
func (a *ApproachInformationList) Each(f func(info *ApproachInformation)) {
	for _, i := range a.values {
		f(i)
	}
}
func (a *ApproachInformationList) SortByRemainTime() {
	sort.Slice(a.values, func(i, j int) bool {
		return a.values[i].RemainMinutes < a.values[j].RemainMinutes
	})
}
