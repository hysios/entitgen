package out

import (
	"database/sql"
	"time"

	pb "github.com/hysios/entitgen/example/gen/proto"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Member struct {
	ID           uint
	UserID       uint
	EnterpriseID uint
	User         *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//go:generate entitgen -type User -O no_models=User
type User struct {
	ID          uint
	Name        string
	Username    string
	Namespace   string
	Nickname    string
	Email       string
	Password    sql.NullString
	Avatar      string
	Phone       string
	Address     string
	Description string
	Score       float64
	Role        string
	Age         uint64
	Money       decimal.Decimal `json:"money" gorm:"type:decimal(10,2)"`
	PersonID    sql.NullInt32
	IsActive    bool
	Period      time.Duration
	ExpiresAt   sql.NullTime
	InScopes    datatypes.JSONSlice[string]
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiredAt   sql.NullTime
	MemberID    uint
	Member      *Member
	Leader      datatypes.JSONType[*pb.Member]
	Friends     []*Friend
}

type Agent struct {
	Agentid          int    `json:"agentid"`
	Name             string `json:"name"`
	RoundLogoURL     string `json:"roundLogoUrl"`
	SquareLogoURL    string `json:"squareLogoUrl"`
	AuthMode         int    `json:"authMode"`
	IsCustomizedApp  bool   `json:"isCustomizedApp"`
	AuthFromThirdapp bool   `json:"authFromThirdapp"`
	// Privilege        datatypes.JSONType[Privilege]  `json:"privilege"`
	// SharedFrom       datatypes.JSONType[SharedFrom] `json:"sharedFrom"`
}
