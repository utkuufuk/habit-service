package tableimage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// TD a table data container
type TD struct {
	Text  string
	Color string
}

// TR the table row
type TR struct {
	BorderColor string
	Tds         []TD
}

type TableImage struct {
	BackgroundColor string
	FileType        string
	FilePath        string
	Header          TR
	Rows            []TR

	width  int
	height int
	img    *image.RGBA
}

const (
	rowBottomPadding = 10
	rowHeight        = 30
	columnPadding    = 16
	columnWidth      = 100
	Png              = "png"
	Jpeg             = "jpg"
)

func Draw(t TableImage) {
	// initialize private fields
	t.width = len(t.Header.Tds) * columnWidth
	t.height = (len(t.Rows)+2)*rowHeight - 2*rowBottomPadding
	t.img = image.NewRGBA(image.Rect(0, 0, t.width, t.height))

	t.drawMask()
	t.drawHeader()
	t.drawRows()
	t.saveFile()
}

func (t *TableImage) saveFile() error {
	f, err := os.Create(t.FilePath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer f.Close()

	if t.FileType == Jpeg {
		return jpeg.Encode(f, t.img, nil)
	}

	if t.FileType == Png {
		return png.Encode(f, t.img)
	}

	return fmt.Errorf("invalid image format: '%s'", t.FileType)
}
