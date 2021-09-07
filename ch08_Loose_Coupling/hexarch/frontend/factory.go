package frontend

import (
	"fmt"
)

//Factoria que retorna uno u otro frontend
func NewFrontEnd(s string) (FrontEnd, error) {
	switch s {
	case "zero":
		return zeroFrontEnd{}, nil

	case "rest":
		return &restFrontEnd{}, nil

	default:
		return nil, fmt.Errorf("no such frontend %s", s)
	}
}
