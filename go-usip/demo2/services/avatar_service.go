package services

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/viper"
)

const (
	fontSize = 400
	imageW   = 640
	imageH   = 480
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

func randomColor(k int16) color.Color {
	loc := uint8(k & 0x03)
	switch loc {
	case 0:
		return color.RGBA{
			R: uint8(k & 0xff),
			G: 0,
			B: 0,
			A: 255,
		}
	case 1:
		return color.RGBA{
			R: 0,
			G: uint8(k & 0xff),
			B: 0,
			A: 255,
		}
	default:
		return color.RGBA{
			R: 0,
			G: 0,
			B: uint8(k & 0xff),
			A: 255,
		}
	}
}

func (s *avatarService) GenerateAvatar(word string) (image.Image, error) {
	wordB := []byte(word)
	rand := int16(wordB[0])
	if len(wordB) > 1 {
		rand += int16(wordB[1])
	}
	m := image.NewRGBA(image.Rect(0, 0, imageW, imageH))
	draw.Draw(m, m.Bounds(), &image.Uniform{randomColor(rand)}, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(s.font)
	c.SetFontSize(fontSize)
	c.SetClip(m.Bounds())
	c.SetDst(m)
	c.SetSrc(image.White)

	pt := freetype.Pt(int(c.PointToFixed((imageW-fontSize)/2)>>6)+10, int(c.PointToFixed(fontSize)>>6)-20)
	_, err := c.DrawString(word, pt)
	if err != nil {
		log.Printf("GenerateAvatar error: %v", err)
		return nil, err
	}

	return m, nil
}
