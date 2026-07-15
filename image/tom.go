package image

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

const outputDir = "output"

func GenerateTom(name string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	const W, H = 600, 400
	dc := gg.NewContext(W, H)

	dc.SetColor(color.RGBA{255, 200, 100, 255})
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace("/usr/share/fonts/Adwaita/AdwaitaSans-Regular.ttf", 36)
	dc.DrawStringAnchored(fmt.Sprintf("%s", name), W/2, H/2-20, 0.5, 0.5)
	dc.DrawStringAnchored("被汤姆嘲笑了！", W/2, H/2+30, 0.5, 0.5)

	filename := filepath.Join(outputDir, "tom.png")
	if err := dc.SavePNG(filename); err != nil {
		return "", err
	}
	return filename, nil
}
