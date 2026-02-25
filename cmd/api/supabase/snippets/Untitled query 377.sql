-- +goose Up
CREATE TABLE jobs(
                     id UUID PRIMARY KEY,
                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                     updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
                     title TEXT NOT NULL,
                     job_role TEXT NOT NULL,
                     job_category TEXT,
                     company_name TEXT NOT NULL,
                     location TEXT NOT NULL,
                     salary TEXT DEFAULT 'Not Disclosed',
                     qualification TEXT NOT NULL,
                     experience TEXT NOT NULL,
                     last_date TEXT NOT NULL,
                     description TEXT NOT NULL,
                     apply_url TEXT NOT NULL,
                     is_active BOOLEAN DEFAULT TRUE
);
