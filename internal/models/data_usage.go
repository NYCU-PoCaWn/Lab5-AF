package models

type RatingGroupDataUsage struct {
	Supi      string `bson:"Supi" json:"supi"`
	Filter    string `bson:"Filter" json:"filter,omitempty"`
	Snssai    string `bson:"Snssai" json:"snssai,omitempty"`
	Dnn       string `bson:"Dnn" json:"dnn,omitempty"`
	TotalVol  int64  `bson:"TotalVol" json:"totalVol,omitempty"`
	UlVol     int64  `bson:"UlVol" json:"ulVol,omitempty"`
	DlVol     int64  `bson:"DlVol" json:"dlVol,omitempty"`
	QuotaLeft int64  `bson:"QuotaLeft" json:"quotaLeft,omitempty"`
	Usage     int64  `bson:"Usage" json:"usage,omitempty"`
}
