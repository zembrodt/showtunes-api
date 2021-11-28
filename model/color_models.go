package model

type ColorResponse struct {
	Color                string `json:"color" example:"#1500ff"`
	UseBlackBarcodeColor bool   `json:"useBlackBarcodeColor" example:"true"`
}
