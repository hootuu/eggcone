package tick

import (
	"github.com/hootuu/eggcone/fdn/tick/schedule"
	"github.com/hootuu/eggcone/fdn/tick/token"
	"time"
)

type Tick struct {
	ID               ID        `gorm:"column:id;not null;size:32"`
	Server           string    `gorm:"column:server;not null;size:64"`
	LstHeartbeatTime time.Time `gorm:"column:lst_heartbeat_time"`
	Living           bool      `gorm:"column:living"`

	Version uint64 `gorm:"column:version"`

	SeqIdx     int64     `gorm:"column:seq_idx"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	ModifiedAt time.Time `gorm:"column:modified_at"`
}

func (m *Tick) TableName() string {
	return "egg_fdn_tick_tick"
}

// Bind 定时器和Token绑定关系
// Bind UK: Token TickID
// Bind UK: Token Available //同一时间只能有一个激活者
type Bind struct {
	Token     token.Token `gorm:"column:token;not null;size:32"`
	TickID    ID          `gorm:"column:tick_id;not null;size:32"`
	BindTime  time.Time   `gorm:"column:bind_time"`
	Available bool        `gorm:"column:available"`

	Version    int64     `gorm:"column:version"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	ModifiedAt time.Time `gorm:"column:modified_at"`
}

func (m *Bind) TableName() string {
	return "egg_fdn_tick_bind"
}

type Record struct {
	ScheduleID schedule.ID `gorm:"column:schedule_id"`
	DealTime   time.Time   `gorm:"column:deal_time"`
	Result     bool        `gorm:"column:result"`
	Listener   string      `gorm:"column:listener"`
	Ctx        any         `gorm:"column:ctx;type:json"`
}

func (m *Record) TableName() string { return "egg_fdn_tick_record" }
