package configmodel

type GetConfig struct {
	LatestBlockNumber int64 `json:"latest_block_number" gorm:"column:latest_block_number;"`
}

func (GetConfig) TableName() string {
	return "config"
}
