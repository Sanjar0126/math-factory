package fonts

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var MplusNormalFont *text.GoTextFace = nil
var mplusFaceSource *text.GoTextFaceSource

func InitFonts() error {
	var err error

	if MplusNormalFont == nil {
		s, err := text.NewGoTextFaceSource(bytes.NewReader(MPlus1pRegular_ttf))
		if err != nil {
			log.Fatal(err)
		}
		mplusFaceSource = s

		MplusNormalFont = &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   12,
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}
