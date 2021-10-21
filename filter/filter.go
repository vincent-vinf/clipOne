package filter

import "clipOne/clipboard/cell"

// Filter Responsible link abstraction
type Filter interface {
	Execute(cell *cell.Cell) *cell.Cell
	SetNext(next Filter)
}
