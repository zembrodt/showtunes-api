package model

type ColorResponse struct {
	Color        string `json:"color" example:"#1500ff"`
	BarCodeColor string `json:"barCodeColor" example:"white,black"`
}
