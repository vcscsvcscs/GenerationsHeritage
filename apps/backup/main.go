// Command backup periodically dumps Memgraph over Bolt and uploads the dump to an S3-compatible bucket.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	usage = `usage: backup [serve | run-once | restore [--force] <object-key|latest> | healthcheck]

With no command, backup serve runs backups on BACKUP_SCHEDULE.`
	attempts       = 3
	initialBackoff = 10 * time.Second
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, log, os.Args[1:]); err != nil {
		log.Error("backup failed", "error", err)

		return 1
	}

	return 0
}

func run(ctx context.Context, log *slog.Logger, args []string) error {
	cmd := "serve"
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	if cmd == "healthcheck" {
		return healthcheck(os.Getenv)
	}

	var (
		force  bool
		target string
	)

	switch cmd {
	case "serve", "run-once":
	case "restore":
		fs := flag.NewFlagSet("restore", flag.ContinueOnError)
		fs.BoolVar(&force, "force", false, "wipe a non-empty database before restoring")

		if err := fs.Parse(args); err != nil {
			return err //nolint:wrapcheck // flag error is self-explanatory
		}

		if fs.NArg() != 1 {
			return errors.New("restore takes exactly one argument: <object-key|latest>")
		}

		target = fs.Arg(0)
	default:
		return fmt.Errorf("unknown command %q\n%s", cmd, usage)
	}

	cfg, err := loadConfig(os.Getenv)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	b, closeDB, err := newBackupper(ctx, cfg, log)
	if err != nil {
		return err
	}
	defer closeDB()

	switch cmd {
	case "run-once":
		return b.runOnce(ctx)
	case "restore":
		return b.restore(ctx, target, force)
	}

	b.serve(ctx, cfg.schedule, cfg.onStart)

	return nil
}

func newBackupper(ctx context.Context, cfg *config, log *slog.Logger) (*backupper, func(), error) {
	rcp, err := cfg.recipient()
	if err != nil {
		return nil, nil, err
	}

	st, err := newS3Store(cfg)
	if err != nil {
		return nil, nil, err
	}

	db, err := newMemgraphDB(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	return &backupper{
		db:             db,
		store:          st,
		log:            log,
		prefix:         cfg.s3Prefix,
		recipient:      rcp,
		identity:       cfg.identity,
		retentionCount: cfg.retentionCount,
		retentionDays:  cfg.retentionDays,
		attempts:       attempts,
		backoff:        initialBackoff,
		healthFile:     cfg.healthFile,
		now:            time.Now,
	}, func() { _ = db.Close(context.WithoutCancel(ctx)) }, nil
}

// serve runs backups on the schedule until ctx is canceled; failed runs are logged and the loop continues.
func (b *backupper) serve(ctx context.Context, sched cron.Schedule, onStart bool) {
	b.touchHealth()
	b.log.Info("backup service started", "retention_count", b.retentionCount, "retention_days", b.retentionDays)

	if onStart {
		b.runLogged(ctx)
	}

	for {
		next := sched.Next(b.now())
		b.log.Info("next backup scheduled", "at", next.Format(time.RFC3339))

		if sleep(ctx, time.Until(next)) != nil {
			b.log.Info("shutting down")

			return
		}

		b.runLogged(ctx)
	}
}

func (b *backupper) runLogged(ctx context.Context) {
	if err := b.runOnce(ctx); err != nil {
		b.log.Error("scheduled backup failed", "error", err)
	}
}

func healthcheck(getenv func(string) string) error {
	sched, err := parseSchedule(getenv)
	if err != nil {
		return err
	}

	file := envOr(getenv, "BACKUP_HEALTH_FILE", defaultHealthFile)

	fi, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("health file: %w", err)
	}

	if age, maxAge := time.Since(fi.ModTime()), healthMaxAge(sched, time.Now()); age > maxAge {
		return fmt.Errorf("last successful backup %s ago, limit %s", age.Round(time.Second), maxAge)
	}

	return nil
}
