package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const usage = `
Jesterx CLI

Uso: jx <comando> [flags]

Comandos:
  migrate           Aplica todas as migrations pendentes
  migrate:status    Lista todas as migrations e seus status
  migrate:create <nome>   Cria uma nova migration vazia
  dev               Inicia o servidor de desenvolvimento (go run cmd/api/main.go)
  db:seed           Executa seeds de desenvolvimento
  routes:list       Lista todas as rotas registradas
  gen:handler <nome>  Gera um novo handler boilerplate

Exemplos:
  jx migrate
  jx migrate:create add_user_settings
  jx dev
  jx gen:handler user
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := os.Args[1]

	switch cmd {
	case "migrate":
		runMigrate(false)
	case "migrate:status":
		runMigrateStatus()
	case "migrate:create":
		if len(os.Args) < 3 {
			fmt.Println("Uso: jx migrate:create <nome>")
			os.Exit(1)
		}
		runMigrateCreate(os.Args[2])
	case "dev":
		runDev()
	case "db:seed":
		runSeed()
	case "routes:list":
		runRoutesList()
	case "gen:handler":
		if len(os.Args) < 3 {
			fmt.Println("Uso: jx gen:handler <nome>")
			os.Exit(1)
		}
		runGenHandler(os.Args[2])
	default:
		fmt.Printf("Comando desconhecido: %s\n", cmd)
		fmt.Print(usage)
		os.Exit(1)
	}
}

func loadDB() *pgxpool.Pool {
	_ = godotenv.Load(".env")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	pass := getEnv("POSTGRES_PASSWORD", "postgres")
	db := getEnv("POSTGRES_DB", "jesterx")
	ssl := getEnv("POSTGRES_SSL", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, db, ssl)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Printf("Erro ao conectar ao banco: %v\n", err)
		os.Exit(1)
	}
	return pool
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrate(dryRun bool) {
	db := loadDB()
	defer db.Close()

	ctx := context.Background()
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Printf("Erro ao criar tabela schema_migrations: %v\n", err)
		os.Exit(1)
	}

	rows, _ := db.Query(ctx, `SELECT version FROM schema_migrations`)
	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		rows.Scan(&v)
		applied[v] = true
	}
	rows.Close()

	entries, err := os.ReadDir("migrations")
	if err != nil {
		fmt.Printf("Erro ao ler pasta migrations: %v\n", err)
		os.Exit(1)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	appliedCount := 0
	for _, fname := range files {
		version := strings.TrimSuffix(fname, ".up.sql")
		if applied[version] {
			fmt.Printf("  ✓ %s (já aplicado)\n", fname)
			continue
		}

		content, err := os.ReadFile(filepath.Join("migrations", fname))
		if err != nil {
			fmt.Printf("  ✗ %s: erro ao ler: %v\n", fname, err)
			continue
		}

		if dryRun {
			fmt.Printf("  → %s (pendente)\n", fname)
			continue
		}

		fmt.Printf("  → Aplicando %s...\n", fname)
		if _, err := db.Exec(ctx, string(content)); err != nil {
			fmt.Printf("  ✗ Erro ao aplicar %s: %v\n", fname, err)
			os.Exit(1)
		}

		db.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, version)
		fmt.Printf("  ✓ %s aplicado\n", fname)
		appliedCount++
	}

	if appliedCount == 0 && !dryRun {
		fmt.Println("  ✓ Nenhuma migration pendente")
	} else if !dryRun {
		fmt.Printf("\n✓ %d migration(s) aplicada(s)\n", appliedCount)
	}
}

func runMigrateStatus() {
	db := loadDB()
	defer db.Close()
	ctx := context.Background()

	_, _ = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`)

	rows, _ := db.Query(ctx, `SELECT version, applied_at FROM schema_migrations ORDER BY applied_at`)
	applied := make(map[string]string)
	for rows.Next() {
		var v, at string
		rows.Scan(&v, &at)
		applied[v] = at
	}
	rows.Close()

	entries, _ := os.ReadDir("migrations")
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	fmt.Println("\nStatus das Migrations:")
	fmt.Println(strings.Repeat("─", 70))
	for _, f := range files {
		version := strings.TrimSuffix(f, ".up.sql")
		if at, ok := applied[version]; ok {
			fmt.Printf("  ✓ %-45s aplicado em %s\n", f, at[:10])
		} else {
			fmt.Printf("  ○ %-45s pendente\n", f)
		}
	}
	fmt.Println(strings.Repeat("─", 70))
}

func runMigrateCreate(name string) {
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)

	entries, _ := os.ReadDir("migrations")
	maxN := 0
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			var n int
			fmt.Sscanf(e.Name(), "%04d_", &n)
			if n > maxN {
				maxN = n
			}
		}
	}

	fname := fmt.Sprintf("%04d_%s.up.sql", maxN+1, name)
	fpath := filepath.Join("migrations", fname)

	content := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Write your SQL here\n",
		name, time.Now().Format("2006-01-02"))

	if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
		fmt.Printf("Erro ao criar migration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Migration criada: %s\n", fpath)
}

func runDev() {
	fmt.Println("Iniciando servidor de desenvolvimento...")
	fmt.Println("Execute: go run cmd/api/main.go")
	fmt.Println("\nOu use o air para hot-reload:")
	fmt.Println("  go install github.com/air-verse/air@latest")
	fmt.Println("  air")
}

func runSeed() {
	fmt.Println("Seeds de desenvolvimento:")
	fmt.Println("  As migrations de seed (0002_Seed_default_data.up.sql, 0007_Seed_themes.up.sql)")
	fmt.Println("  são aplicadas automaticamente com 'jx migrate'")
	fmt.Println("\nPara seeds manuais, crie arquivos em migrations/ e execute 'jx migrate'")
}

func runRoutesList() {
	routes := []struct{ method, path, desc string }{
		{"POST", "/api/v1/auth/register", "Registrar usuário"},
		{"GET", "/api/v1/auth/verify/{token}", "Verificar email"},
		{"POST", "/api/v1/auth/login", "Login"},
		{"GET", "/api/v1/auth/refresh", "Refresh token"},
		{"GET", "/api/v1/auth/me", "Perfil do usuário"},
		{"PATCH", "/api/v1/auth/me", "Atualizar perfil"},
		{"GET", "/api/v1/auth/logout", "Logout"},
		{"GET", "/api/v1/websites", "Listar sites"},
		{"POST", "/api/v1/websites", "Criar site"},
		{"DELETE", "/api/v1/sites/{siteID}", "Deletar site"},
		{"GET", "/api/v1/sites/{siteID}/routes", "Listar rotas"},
		{"POST", "/api/v1/sites/{siteID}/routes", "Definir rotas"},
		{"GET", "/api/v1/sites/{siteID}/versions", "Listar versões"},
		{"POST", "/api/v1/sites/{siteID}/versions", "Criar versão"},
		{"POST", "/api/v1/sites/{siteID}/publish/{v}", "Publicar versão"},
		{"POST", "/api/v1/sites/{siteID}/products", "Criar produto"},
		{"GET", "/api/v1/sites/{siteID}/products", "Listar produtos"},
		{"PATCH", "/api/v1/sites/{siteID}/products/{id}", "Atualizar produto"},
		{"DELETE", "/api/v1/sites/{siteID}/products/{id}", "Deletar produto"},
		{"GET", "/api/v1/sites/{siteID}/orders", "Listar pedidos (admin)"},
		{"POST", "/api/store/{siteID}/orders", "Criar pedido (público)"},
		{"GET", "/api/store/{siteID}/products", "Listar produtos (público)"},
		{"GET", "/api/store/{siteID}/products/{id}", "Detalhe produto (público)"},
		{"POST", "/api/v1/upload", "Upload de arquivo"},
		{"GET", "/api/v1/themes", "Listar temas"},
		{"GET", "/api/v1/plans", "Listar planos"},
		{"POST", "/api/v1/payments/checkout", "Criar checkout"},
		{"GET", "/api/v1/payments/confirm", "Confirmar pagamento"},
		{"POST", "/api/v1/payments/webhook", "Stripe webhook"},
		{"GET", "/p/{siteID}/{path...}", "Página pública"},
	}

	fmt.Println("\nRotas Registradas:")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("  %-8s %-45s %s\n", "MÉTODO", "ROTA", "DESCRIÇÃO")
	fmt.Println(strings.Repeat("─", 80))
	for _, r := range routes {
		fmt.Printf("  %-8s %-45s %s\n", r.method, r.path, r.desc)
	}
	fmt.Println(strings.Repeat("─", 80))
}

func runGenHandler(name string) {
	name = strings.ToLower(name)
	capitalized := strings.ToUpper(name[:1]) + name[1:]

	tmpl := fmt.Sprintf(`package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
)

type %sHandler struct {
	// TODO: add service dependency
}

func New%sHandler() *%sHandler {
	return &%sHandler{}
}

// GET /api/v1/%ss
func (h *%sHandler) List(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: implement
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": []any{}})
}

// POST /api/v1/%ss
func (h *%sHandler) Create(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// TODO: implement
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "%s created"})
}

// DELETE /api/v1/%ss/{id}
func (h *%sHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := strings.TrimSpace(r.PathValue("id"))
	_ = id
	// TODO: implement

	w.WriteHeader(http.StatusNoContent)
}
`,
		capitalized, capitalized, capitalized, capitalized,
		name, capitalized,
		name, capitalized,
		name,
		name, capitalized,
	)

	fname := fmt.Sprintf("internal/http/handlers/%s_handler.go", name)
	if _, err := os.Stat(fname); err == nil {
		fmt.Printf("Arquivo já existe: %s\n", fname)
		os.Exit(1)
	}

	if err := os.WriteFile(fname, []byte(tmpl), 0644); err != nil {
		fmt.Printf("Erro ao criar handler: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Handler criado: %s\n", fname)
	fmt.Printf("  Registre as rotas em internal/http/router.go\n")
}
