package model

// Production represents a production process
type Production struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	InputType  []RelProduct `json:"inputType"`
	OutputType []RelProduct `json:"outputType"`
	Duration   int          `json:"duration"`
	Cost       float64      `json:"cost"`
}

type RelProduct struct {
	ProductID int `json:"productId"`
	Amount    int `json:"amount"`
}
