-- name: CreateJob :one
INSERT INTO jobs(id, created_at, updated_at, title, job_role, job_category ,company_name, location, salary, qualification, experience, last_date, description, apply_url)
VALUES($1, $2 ,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    RETURNING *;

-- name: GetAllJobs :many
SELECT * FROM jobs
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: GeyJobById :one
SELECT * FROM jobs WHERE id = $1 LIMIT 1;

