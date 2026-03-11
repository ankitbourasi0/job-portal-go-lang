package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ankitbourasi0/job-portal/internal/database"
	"github.com/ankitbourasi0/job-portal/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobHandler struct {
	Repo *repository.JobRepository
}

type RequestBody struct {
	Title         string `json:"title"`
	JobRole       string `json:"job_role"`
	JobCategory   string `json:"job_category"`
	CompanyName   string `json:"company_name"`
	Location      string `json:"location"`
	Salary        string `json:"salary"`
	Qualification string `json:"qualification "`
	Experience    string `json:"experience"`
	LastDate      string `json:"last_date"`
	Description   string `json:"description"`
	ApplyUrl      string `json:"apply_url"`
}

type RequestBodyPointer struct {
	Title         *string `json:"title"`
	JobRole       *string `json:"job_role"`
	JobCategory   *string `json:"job_category"`
	CompanyName   *string `json:"company_name"`
	Location      *string `json:"location"`
	Salary        *string `json:"salary"`
	Qualification *string `json:"qualification"`
	Experience    *string `json:"experience"`
	LastDate      *string `json:"last_date"`
	Description   *string `json:"description"`
	ApplyUrl      *string `json:"apply_url"`
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	//1. request body for struct
	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newID := uuid.New()
	var pgID pgtype.UUID
	pgID.Bytes = newID
	pgID.Valid = true

	//2. prepare SQLC params: should match Models.go
	job, err := h.Repo.CreateNewJob(r.Context(), database.CreateJobParams{
		ID:          pgID,
		CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		Title:       req.Title,
		JobRole:     req.JobRole,
		JobCategory: pgtype.Text{String: req.JobCategory, Valid: true},
		CompanyName: req.CompanyName,
		Location:    req.Location,
		Salary:      pgtype.Text{String: req.Salary, Valid: true},

		Qualification: req.Qualification,
		Experience:    req.Experience,
		LastDate:      req.LastDate,
		Description:   req.Description,
		ApplyUrl:      req.ApplyUrl,
	})

	if err != nil {
		http.Error(w, "Failed to create job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Success Response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)

	if err != nil {
		http.Error(w, "Failed to Encode", http.StatusInternalServerError)

		return
	}

}

func (h *JobHandler) HandleGetAllJob(res http.ResponseWriter, req *http.Request) {
	//Get Data from Repository
	jobs, err := h.Repo.GetAllJobs(req.Context())
	if err != nil {
		http.Error(res, "Failed to get jobs", http.StatusInternalServerError)
	}

	//Convert Response in to JSON
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(jobs); err != nil {
		http.Error(res, "Failed to encode jobs", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleGetJobById(res http.ResponseWriter, req *http.Request) {
	//GET ID from URL PARAMS
	IdInString := chi.URLParam(req, "id")
	if IdInString == "" {
		http.Error(res, "ID is required", http.StatusBadRequest)
		return
	}
	//Convert URL PARAM(String) ID into PgUUID , PG expect PG UUID
	var pgID pgtype.UUID
	err := pgID.Scan(IdInString) //scan fn checks -  string in UUID format
	if err != nil {
		http.Error(res, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	//Get Job by ID from Repository
	job, err := h.Repo.GetJobById(req.Context(), pgID)

	//Validate Job
	if err != nil {
		http.Error(res, "Failed to get job by id", http.StatusInternalServerError)
	}

	//Encode Job in JSON
	res.Header().Set("Content-Type", "application/json")

	//Return
	if err := json.NewEncoder(res).Encode(job); err != nil {
		http.Error(res, "Failed to encode job", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleUpdateJobById(res http.ResponseWriter, req *http.Request) {

	//Get ID from URL Params
	IdInString := chi.URLParam(req, "id")
	if IdInString == "" {
		http.Error(res, "ID is required", http.StatusBadRequest)
		return
	}
	//Parse ID - String to UUID
	var pgID pgtype.UUID
	err := pgID.Scan(IdInString)
	//Validate ID
	if err != nil {
		http.Error(res, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	//Request Body - Updated Data
	var updatedRequest RequestBody
	//Parse JSON data into Struct
	if err := json.NewDecoder(req.Body).Decode(&updatedRequest); err != nil {
		http.Error(res, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Pass Request body Struct to SQLC Struct
	params := database.UpdateJobByIdParams{
		ID:      pgID,
		Title:   updatedRequest.Title,
		JobRole: updatedRequest.JobRole,
		JobCategory: pgtype.Text{
			String: updatedRequest.JobCategory,
		},
		CompanyName: updatedRequest.CompanyName,
		Location:    updatedRequest.Location,
		Salary: pgtype.Text{
			String: updatedRequest.Salary,
			Valid:  true, //so in db salary wont be null
		},
		Qualification: updatedRequest.Qualification,
		Experience:    updatedRequest.Experience,
		LastDate:      updatedRequest.LastDate,
		Description:   updatedRequest.Description,
		ApplyUrl:      updatedRequest.ApplyUrl,
	}

	//Pass updated params , Update Job by ID
	job, err := h.Repo.Queries.UpdateJobById(req.Context(), params)
	if err != nil {
		http.Error(res, "Failed to update job: "+err.Error(), http.StatusInternalServerError)

		return
	}
	//Convert Response into JSON
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(job); err != nil {
		http.Error(res, "Failed to encode job", http.StatusInternalServerError)
		return
	}

}
func (h *JobHandler) HandlePartialUpdateJob(res http.ResponseWriter, req *http.Request) {
	// 1. Decode Request Body (using pointers to detect missing fields)
	var request RequestBodyPointer

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Parse ID from URL params
	idInString := chi.URLParam(req, "id")
	var pgID pgtype.UUID
	err := pgID.Scan(idInString)
	if err != nil {
		http.Error(res, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	//Map pointers to pgtype.TEXT
	//if pointer is nil,then Valid=false(SQL COALESCE maintain old data)
	params := database.PartialUpdateJobParams{
		ID:            pgID,
		Title:         pgtype.Text{String: getString(request.Title), Valid: request.Title != nil},
		JobRole:       pgtype.Text{String: getString(request.JobRole), Valid: request.JobRole != nil},
		JobCategory:   pgtype.Text{String: getString(request.JobCategory), Valid: request.JobCategory != nil},
		CompanyName:   pgtype.Text{String: getString(request.CompanyName), Valid: request.CompanyName != nil},
		Location:      pgtype.Text{String: getString(request.Location), Valid: request.Location != nil},
		Salary:        pgtype.Text{String: getString(request.Salary), Valid: request.Salary != nil},
		Qualification: pgtype.Text{String: getString(request.Qualification), Valid: request.Qualification != nil},
		Experience:    pgtype.Text{String: getString(request.Experience), Valid: request.Experience != nil},
		LastDate:      pgtype.Text{String: getString(request.LastDate), Valid: request.LastDate != nil},
		Description:   pgtype.Text{String: getString(request.Description), Valid: request.Description != nil},
		ApplyUrl:      pgtype.Text{String: getString(request.ApplyUrl), Valid: request.ApplyUrl != nil},
	}

	//Call Repository
	updatedJob, err := h.Repo.PartialUpdateJob(req.Context(), params)
	if err != nil {
		http.Error(res, "Failed to update job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Response
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(updatedJob)
	if err != nil {
		http.Error(res, "Failed to encode job", http.StatusInternalServerError)
		return
	}
}

// helper function to safely dereference pointers
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (h *JobHandler) HandleGetJobsByLocation(res http.ResponseWriter, req *http.Request) {
	//Extract Query Parameter
	location := req.URL.Query().Get("location")
	//validate Query Params
	if location == "" {
		http.Error(res, "Location is required", http.StatusBadRequest)
		return
	}
	//Repository Call
	var pgLocation pgtype.Text
	err := pgLocation.Scan(location)
	if err != nil {
		http.Error(res, "Invalid Location format", http.StatusBadRequest)
		return
	}
	jobs, err := h.Repo.GetJobsByLocation(req.Context(), pgLocation)

	//validate jobs
	if err != nil {
		http.Error(res, "Failed to get jobs by location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Convert Response to JSON
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(jobs); err != nil {
		http.Error(res, "Failed to encode jobs", http.StatusInternalServerError)
		return
	}
	//Return
}

func (h *JobHandler) HandleGetAllLocation(res http.ResponseWriter, req *http.Request) {
	//Call Repo
	locations, err := h.Repo.GetAllLocation(req.Context())
	//Validate Location
	if err != nil {
		http.Error(res, "Failed to get locations", http.StatusInternalServerError)
		return
	}

	//Convert Response to JSON
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(locations); err != nil {
		http.Error(res, "Failed to encode locations", http.StatusInternalServerError)
		return
	}
	//Return
}

func (h *JobHandler) HandleSearchJobs(res http.ResponseWriter, req *http.Request) {
	//Extract Query Parameters : title , location
	title := req.URL.Query().Get("title")
	location := req.URL.Query().Get("location")
	//Repository call
	jobs, err := h.Repo.SearchJobs(req.Context(), title, location)
	//Validate
	if err != nil {
		http.Error(res, "Failed to search jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//Convert Response to JSON
	res.Header().Set("Content-Type", "application/json")
	//Return
	if err := json.NewEncoder(res).Encode(jobs); err != nil {
		http.Error(res, "Failed to encode jobs", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleGetJobWithPagination(res http.ResponseWriter, req *http.Request) {
	//Formula Offset = (PageNumber -1 ) X Limit
	//Extract page and Limit
	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("limit")

	//Parse String -> Integer
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	} //Default to page 1
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	} //Default 10 per page

	//Formula Apply
	offset := (page - 1) * limit

	//Repository Call
	jobs, total, err := h.Repo.GetJobsWithPagination(req.Context(), int32(offset), int32(limit))
	if err != nil {
		http.Error(res, "Failed to get jobs with pagination: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//Profession Response (Data + MetaData)
	response := map[string]interface{}{"data": jobs,
		"meta": map[string]interface{}{
			"current_page": page,
			"limit":        limit,
			"total_items":  total,
		},
	}
	//Return
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
