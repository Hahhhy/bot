package image

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

const outputDir = "output"

func drawCatFace(dc *gg.Context) {
	dc.Push()
	// cat head - yellow circle
	dc.SetColor(color.RGBA{255, 220, 100, 255})
	dc.DrawCircle(300, 160, 80)
	dc.Fill()

	// left ear - triangle
	dc.MoveTo(240, 100)
	dc.LineTo(260, 50)
	dc.LineTo(290, 90)
	dc.ClosePath()
	dc.Fill()

	// right ear - triangle
	dc.MoveTo(310, 90)
	dc.LineTo(340, 50)
	dc.LineTo(360, 100)
	dc.ClosePath()
	dc.Fill()

	// inner ears (pink)
	dc.SetColor(color.RGBA{255, 150, 150, 255})
	dc.MoveTo(252, 92)
	dc.LineTo(262, 62)
	dc.LineTo(280, 90)
	dc.ClosePath()
	dc.Fill()
	dc.MoveTo(320, 90)
	dc.LineTo(338, 62)
	dc.LineTo(348, 92)
	dc.ClosePath()
	dc.Fill()

	// eyes - white circles
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.DrawEllipse(270, 150, 18, 22)
	dc.Fill()
	dc.DrawEllipse(330, 150, 18, 22)
	dc.Fill()

	// pupils - black circles
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(273, 148, 6)
	dc.Fill()
	dc.DrawCircle(333, 148, 6)
	dc.Fill()

	// nose - small pink triangle
	dc.SetColor(color.RGBA{255, 100, 100, 255})
	dc.MoveTo(300, 165)
	dc.LineTo(290, 175)
	dc.LineTo(310, 175)
	dc.ClosePath()
	dc.Fill()

	// mouth
	dc.SetRGB(0, 0, 0)
	dc.DrawLine(300, 175, 300, 185)
	dc.DrawLine(300, 185, 280, 195)
	dc.DrawLine(300, 185, 320, 195)
	dc.SetLineWidth(2)
	dc.Stroke()

	// whiskers
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1.5)
	w := []struct{ x1, y1, x2, y2 float64 }{
		{220, 155, 255, 160},
		{220, 165, 255, 168},
		{345, 160, 380, 155},
		{345, 168, 380, 165},
	}
	for _, v := range w {
		dc.DrawLine(v.x1, v.y1, v.x2, v.y2)
		dc.Stroke()
	}

	// pointing arm (a line + small paw)
	dc.SetColor(color.RGBA{255, 220, 100, 255})
	dc.SetLineWidth(12)
	dc.DrawLine(350, 200, 420, 240)
	dc.Stroke()
	// paw
	dc.SetLineWidth(3)
	for i := 0; i < 3; i++ {
		angle := -math.Pi/4 + float64(i)*math.Pi/8
		lx := 420 + 12*math.Cos(angle)
		ly := 240 + 12*math.Sin(angle)
		dc.DrawLine(420, 240, lx, ly)
		dc.Stroke()
	}
	dc.Pop()
}

func GenerateTom(name string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	const W, H = 600, 400
	dc := gg.NewContext(W, H)

	// background - light blue
	dc.SetColor(color.RGBA{200, 230, 255, 255})
	dc.Clear()

	// draw the cat
	drawCatFace(dc)

	// text
	dc.LoadFontFace("/usr/share/fonts/Adwaita/AdwaitaSans-Regular.ttf", 28)
	dc.SetRGB(0, 0, 0)
	text := fmt.Sprintf("%s 被汤姆嘲笑了！", name)
	dc.DrawStringAnchored(text, W/2, 380, 0.5, 0.5)

	filename := filepath.Join(outputDir, "tom.png")
	if err := dc.SavePNG(filename); err != nil {
		return "", err
	}
	return filename, nil
}
