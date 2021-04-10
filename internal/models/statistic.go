package models

type Statistic struct {
	Zero   int64 `json:"zero" bson:"zero"`
	Low    int64 `json:"low" bson:"low"`
	Hard   int64 `json:"hard" bson:"hard"`
	Driver int64 `json:"driver" bson:"driver"`
	Total  int64 `json:"total" bson:"total"`
}
