package filter

import (
	"bytes"
	"clipOne/clipboard"
	"clipOne/clipboard/cell"
	"net/url"
	"strings"
)

type TaobaoLink struct {
	next Filter
}

func (t *TaobaoLink) Execute(cell *cell.Cell) *cell.Cell {
	if cell.Types == clipboard.TypeText && isTaobaoLink(cell.Data) {
		cell = nil
		return nil
	}
	// tail
	if t.next == nil {
		return cell
	}
	// next
	return t.next.Execute(cell)
}
func (t *TaobaoLink) SetNext(next Filter) {
	t.next = next
}

func isTaobaoLink(str []byte) bool {
	if len(str) == 0 {
		return false
	}
	arr := bytes.Split(str, []byte(" "))
	if len(arr) >= 3 && len(arr[0]) > 0 && isLink(arr[1]) {
		return true
	}
	return false
}

func isLink(str []byte) bool {
	u, err := url.Parse(string(str))
	if err != nil {
		return false
	}
	if strings.EqualFold(u.Host, "m.tb.cn") {
		return true
	}
	return false
}
