package models

type AttributeRow struct {
	Name   string `json:"name"`
	Info   string `json:"info"`
	Suffix string `json:"suffix"`
}

type Attribute struct {
	Phones []AttributeRow `json:"phones"`
	Urls   []AttributeRow `json:"urls"`
}

type FieldRow struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
