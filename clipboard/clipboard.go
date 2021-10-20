package clipboard

import (
	"bytes"
	"clipOne/clipboard/cell"
	"clipOne/converter"
	"context"
	"fmt"
	"golang.design/x/clipboard"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

const (
	TypeText = "text"
	TypeImg  = "image"
)

var (
	convert converter.Converter
)

func init() {
	convert = &converter.BaseConverter{}
}

func UseCompress() {
	convert = &converter.Compressor{Converter: convert}
}

func UseEncryptor(key []byte) {
	e := &converter.Encryptor{
		Converter: convert,
	}
	e.SetKey(key)
	convert = e
}

func Encode(c *cell.Cell) ([]byte, error) {
	data, err := proto.Marshal(c)
	if err != nil {
		return nil, err
	}
	out, err := convert.Encode(data)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Decode(data []byte) (*cell.Cell, error) {
	data, err := convert.Decode(data)
	if err != nil {
		return nil, err
	}
	c := &cell.Cell{}
	err = proto.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, err
}

type Clipboard struct {
	CellChan chan *cell.Cell
	backCell *cell.Cell
	rwLock   sync.RWMutex
}

func New() *Clipboard {
	return &Clipboard{
		CellChan: make(chan *cell.Cell, 1),
	}
}

func (c *Clipboard) Write(cell *cell.Cell) error {
	c.rwLock.Lock()
	c.backCell = cell
	c.rwLock.Unlock()
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
			c.CellChan <- &cell.Cell{
				Time:  timestamppb.Now(),
				Types: TypeText,
				Data:  data,
			}
		case data := <-ch2:
			c.rwLock.RLock()
			if c.backCell != nil && bytes.Equal(c.backCell.Data, data) {
				c.rwLock.RUnlock()
				continue
			}
			c.rwLock.RUnlock()
			c.CellChan <- &cell.Cell{
				Time:  timestamppb.Now(),
				Types: TypeImg,
				Data:  data,
			}
		}
	}
}

func (c *Clipboard) Clean() {
	for len(c.CellChan) > 0 {
		<-c.CellChan
	}
}
