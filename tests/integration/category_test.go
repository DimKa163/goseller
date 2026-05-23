package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DimKa163/goseller/internal/category/domain"
	"github.com/DimKa163/goseller/internal/category/infrastructure"
	categoryinterfaces "github.com/DimKa163/goseller/internal/category/interfaces"
	"github.com/DimKa163/goseller/internal/category/usecase"
	"github.com/DimKa163/goseller/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestCategoryAPI_CreateShouldBeSuccess(t *testing.T) {
	ctx := context.Background()
	container, srv, db, err := run_categories(ctx, func(pool *pgxpool.Pool) error {
		// Здесь можно выполнить начальную настройку базы данных, если это необходимо
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	}()
	defer srv.Close()
	defer db.Close()

	req := &usecase.CategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	buffer := bytes.NewBuffer(data)
	response, err := makeRequest(ctx, "POST", fmt.Sprintf("%s/category", srv.URL), buffer)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	location := response.Header.Get("Location")
	assert.NotEmpty(t, location)
	var s struct {
		ID domain.CategoryID `json:"id"`
	}
	if err := readResponseBody(response, &s); err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}
	assert.NotEmpty(t, s.ID)

	response, err = makeRequest(ctx, "GET", fmt.Sprintf("%s%s", srv.URL, location), nil)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Greater(t, response.ContentLength, int64(0))
	var category domain.Category
	if err := readResponseBody(response, &category); err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}

	assert.Equal(t, req.Name, category.Name)
	assert.Equal(t, req.Description, category.Description)

	req = &usecase.CategoryRequest{
		Name:        "Updated Category",
		Description: "Updated Description",
	}
	data, err = json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	buffer = bytes.NewBuffer(data)
	response, err = makeRequest(ctx, "PUT", fmt.Sprintf("%s/category/%s", srv.URL, category.ID.String()), buffer)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Greater(t, response.ContentLength, int64(0))
	var newCategory domain.Category
	if err := readResponseBody(response, &newCategory); err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}

	assert.Equal(t, req.Name, newCategory.Name)
	assert.Equal(t, req.Description, newCategory.Description)

	response, err = makeRequest(ctx, "DELETE", fmt.Sprintf("%s/category/%s", srv.URL, category.ID.String()), nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	assert.Equal(t, http.StatusNoContent, response.StatusCode)

	response, err = makeRequest(ctx, "GET", fmt.Sprintf("%s/category/%s", srv.URL, category.ID.String()), nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	var inactive bool
	if err = db.QueryRow(ctx, `SELECT inactive FROM categories WHERE id = $1`, category.ID).Scan(&inactive); err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	assert.True(t, inactive)
}

func run_categories(ctx context.Context, beFn func(pool *pgxpool.Pool) error) (testcontainers.Container, *httptest.Server, *pgxpool.Pool, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForListeningPort("5432").WithStartupTimeout(160 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, nil, err
	}
	databaseURL := fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", host, port.Port())

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, nil, nil, err
	}
	if err := infrastructure.Migrate(db, "../../internal/category/migrations"); err != nil {
		return nil, nil, nil, err
	}
	if err := beFn(db); err != nil {
		return nil, nil, nil, err
	}
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, nil, nil, err
	}
	router := gin.New()
	router.Use(middleware.LoggingMiddleware(log))
	controller := categoryinterfaces.NewCategoryController(log, usecase.NewCategoryAppService(log, infrastructure.NewCategoryRepository(db)))
	controller.Map(router)
	return container, httptest.NewServer(router), db, nil
}
