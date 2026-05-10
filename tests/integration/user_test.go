package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/DimKa163/goseller/internal/middleware"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/shared/resterror"
	"github.com/DimKa163/goseller/internal/user/domain"
	"github.com/DimKa163/goseller/internal/user/infrastructure"
	userinterfaces "github.com/DimKa163/goseller/internal/user/interface"
	"github.com/DimKa163/goseller/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestUserApi_CreateShouldBeSuccess(t *testing.T) {
	ctx := context.Background()
	addr := ":8081"
	container, db, err := run(ctx, addr, func(pool *pgxpool.Pool) error {
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

	defer db.Close()
	rep := infrastructure.NewUserRepository(db)
	req := &usecase.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "1234567890",
	}

	client := http.Client{Timeout: 30 * time.Second}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	buffer := bytes.NewBuffer(data)
	requestMessage, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost%s/user", addr), buffer)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	response, err := client.Do(requestMessage)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	location := response.Header.Get("Location")
	assert.NotEmpty(t, location)

	defer response.Body.Close()
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}
	var s struct {
		ID int64 `json:"id"`
	}
	if err = json.Unmarshal(respData, &s); err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}
	user, err := rep.GetByID(ctx, domain.UserID(s.ID))

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserAPI_CreateShouldBeFailed(t *testing.T) {
	ctx := context.Background()
	addr := ":8082"
	container, db, err := run(ctx, addr, func(pool *pgxpool.Pool) error {
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

	defer db.Close()
	rep := infrastructure.NewUserRepository(db)
	cases := []struct {
		Name               string
		Request            *usecase.CreateUserRequest
		ExprectedError     resterror.ErrorResponse
		ExpectedStatusCode int
		Before             func(db domain.UserRepository) (domain.UserID, error)
		After              func(db *pgxpool.Pool, id domain.UserID) error
	}{
		{
			Name: "Invalid phone",
			Request: &usecase.CreateUserRequest{
				Name:  "Test",
				Phone: domain.Phone("invalid phone"),
				Email: domain.Email("test@test.com"),
			},
			ExprectedError: resterror.ErrorResponse{
				Error: &resterror.Error{
					Code:    int(shared.ErrorCodeBadInputData),
					Message: "invalid request",
					Details: []*resterror.ErrorDetail{
						{
							Field:   "Phone",
							Message: "invalid format",
							Value:   "invalid phone",
						},
					},
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			Before:             func(db domain.UserRepository) (domain.UserID, error) { return domain.UserID(-1), nil },
			After: func(db *pgxpool.Pool, id domain.UserID) error {
				return nil
			},
		},
		{
			Name: "Invalid email",
			Request: &usecase.CreateUserRequest{
				Name:  "Test",
				Phone: domain.Phone("79270192274"),
				Email: domain.Email("invalid-email"),
			},
			ExprectedError: resterror.ErrorResponse{
				Error: &resterror.Error{
					Code:    int(shared.ErrorCodeBadInputData),
					Message: "invalid request",
					Details: []*resterror.ErrorDetail{
						{
							Field:   "Email",
							Message: "invalid format",
							Value:   "invalid-email",
						},
					},
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			Before:             func(db domain.UserRepository) (domain.UserID, error) { return domain.UserID(-1), nil },
			After: func(db *pgxpool.Pool, id domain.UserID) error {
				return nil
			},
		},
		{
			Name: "Invalid email and phone",
			Request: &usecase.CreateUserRequest{
				Name:  "Test",
				Phone: domain.Phone("invalid phone"),
				Email: domain.Email("invalid-email"),
			},
			ExprectedError: resterror.ErrorResponse{
				Error: &resterror.Error{
					Code:    int(shared.ErrorCodeBadInputData),
					Message: "invalid request",
					Details: []*resterror.ErrorDetail{
						{
							Field:   "Email",
							Message: "invalid format",
							Value:   "invalid-email",
						},
						{
							Field:   "Phone",
							Message: "invalid format",
							Value:   "invalid phone",
						},
					},
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			Before:             func(db domain.UserRepository) (domain.UserID, error) { return domain.UserID(-1), nil },
			After: func(db *pgxpool.Pool, id domain.UserID) error {
				return nil
			},
		},
		{
			Name: "Email already exists",
			Request: &usecase.CreateUserRequest{
				Name:  "Test",
				Phone: domain.Phone("79270192274"),
				Email: domain.Email("test@test.com"),
			},
			ExprectedError: resterror.ErrorResponse{
				Error: &resterror.Error{
					Code:    int(shared.ErrorCodeResourceAlreadyExists),
					Message: "resource already exists",
					Details: []*resterror.ErrorDetail{
						{
							Message: "email already exist",
						},
					},
				},
			},
			ExpectedStatusCode: http.StatusConflict,
			Before: func(db domain.UserRepository) (domain.UserID, error) {
				return db.Create(context.Background(), domain.CreateNewUser("Test email already exist", domain.Email("test@test.com"), domain.Phone("79270192275")))
			},
			After: func(db *pgxpool.Pool, id domain.UserID) error {
				if _, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "Phone already exists",
			Request: &usecase.CreateUserRequest{
				Name:  "Test",
				Phone: domain.Phone("79270192274"),
				Email: domain.Email("test@test.com"),
			},
			ExprectedError: resterror.ErrorResponse{
				Error: &resterror.Error{
					Code:    int(shared.ErrorCodeResourceAlreadyExists),
					Message: "resource already exists",
					Details: []*resterror.ErrorDetail{
						{
							Message: "phone already exist",
						},
					},
				},
			},
			ExpectedStatusCode: http.StatusConflict,
			Before: func(db domain.UserRepository) (domain.UserID, error) {
				return db.Create(context.Background(), domain.CreateNewUser("Test email already exist", domain.Email("test1@test.com"), domain.Phone("79270192274")))
			},
			After: func(db *pgxpool.Pool, id domain.UserID) error {
				if _, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id); err != nil {
					return err
				}
				return nil
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			id, err := c.Before(rep)
			if err != nil {
				t.Fatalf("Failed to prepare test: %v", err)
			}
			client := http.Client{Timeout: 30 * time.Second}

			data, err := json.Marshal(c.Request)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			buffer := bytes.NewBuffer(data)
			requestMessage, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost%s/user", addr), buffer)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			response, err := client.Do(requestMessage)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			assert.Equal(t, c.ExpectedStatusCode, response.StatusCode)

			defer response.Body.Close()
			respData, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Failed to parse body: %v", err)
			}
			var errResponse resterror.ErrorResponse
			if err := json.Unmarshal(respData, &errResponse); err != nil {
				t.Fatalf("Failed to unmarshal error response body: %v", err)
			}

			assert.Equal(t, c.ExprectedError, errResponse, "error should be equal")
			if err := c.After(db, id); err != nil {
				t.Fatalf("failed cleanup %v", err)
			}
		})
	}
}

func run(ctx context.Context, addr string, beFn func(pool *pgxpool.Pool) error) (testcontainers.Container, *pgxpool.Pool, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
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
		return nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}
	databaseURL := fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", host, port.Port())

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, nil, err
	}
	if err := infrastructure.Migrate(db, "../../internal/user/migrations"); err != nil {
		return nil, nil, err
	}
	if err := beFn(db); err != nil {
		return nil, nil, err
	}
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, nil, err
	}
	router := gin.New()
	router.Use(middleware.LoggingMiddleware(log))
	controller := userinterfaces.NewUserController(log, usecase.NewUser(infrastructure.NewUserRepository(db), log))
	controller.Map(router)
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatal("Failed to run server", zap.Error(err))
		}
	}()
	return container, db, nil
}
