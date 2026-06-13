package click

import (
	"strings"
)

// Addressed интерфейс предоставляет ID клика в трекере и адрес по которому можно определить трекер.
type Addressed interface {
	ID() string
	To() []string
	String() string
}

// NewAddress создает из clickID формата {to}_{id}.
func NewAddress(clickID string) Addressed {
	return AddressedClick(clickID)
}

type addressedClick []string

// AddressedClick разбирает clickID и создает *addressedClick.
func AddressedClick(clickID string) *addressedClick {
	var aclk addressedClick

	aclk = strings.Split(clickID, "_")
	return &aclk
}

// ID клика в трекере.
func (aclk *addressedClick) ID() string {
	if len(*aclk) < 1 {
		return ""
	}

	return (*aclk)[len(*aclk)-1]
}

// To возвращает идентификатор трекера из которого пришел клик.
// Либо возвращает nil в случае отсутствия адреса to.
func (aclk *addressedClick) To() []string {
	if len(*aclk) < 2 {
		return nil
	}

	return (*aclk)[:len(*aclk)-1]
}

// String возвращает строковое представление клика вида {to}_{id}.
func (aclk *addressedClick) String() string {
	return strings.Join(*aclk, "_")
}
