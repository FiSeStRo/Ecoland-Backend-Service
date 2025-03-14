package model

// Production represents a production process
type Production struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	InputType  []RelProduct `json:"inputs"`
	OutputType []RelProduct `json:"outputs"`
	Duration   int          `json:"duration"`
	Cost       float64      `json:"cost"`
}

type RelProduct struct {
	ProductID int `json:"productId"`
	Amount    int `json:"amount"`
}
