package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRep struct {
	db *pgxpool.Pool
}
