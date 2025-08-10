package db

// db.go (conexões)
import (
	"context"
	"database/sql"
	"time"

	// "github.com/go-redis/redis"
	"github.com/redis/go-redis/v9"

	_ "github.com/jackc/pgx/v5/stdlib"
	// "github.com/kovarike/protocol/src/go/worker"
)

var (
	Pg  *sql.DB
	Rdb *redis.Client
)

func InitDB() error {
	var err error
	Pg, err = sql.Open("pgx", "postgres://user:pass@localhost:5432/protocol?sslmode=disable")
	if err != nil {
		return err
	}
	Pg.SetMaxOpenConns(20)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // sem senha
		DB:       0,  // usar DB padrão
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	// modemPort := "/dev/ttyUSB0"

	// // Executa worker (bloqueante)
	// worker.WorkerLoop(ctx, modemPort)
	defer cancel()
	return Rdb.Ping(ctx).Err()
}
