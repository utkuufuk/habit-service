package tableimage

import (
	"image"
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
	width           int
	height          int
	th              TR
	trs             []TR
	backgroundColor string
	fileType        string
	filePath        string
	img             *image.RGBA
}

const (
	rowSpace         = 26
	tablePadding     = 20
	letterPerPx      = 10
	separatorPadding = 10
	wrapWordsLen     = 12
	columnSpace      = wrapWordsLen * letterPerPx
	Png              = "png"
	Jpeg             = "jpg"
)

// Init initialise the table image receiver
func Init(backgroundColor string, fileType string, filePath string) TableImage {
	ti := TableImage{
		backgroundColor: backgroundColor,
		fileType:        fileType,
		filePath:        filePath,
	}
	ti.setRgba()
	return ti
}

// AddTH adds the table header
func (ti *TableImage) AddTH(th TR) {
	ti.th = th
}

// AddTRs add the table rows
func (ti *TableImage) AddTRs(trs []TR) {
	ti.trs = trs
}

// Save saves the table image
func (ti *TableImage) Save() {
	ti.calculateHeight()
	ti.calculateWidth()
	ti.setRgba()
	ti.drawTH()
	ti.drawTR()
	ti.saveFile()
}
