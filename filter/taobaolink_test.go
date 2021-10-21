package filter

import (
	"clipOne/clipboard"
	"clipOne/clipboard/cell"
	"log"
	"testing"
)

func TestTaobaoLink_Execute(t *testing.T) {
	l := &TaobaoLink{}
	c := &cell.Cell{
		Time:  nil,
		Types: clipboard.TypeText,
		Data:  []byte("1.0SxNtXtwblzt https://m.tb.cn/h.fVGnWH9?sm=ff9bc2  【官方旗舰】罗技 ERGO M575无线轨迹球"),
	}
	if l.Execute(c) == nil {
		log.Println("nil")
	}
}
