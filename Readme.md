
## **Backend:**

1.  **Programming** : Go(1.25.5)
2.  **Database ORM**: SQLC
3.  **Database :** Supabase CLI(Postgres)
4.  **Production Database:** Supabase Cloud
5.  **Routing:** Chi Router(v5)
6.  **Database Migration Tool - **Goose** ya **Atlas**.**

Project Structure:

```Bash
job-portal/
├── cmd/
│   └── api/
│       └── main.go          # Entry point
├── internal/
│   ├── database/            # SQLC generated code
│   ├── repository/          # Data logic
│   ├── service/             # Business logic
│   └── handler/             # HTTP handlers
|   └── middleware/          # CORS, Logging, Auth middleware
|   └── auth/ 		     # Supabase JWT validation logic
├── sql/
│   ├── schema/              # Database migrations
│   └── queries/             # Raw SQL queries for SQLC
├── sqlc.yaml                # SQLC configuration
├── go.mod
└── .env                     # Environment variables
├── Dockerfile               # For Cloud Run/Racknerd deployment
```

Dependencies for Backend:

1.  **Chi Router** - go get github.com/go-chi/chi/v5
2.  **Driver Postgres**\- go get github.com/lib/pq
3.  **UUID** -  go get github.com/google/uuid
4.  **CORS Middleware -** go get github.com/go-chi/cors
5.  **Envrionment Variable(.env)** : go get github.com/joho/godotenv
6.  **SQLC** - go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

* * *