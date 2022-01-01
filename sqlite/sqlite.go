package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	pa "github.com/Lambels/patrickarvatu.com"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	userCountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_users",
		Help: "total number of db users",
	})

	blogCountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_blogs",
		Help: "total number of blogs",
	})

	subBlogCountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_sub_blogs",
		Help: "total number of sub blogs",
	})

	commentCountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_comments",
		Help: "total number of comments",
	})

	subscriptionCountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_subscriptions",
		Help: "combined total of total number of blog subscriptions and sub blog subscriptions",
	})
)

// migration fs ./migration/*.sql
var migrationFS embed.FS

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()

	DSN string

	// sqlite package will use events, ie: blog:new, blog:sub_blog:new etc
	EventService pa.EventService

	Now func() time.Time
}

type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN:          dsn,
		Now:          time.Now,
		EventService: pa.NewNOPEventService(),
	}

	return db
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			return err
		}
	}

	if db.db, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}

	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	go db.monitor()

	return nil
}

func (db *DB) migrate() error {
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFS, "./migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("mgrateFile: err=%w name=%q", err, name)
		}
	}

	return nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) Close() error {
	db.cancel() // cancel ctx

	return db.db.Close()
}

func (db *DB) BeginTX(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// testing only
func (db *DB) MustBeginTX(ctx context.Context, opts *sql.TxOptions) *Tx {
	tx, err := db.BeginTX(ctx, opts)
	if err != nil {
		panic(err)
	}

	return tx
}

func (db *DB) monitor() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-db.ctx.Done(): // kill gorutine / stop monitoring when context is canceled
			return

		case <-ticker.C: // each tick represents one monitoring rutine
		}

		if err := db.update(db.ctx); err != nil {
			log.Printf("stats err: %s", err)
		}
	}
}

func (db *DB) update(ctx context.Context) error {
	tx, err := db.BeginTX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM users;`).Scan(&n); err != nil {
		return fmt.Errorf("user count: %w", err)
	}
	userCountGauge.Set(float64(n))

	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM blogs;`).Scan(&n); err != nil {
		return fmt.Errorf("blog count: %w", err)
	}
	blogCountGauge.Set(float64(n))

	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sub_blogs;`).Scan(&n); err != nil {
		return fmt.Errorf("sub_blog count: %w", err)
	}
	subBlogCountGauge.Set(float64(n))

	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM comments;`).Scan(&n); err != nil {
		return fmt.Errorf("comment count: %w", err)
	}
	userCountGauge.Set(float64(n))

	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM blog_subscriptions;`).Scan(&n); err != nil {
		return fmt.Errorf("dial count: %w", err)
	}
	blogSubscriptions := n

	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sub_blog_subscriptions;`).Scan(&n); err != nil {
		return fmt.Errorf("dial membership count: %w", err)
	}
	subscriptionCountGauge.Set(float64(blogSubscriptions + n))

	return nil
}

type NullTime time.Time

func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if value, ok := value.(string); ok {
		*(*time.Time)(n), _ = time.Parse(time.RFC3339, value)
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format(time.RFC3339), nil
}

func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}
