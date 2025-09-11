package order

import (
	"time"

	"github.com/EduardoMark/gastro-api/internal/dishes"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	STATUS_NEW            Status = "new"
	STATUS_IN_PREPARATION Status = "in preparation"
	STATUS_FINISHED       Status = "finished"
)

type Order struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Status      Status          `json:"status" gorm:"type:varchar(100);not null"`
	TotalAmount decimal.Decimal `json:"total_amount" gorm:"type:numeric"`
	Items       []OrderItem     `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID       uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	OrderID  uuid.UUID       `json:"order_id" gorm:"type:uuid;not null"`
	DishID   uuid.UUID       `json:"dish_id" gorm:"type:uuid;not null"`
	Quantity int             `json:"quantity"`
	Price    decimal.Decimal `json:"price" gorm:"type:numeric"`
	SubTotal decimal.Decimal `json:"sub_total" gorm:"type:numeric"`
	Dish     dishes.Dish     `json:"dish" gorm:"foreignKey:DishID;references:ID"`
}
