package services

import (
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/viper"
)

const (
	fontSize = 76
	imageW   = 160
	imageH   = 160
)

type AvatarService interface {
	GenerateAvatar(word string) (image.Image, error)
}

type avatarService struct {
	font *truetype.Font
}

func NewAvatarService() AvatarService {
	// Read the font data.
	fontBytes, err := os.ReadFile(viper.GetString("font"))
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	return &avatarService{
		font: font,
	}
}

func colorPalette(seed string) (bg color.RGBA, ring color.RGBA) {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(seed))
	sum := hasher.Sum32()

	base := uint8(sum % 120)
	bg = color.RGBA{
		R: uint8(120 + base/3),
		G: uint8(145 + base/4),
		B: uint8(185 + base/5),
		A: 255,
	}
	ring = color.RGBA{
		R: uint8(40 + base/4),
		G: uint8(75 + base/3),
		B: uint8(120 + base/2),
		A: 255,
	}
	return
}

func (s *avatarService) GenerateAvatar(word string) (image.Image, error) {
	word = strings.TrimSpace(word)
	if word == "" {
		word = "U"
	}

	first, _ := utf8.DecodeRuneInString(word)
	char := strings.ToUpper(string(first))
	bg, ring := colorPalette(word)

	m := image.NewRGBA(image.Rect(0, 0, imageW, imageH))
	draw.Draw(m, m.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Draw a subtle circle so avatar feels intentional on modern UI.
	cx, cy := imageW/2, imageH/2
	radius := imageW / 2
	inner := radius - 6
	for y := 0; y < imageH; y++ {
		for x := 0; x < imageW; x++ {
			dx := x - cx
			dy := y - cy
			d2 := dx*dx + dy*dy
			if d2 <= radius*radius && d2 >= inner*inner {
				m.Set(x, y, ring)
			}
		}
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(s.font)
	c.SetFontSize(fontSize)
	c.SetClip(m.Bounds())
	c.SetDst(m)
	c.SetSrc(image.White)

	pt := freetype.Pt(44, 106)
	_, err := c.DrawString(char, pt)
	if err != nil {
		log.Printf("GenerateAvatar error: %v", err)
		return nil, err
	}

	return m, nil
}
