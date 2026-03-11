package repository

import (
	"context"

	"github.com/ankitbourasi0/job-portal/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobRepository struct {
	Queries *database.Queries
}

// Constructor : Go does not have constructors i.e we use New.. Functions as Cts...
// This fn accept dbQueries, it doest not know who and how the database source,
// but still it will work bcuz of dependency injection we put those queries JobRepository Queries
func NewJobRepository(dbQueries *database.Queries) *JobRepository {
	return &JobRepository{Queries: dbQueries}
}

// Interface for Service Layer, irrespective of ORM or Framework,
// they will call always this to create a job: Dependency Injection
func (r *JobRepository) CreateNewJob(ctx context.Context, arg database.CreateJobParams) (database.Job, error) {
	return r.Queries.CreateJob(ctx, arg)
}

func (r *JobRepository) GetAllJobs(ctx context.Context) ([]database.Job, error) {
	return r.Queries.GetAllJobs(ctx)
}

// GetJobByID fetches a single job by its UUID
func (r *JobRepository) GetJobById(ctx context.Context, id pgtype.UUID) (database.Job, error) {
	return r.Queries.GetJobById(ctx, id)
}

// here we use COALESCE which does not update the field in db if we pass null, maintain old values, work in 1 db call
func (r *JobRepository) PartialUpdateJob(ctx context.Context, arg database.PartialUpdateJobParams) (database.Job, error) {
	return r.Queries.PartialUpdateJob(ctx, arg)
}

func (r *JobRepository) GetJobsByLocation(ctx context.Context, location pgtype.Text) ([]database.Job, error) {
	return r.Queries.GetJobsByLocation(ctx, location)
}

// we used Distinct keyword in queries which get all unique values of location,
func (r *JobRepository) GetAllLocation(ctx context.Context) ([]string, error) {
	return r.Queries.GetAllLocation(ctx)
}

func (r *JobRepository) SearchJobs(ctx context.Context, title, location string) ([]database.Job, error) {
	return r.Queries.SearchJobs(ctx, database.SearchJobsParams{
		Title:    title,
		Location: location,
	})
}

func (r *JobRepository) GetJobsWithPagination(ctx context.Context, limit, offset int32) ([]database.Job, int64, error) {
	//1. Fetch Jobs
	jobs, err := r.Queries.GetJobsWithPagination(ctx, database.GetJobsWithPaginationParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, err
	}

	//Get Total Count(for Frontend Pagination)
	total, err := r.Queries.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}
