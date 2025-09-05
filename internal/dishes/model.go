package dishes

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Dish struct {
	ID          uuid.UUID       `json:"id" gorm:"default:gen_random_uuid();primary key"`
	Name        string          `json:"name" gorm:"varchar(100);not null;unique"`
	Description string          `json:"description" gorm:"text;not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:numeric(10,2);not null"`
	Category    string          `json:"category" gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}
