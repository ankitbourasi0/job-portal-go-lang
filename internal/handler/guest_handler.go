package handler

import (
	"net/http"

	"github.com/ankitbourasi0/job-portal/internal/repository"
)

type GuestHandler struct {
	Repo *repository.GuestRepository
}

func (g GuestHandler) HandleAnalyzeResumeForAtsScore(w http.ResponseWriter, r *http.Request) {}
