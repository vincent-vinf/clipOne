package filter

import (
	"clipOne/clipboard"
	"clipOne/clipboard/cell"
	"log"
	"regexp"
)

type VerificationCode struct {
	next   Filter
	regexp *regexp.Regexp
}

func NewVerificationCode() *VerificationCode {
	r, err := regexp.Compile("^[0-9]{4,6}$")
	if err != nil {
		panic(err)
	}
	return &VerificationCode{
		regexp: r,
	}
}

func (c *VerificationCode) Execute(cell *cell.Cell) *cell.Cell {
	if cell.Types == clipboard.TypeText && c.regexp.Match(cell.Data) {
		log.Println("match code")
		return nil
	}

	if c.next == nil {
		return cell
	}
	// next
	return c.next.Execute(cell)
}

func (c *VerificationCode) SetNext(next Filter) {
	c.next = next
}
