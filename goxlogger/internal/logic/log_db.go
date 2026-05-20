package logic

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// LogPostgresDb contains all the logic to save, retrieve
// and manage logs in a PostgreSQL database
type LogPostgresDb struct {
	ctx  context.Context
	conn *pgx.Conn
}

type LogDbInterface interface {
	SaveLogs(logLines []LogLine) error
	GetLogs(options LogSelectOptions) ([]LogLine, error)
	DeleteLogs(options LogSelectOptions) error
	Close() error
	Initialize() error
	BroadcastLog(logLine LogLine)
	StartBroadcaster() <-chan error
}

type LogSelectOptions struct {
	StatusCode int
	Method     string
}

func (l *LogPostgresDb) Initialize() error {
	// Create the logs table if it doesn't exist
	_, err := l.conn.Exec(l.ctx, `
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMPTZ NOT NULL,
			log JSONB NOT NULL
		)
	`)

	// Create indexes on the timestamp and
	// log fields for faster querying
	_, err = l.conn.Exec(l.ctx, `
		CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs (timestamp);
		CREATE INDEX IF NOT EXISTS idx_logs_log ON logs USING gin (log);
	`)
	if err != nil {
		return err
	}

	return err
}

func (l *LogPostgresDb) SaveLogs(logLines []LogLine) error {
	return nil
}

func (l *LogPostgresDb) GetLogs(options LogSelectOptions) ([]LogLine, error) {
	// whereClause := []string{}

	// if options.StatusCode != 0 {
	// 	whereClause = append(whereClause, `log->>'status_code' = $1`)
	// }

	// if options.Method != "" {
	// 	whereClause = append(whereClause, `log->>'method' = $2`)
	// }

	// whereSql := strings.Join(whereClause, " AND ")
	// sql := fmt.Sprintf(`
	// 	SELECT log
	// 	FROM logs
	// 	WHERE %s
	// `,
	// 	whereSql,
	// )

	// cmd, err := l.conn.Exec(l.ctx, sql)

	// for cmd.Next() {
	// 	var logLine LogLine
	// 	err = cmd.Scan(&logLine)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return nil, nil
}

func (l *LogPostgresDb) DeleteLogs(options LogSelectOptions) error {
	return nil
}

func (l *LogPostgresDb) Close() error {
	return l.conn.Close(context.Background())
}

func (l *LogPostgresDb) BroadcastLog(logLine LogLine) {
	// No broadcasting logic for PostgreSQL
}

func (l *LogPostgresDb) StartBroadcaster() <-chan error {
	ch := make(chan error, 1)
	close(ch) // No broadcaster, so we close the channel immediately
	return ch
}

func NewLogPostgresDb(ctx context.Context) LogDbInterface {
	conn, err := pgx.Connect(ctx, "postgres://user:password@localhost:5432/logs_db")

	if err != nil {
		panic(err)
	}

	return &LogPostgresDb{
		ctx:  ctx,
		conn: conn,
	}
}
