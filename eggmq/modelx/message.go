package modelx

type MessageM struct {
	Basic
	ID         string `gorm:"column:id;primaryKey;not null;size:100"`
	Topic      string `gorm:"column:topic;index;not null;size:200"`
	Payload    string `gorm:"type:text"`
	RetryCount int    `gorm:"column:retry_count;default:0"`
	Status     int    `gorm:"column:status;"`
}

func (model *MessageM) TableName() string {
	return "eggmq_message"
}
