package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// use connection pool
var DB *pgxpool.Pool

func Setup() (err error) {
	DB, err = pgxpool.Connect(context.Background(), fmt.Sprintf(`postgres://%s:%s@host.docker.internal:6432/%s`, viper.GetString("postgres_user"), viper.GetString("postgres_password"), viper.GetString("postgres_db")))
	if err != nil {
		zap.S().Fatalf("failed to connect to database %+v", err)
		return
	}
	if err = DB.Ping(context.Background()); err != nil {
		zap.S().Fatalf("failed to ping database %+v", err.Error())
	}
	return
}

func Store(commands int, cleaned int, seconds float64) (ID int, timestamp time.Time, err error) {
	row := DB.QueryRow(context.Background(), `INSERT INTO executions ("timestamp",commands,"result",duration) VALUES (NOW(),$1,$2,$3) RETURNING id, "timestamp"`, commands, cleaned, seconds)
	err = row.Scan(&ID, &timestamp)
	if err != nil {
		zap.S().Errorf("failed to insert into executions - %w", err)
		return
	}
	return
}

func List(offset int, limit int) (rows []Row, err error) {
	var rowsdb pgx.Rows
	rowsdb, err = DB.Query(context.Background(), `SELECT id,"timestamp",commands,"result",duration FROM executions ORDER BY id ASC OFFSET $1 LIMIT $2`, offset, limit)
	if err != nil {
		zap.S().Errorf("failed to select rows: %w", err)
		return
	}

	for rowsdb.Next() {
		var ID int
		var timestamp time.Time
		var commands int
		var result int
		var duration float64
		err = rowsdb.Scan(&ID, &timestamp, &commands, &result, &duration)
		if err != nil {
			zap.S().Errorf("failed to scan row: %w", err)
			return
		}
		rows = append(rows, Row{ID: ID, Timestamp: timestamp, Commands: commands, Result: result, Duration: duration})
	}
	return
}

type Row struct {
	ID        int
	Timestamp time.Time
	Commands  int
	Result    int
	Duration  float64
}
