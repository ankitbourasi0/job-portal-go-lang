-- name: CreateJob :one
INSERT INTO jobs(id, created_at, updated_at, title, job_role, job_category ,company_name, location, salary, qualification, experience, last_date, description, apply_url)
VALUES($1, $2 ,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    RETURNING *;

-- name: GetAllJobs :many
SELECT * FROM jobs
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: GetJobById :one
SELECT * FROM jobs WHERE id = $1 LIMIT 1;

-- name: UpdateJobById :one
UPDATE jobs
SET
    title = $2,
    job_role = $3,
    job_category = $4,
    company_name = $5,
    location = $6,
    salary = $7,
    qualification = $8,
    experience = $9,
    last_date = $10,
    description = $11,
    apply_url = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PartialUpdateJob :one
UPDATE jobs
SET
    title = COALESCE(sqlc.narg('title'), title), --Coalesce fn skip the column if value is null comed, otherwise update , helpful in patch
    job_role = COALESCE(sqlc.narg('job_role'), job_role),
    job_category = COALESCE(sqlc.narg('job_category'), job_category),
    company_name = COALESCE(sqlc.narg('company_name'), company_name),
    location = COALESCE(sqlc.narg('location'), location),
    salary = COALESCE(sqlc.narg('salary'), salary),
    qualification = COALESCE(sqlc.narg('qualification'), qualification),
    experience = COALESCE(sqlc.narg('experience'), experience),
    last_date = COALESCE(sqlc.narg('last_date'), last_date),
    description = COALESCE(sqlc.narg('description'), description),
    apply_url = COALESCE(sqlc.narg('apply_url'), apply_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetJobsByLocation :many
SELECT * FROM jobs
WHERE location ILIKE '%' || $1 || '%' --this '%' and ILIKE actually means case senstive and search data from multiwords
AND is_active = true;

-- name: GetAllLocation :many
SELECT DISTINCT location --Distinct means get unique records
FROM jobs
WHERE is_active =true
ORDER BY  location ASC;


-- name: SearchJobs :many
SELECT * FROM jobs
WHERE (title ILIKE '%' || sqlc.arg(title)::text || '%')
    AND (location ILIKE '%' || sqlc.arg(location)::text || '%')
    AND is_active =true
ORDER BY created_at DESC ;

-- name: GetJobsWithPagination :many
SELECT * FROM jobs
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2; --Limit tells how many records need to get e.g. 10

-- name: GetTotalCount :one
SELECT COUNT(*) FROM jobs WHERE is_active = true; --It tells how many records need to skip
