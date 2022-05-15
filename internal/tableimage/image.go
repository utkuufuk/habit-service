package tableimage

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func (t *TableImage) drawMask() {
	draw.Draw(
		t.img,
		t.img.Bounds(),
		&image.Uniform{getColorByHex(t.BackgroundColor)},
		image.ZP,
		draw.Src,
	)
}

func (t *TableImage) drawHeader() {
	for colNo, td := range t.Header.Tds {
		t.addString(colNo*columnWidth+columnPadding, rowHeight, td.Text, td.Color)
		t.addLine(colNo*columnWidth, 0, colNo*columnWidth, t.height, "#000000")
	}

	t.addLine(1, rowHeight+rowBottomPadding, t.width, rowHeight+rowBottomPadding, "#000000")
	t.addLine(1, rowHeight+rowBottomPadding+2, t.width, rowHeight+rowBottomPadding+2, "#000000")
}

func (t *TableImage) drawRows() {
	fRowNo := 2
	for _, tds := range t.Rows {
		for colNo, td := range tds.Tds {
			colHeight := fRowNo * rowHeight
			colWidth := colNo*columnWidth + columnPadding
			if colNo == 0 {
				colWidth = columnWidth - columnWidth + columnPadding
			}

			t.addString(colWidth, colHeight, td.Text, td.Color)
		}

		t.addLine(1, fRowNo*rowHeight+rowBottomPadding, t.width, fRowNo*rowHeight+rowBottomPadding, "#000")
		fRowNo++
	}
}

func (t *TableImage) addString(x, y int, label string, color string) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  t.img,
		Src:  image.NewUniform(getColorByHex(color)),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func (t *TableImage) addLine(x1, y1, x2, y2 int, color string) {
	var dx, dy, e, slope int
	col := getColorByHex(color)
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	if dy < 0 {
		dy = -dy
	}

	switch {
	case x1 == x2 && y1 == y2:
		t.img.Set(x1, y1, col)

	case y1 == y2:
		for ; dx != 0; dx-- {
			t.img.Set(x1, y1, col)
			x1++
		}
		t.img.Set(x1, y1, col)

	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			t.img.Set(x1, y1, col)
			y1++
		}
		t.img.Set(x1, y1, col)

	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				t.img.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				t.img.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		t.img.Set(x1, y1, col)

	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				t.img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				t.img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		t.img.Set(x2, y2, col)

	default:
		if y1 < y2 {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				t.img.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				t.img.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		t.img.Set(x2, y2, col)
	}
}

func getColorByHex(hex string) color.RGBA {
	var r, g, b int
	a := 255

	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(hex, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}

	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	}

	if len(hex) == 8 {
		fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
