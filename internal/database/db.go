package database

import (
	"database/sql"
	"fmt"

	"github.com/EduardoMark/gastro-api/internal/config"
	"github.com/EduardoMark/gastro-api/internal/dishes"
	"github.com/EduardoMark/gastro-api/internal/order"
	"github.com/EduardoMark/gastro-api/internal/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(env *config.Env) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DbUser,
		env.DbPassword,
		env.DbHost,
		env.DbPort,
		env.DbName,
	)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open pgx connection: %v", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		users.User{},
		dishes.Dish{},
		order.Order{},
		order.OrderItem{},
	)
}
