package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/application/command"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/application/query"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/company"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/infrastructure/postgres"
	redisInfra "github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/infrastructure/redis"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/interfaces/http/handler"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/company?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	if err := runMigrations(pool); err != nil {
		log.Printf("Warning: migration error: %v", err)
	}

	redisClient, err := redisInfra.NewRedisClient(ctx)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}
	_ = redisClient

	companyRepo := postgres.NewCompanyPostgresRepository(pool)
	branchRepo := postgres.NewBranchPostgresRepository(pool)

	publisher := &mockPublisher{}

	createCompanyCmd := command.NewCreateCompanyHandler(companyRepo, publisher)
	updateCompanyCmd := command.NewUpdateCompanyHandler(companyRepo)
	addBranchCmd := command.NewAddBranchHandler(companyRepo, branchRepo, publisher)
	renameBranchCmd := command.NewRenameBranchHandler(branchRepo)

	getCompanyQuery := query.NewGetCompanyHandler(companyRepo)
	listCompaniesQuery := query.NewListCompaniesHandler(companyRepo)
	listBranchesQuery := query.NewListBranchesHandler(branchRepo)
	getBranchQuery := query.NewGetBranchHandler(branchRepo)

	companyHandler := handler.NewCompanyHandler(
		createCompanyCmd, updateCompanyCmd,
		addBranchCmd, renameBranchCmd,
		getCompanyQuery, listCompaniesQuery,
		listBranchesQuery, getBranchQuery,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /v1/companies", companyHandler.Create)
	mux.HandleFunc("GET /v1/companies", companyHandler.List)
	mux.HandleFunc("GET /v1/companies/{id}", companyHandler.GetByID)
	mux.HandleFunc("PATCH /v1/companies/{id}", companyHandler.Update)
	mux.HandleFunc("POST /v1/companies/{id}/branches", companyHandler.AddBranch)
	mux.HandleFunc("GET /v1/companies/{id}/branches", companyHandler.ListBranches)
	mux.HandleFunc("GET /v1/companies/{id}/branches/{branchId}", companyHandler.GetBranch)
	mux.HandleFunc("PATCH /v1/companies/{id}/branches/{branchId}", companyHandler.RenameBranch)

	loggingMux := loggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggingMux,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		log.Println("Shutting down gracefully...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Company Service starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations(pool *pgxpool.Pool) error {
	migration := `
	CREATE TABLE IF NOT EXISTS companies (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		cnpj VARCHAR(14) NOT NULL UNIQUE,
		street VARCHAR(255) NOT NULL,
		city VARCHAR(100) NOT NULL,
		state VARCHAR(100) NOT NULL,
		postal_code VARCHAR(20) NOT NULL,
		country VARCHAR(100) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS company_branches (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		company_id VARCHAR(64) NOT NULL REFERENCES companies(id),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_companies_cnpj ON companies(cnpj);
	CREATE INDEX IF NOT EXISTS idx_branches_company_id ON company_branches(company_id);
	`
	_, err := pool.Exec(context.Background(), migration)
	return err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed: %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

type mockPublisher struct{}

func (m *mockPublisher) Publish(ctx context.Context, event company.DomainEvent) error {
	log.Printf("[PubSub] Event: %s Type: %s", event.EventID, event.EventType)
	return nil
}
