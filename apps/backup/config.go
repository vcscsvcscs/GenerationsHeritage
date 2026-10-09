package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	defaultMemgraphURI = "bolt://memgraph:7687"
	defaultSchedule    = "0 3 * * *"
	defaultPrefix      = "memgraph/"
	defaultRegion      = "auto"
	defaultRetention   = 30
	defaultHealthFile  = "/tmp/last_success"
)

type config struct {
	schedule cron.Schedule

	memgraphURI  string
	memgraphUser string
	memgraphPass string

	s3Endpoint  string
	s3Bucket    string
	s3AccessKey string
	s3SecretKey string
	s3Region    string
	s3Prefix    string

	healthFile string

	ageRecipient string
	ageIdentity  string
	passphrase   string

	retentionCount int
	retentionDays  int
	s3PathStyle    bool
	onStart        bool
}

func envOr(getenv func(string) string, key, def string) string {
	if v := getenv(key); v != "" {
		return v
	}

	return def
}

func envInt(getenv func(string) string, key string, def int) (int, error) {
	v := getenv(key)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer, got %q", key, v)
	}

	return n, nil
}

func envBool(getenv func(string) string, key string) (bool, error) {
	v := getenv(key)
	if v == "" {
		return false, nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean, got %q", key, v)
	}

	return b, nil
}

func parseSchedule(getenv func(string) string) (cron.Schedule, error) {
	expr := envOr(getenv, "BACKUP_SCHEDULE", defaultSchedule)

	s, err := cron.ParseStandard(expr)
	if err != nil {
		return nil, fmt.Errorf("BACKUP_SCHEDULE %q: %w", expr, err)
	}

	return s, nil
}

func loadConfig(getenv func(string) string) (*config, error) {
	c := &config{
		memgraphURI:  envOr(getenv, "MEMGRAPH_URI", defaultMemgraphURI),
		memgraphUser: getenv("MEMGRAPH_USER"),
		memgraphPass: getenv("MEMGRAPH_PASSWORD"),
		s3Endpoint:   getenv("S3_ENDPOINT"),
		s3Bucket:     getenv("S3_BUCKET"),
		s3AccessKey:  getenv("S3_ACCESS_KEY_ID"),
		s3SecretKey:  getenv("S3_SECRET_ACCESS_KEY"),
		s3Region:     envOr(getenv, "S3_REGION", defaultRegion),
		s3Prefix:     normalizePrefix(envOr(getenv, "S3_PREFIX", defaultPrefix)),
		healthFile:   envOr(getenv, "BACKUP_HEALTH_FILE", defaultHealthFile),
		ageRecipient: getenv("BACKUP_AGE_RECIPIENT"),
		ageIdentity:  getenv("BACKUP_AGE_IDENTITY"),
		passphrase:   getenv("BACKUP_ENCRYPTION_PASSPHRASE"),
	}

	var errs []error

	for _, k := range []string{"S3_BUCKET", "S3_ACCESS_KEY_ID", "S3_SECRET_ACCESS_KEY"} {
		if getenv(k) == "" {
			errs = append(errs, fmt.Errorf("%s is required", k))
		}
	}

	if c.ageRecipient != "" && c.passphrase != "" {
		errs = append(errs, errors.New("BACKUP_AGE_RECIPIENT and BACKUP_ENCRYPTION_PASSPHRASE are mutually exclusive"))
	}

	var err error
	if c.schedule, err = parseSchedule(getenv); err != nil {
		errs = append(errs, err)
	}

	if c.s3PathStyle, err = envBool(getenv, "S3_PATH_STYLE"); err != nil {
		errs = append(errs, err)
	}

	if c.onStart, err = envBool(getenv, "BACKUP_ON_START"); err != nil {
		errs = append(errs, err)
	}

	if c.retentionCount, err = envInt(getenv, "BACKUP_RETENTION_COUNT", defaultRetention); err != nil {
		errs = append(errs, err)
	}

	if c.retentionDays, err = envInt(getenv, "BACKUP_RETENTION_DAYS", 0); err != nil {
		errs = append(errs, err)
	}

	return c, errors.Join(errs...)
}

func normalizePrefix(p string) string {
	p = strings.TrimLeft(p, "/")
	if p != "" && !strings.HasSuffix(p, "/") {
		p += "/"
	}

	return p
}

// healthMaxAge is twice the schedule interval plus an hour of slack.
func healthMaxAge(s cron.Schedule, now time.Time) time.Duration {
	next := s.Next(now)

	return 2*s.Next(next).Sub(next) + time.Hour
}
