package tokyu_bus_approaching

import (
	"image/color"
	"io/ioutil"
	"log"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/rakyll/statik/fs"
	_ "github.com/yuki-eto/tokyu-bus-approaching/statik"
)

type CustomTheme struct {
	baseTheme fyne.Theme
	font      fyne.Resource
}

func NewCustomTheme() *CustomTheme {
	base := theme.DarkTheme()
	return &CustomTheme{
		baseTheme: base,
		font:      nil,
	}
}
func (c *CustomTheme) LoadFont() error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	name := "/migmix-1m.ttf"
	f, err := statikFS.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	res := &fyne.StaticResource{
		StaticName:    path.Base(name),
		StaticContent: b,
	}
	log.Printf("loading font: %s [%dbytes]", res.Name(), len(b))

	c.font = res
	return nil
}

func (c *CustomTheme) BackgroundColor() color.Color {
	return c.baseTheme.BackgroundColor()
}

func (c *CustomTheme) ButtonColor() color.Color {
	return c.baseTheme.ButtonColor()
}

func (c *CustomTheme) DisabledButtonColor() color.Color {
	return c.baseTheme.DisabledButtonColor()
}

func (c *CustomTheme) HyperlinkColor() color.Color {
	return c.baseTheme.HyperlinkColor()
}

func (c *CustomTheme) TextColor() color.Color {
	return c.baseTheme.TextColor()
}

func (c *CustomTheme) DisabledTextColor() color.Color {
	return c.baseTheme.DisabledTextColor()
}

func (c *CustomTheme) IconColor() color.Color {
	return c.baseTheme.IconColor()
}

func (c *CustomTheme) DisabledIconColor() color.Color {
	return c.baseTheme.DisabledIconColor()
}

func (c *CustomTheme) PlaceHolderColor() color.Color {
	return c.baseTheme.PlaceHolderColor()
}

func (c *CustomTheme) PrimaryColor() color.Color {
	return c.baseTheme.PrimaryColor()
}

func (c *CustomTheme) HoverColor() color.Color {
	return c.baseTheme.HoverColor()
}

func (c *CustomTheme) FocusColor() color.Color {
	return c.baseTheme.FocusColor()
}

func (c *CustomTheme) ScrollBarColor() color.Color {
	return c.baseTheme.ScrollBarColor()
}

func (c *CustomTheme) ShadowColor() color.Color {
	return c.baseTheme.ShadowColor()
}

func (c *CustomTheme) TextSize() int {
	return c.baseTheme.TextSize()
}

func (c *CustomTheme) TextFont() fyne.Resource {
	return c.font
}

func (c *CustomTheme) TextBoldFont() fyne.Resource {
	return c.font
}

func (c *CustomTheme) TextItalicFont() fyne.Resource {
	return c.font
}

func (c *CustomTheme) TextBoldItalicFont() fyne.Resource {
	return c.font
}

func (c *CustomTheme) TextMonospaceFont() fyne.Resource {
	return c.font
}

func (c *CustomTheme) Padding() int {
	return c.baseTheme.Padding()
}

func (c *CustomTheme) IconInlineSize() int {
	return c.baseTheme.IconInlineSize()
}

func (c *CustomTheme) ScrollBarSize() int {
	return c.baseTheme.ScrollBarSize()
}

func (c *CustomTheme) ScrollBarSmallSize() int {
	return c.baseTheme.ScrollBarSmallSize()
}
