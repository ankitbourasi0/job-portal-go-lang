package repository

import "github.com/ankitbourasi0/job-portal/internal/database"

type GuestRepository struct {
	Queries *database.Queries
}

func NewGuestRepository(dbQueries *database.Queries) *GuestRepository {
	return &GuestRepository{Queries: dbQueries}
}
