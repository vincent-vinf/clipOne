package clipboard

import (
	"bytes"
	"clipOne/converter"
	"clipOne/deviceid"
	"context"
	"encoding/gob"
	"fmt"
	"golang.design/x/clipboard"

	"time"
)

const (
	TypeText = "text"
	TypeImg  = "image"
)

var (
	convert converter.Converter
)

type Cell struct {
	Time   time.Time
	Types  string
	Data   []byte
	Source deviceid.DeviceID
}

func init()  {
	convert = &converter.Compressor{}
}

func UseEncryptor(key []byte)  {
	e := &converter.Encryptor{
		Converter: convert,
	}
	e.SetKey(key)
	convert = e
}

func (c *Cell) Encode() ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(c)
	if err != nil {
		return nil, err
	}
	out, err := convert.Encode(buff.Bytes())
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Decode(data []byte) (*Cell, error) {
	data, err := convert.Decode(data)
	if err != nil {
		return nil, err
	}
	c := &Cell{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(c)
	return c, err
}

type Clipboard struct {
	CellChan chan *Cell
	backCell *Cell
}

func New() *Clipboard {
	return &Clipboard{
		CellChan: make(chan *Cell, 1),
	}
}

func (c *Clipboard) Write(cell *Cell) error {
	c.backCell = cell
	switch cell.Types {
	case TypeText:
		clipboard.Write(clipboard.FmtText, cell.Data)
		return nil
	case TypeImg:
		clipboard.Write(clipboard.FmtImage, cell.Data)
		return nil
	default:
		return fmt.Errorf("[ClipboardManager]:error type")
	}
}

func (c *Clipboard) Watching(ctx context.Context) {
	ch1 := clipboard.Watch(ctx, clipboard.FmtText)
	ch2 := clipboard.Watch(ctx, clipboard.FmtImage)

	for {
		select {
		case data := <-ch1:
			if c.backCell != nil && bytes.Equal(c.backCell.Data, data) {
				continue
			}
			c.CellChan <- &Cell{
				Time:   time.Now(),
				Types:  TypeText,
				Data:   data,
			}
		case data := <-ch2:
			if c.backCell != nil && bytes.Equal(c.backCell.Data, data) {
				continue
			}
			c.CellChan <- &Cell{
				Time:   time.Now(),
				Types:  TypeImg,
				Data:   data,
			}
		}
	}
}

func (c *Clipboard) Clean() {
	for len(c.CellChan) > 0 {
		<-c.CellChan
	}
}
