package main

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"filippo.io/age"
)

const (
	namePrefix = "gheritage-"
	nameSuffix = ".cypher.gz"

	healthFileMode = 0o600
)

type object struct {
	modified time.Time
	key      string
	size     int64
}

type database interface {
	Dump(ctx context.Context, emit func(stmt string) error) error
	Exec(ctx context.Context, stmt string) error
	IsEmpty(ctx context.Context) (bool, error)
	Wipe(ctx context.Context) error
}

type store interface {
	Put(ctx context.Context, key string, r io.Reader) error
	Stat(ctx context.Context, key string) (int64, error)
	List(ctx context.Context, prefix string) ([]object, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, keys []string) error
}

type backupper struct {
	db             database
	store          store
	recipient      age.Recipient
	log            *slog.Logger
	identity       func() (age.Identity, error)
	now            func() time.Time
	prefix         string
	healthFile     string
	retentionCount int
	retentionDays  int
	attempts       int
	backoff        time.Duration
}

type countingReader struct {
	r io.Reader
	n atomic.Int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n.Add(int64(n))

	return n, err //nolint:wrapcheck // pass-through reader
}

func objectKey(prefix string, t time.Time, encrypted bool) string {
	t = t.UTC()
	key := prefix + t.Format("2006/01/02/") + namePrefix + t.Format("20060102T150405Z") + nameSuffix
	if encrypted {
		key += ageSuffix
	}

	return key
}

func isBackupKey(key string) bool {
	base := path.Base(key)

	return strings.HasPrefix(base, namePrefix) &&
		(strings.HasSuffix(base, nameSuffix) || strings.HasSuffix(base, nameSuffix+ageSuffix))
}

// backups returns the backup objects under the prefix, newest first.
func (b *backupper) backups(ctx context.Context) ([]object, error) {
	objs, err := b.store.List(ctx, b.prefix)
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	out := objs[:0]

	for _, o := range objs {
		if isBackupKey(o.key) {
			out = append(out, o)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].key > out[j].key })

	return out, nil
}

func (b *backupper) touchHealth() {
	if err := os.WriteFile(b.healthFile, []byte(b.now().UTC().Format(time.RFC3339)), healthFileMode); err != nil {
		b.log.Warn("write health file", "path", b.healthFile, "error", err)
	}
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err() //nolint:wrapcheck // plain context error
	case <-t.C:
		return nil
	}
}

// runOnce dumps, uploads and verifies one backup (retrying with exponential backoff), then prunes old ones.
func (b *backupper) runOnce(ctx context.Context) error {
	key := objectKey(b.prefix, b.now(), b.recipient != nil)
	start := time.Now()

	var err error

	for attempt := 1; attempt <= b.attempts; attempt++ {
		if err = b.upload(ctx, key); err == nil {
			break
		}

		b.log.Warn("backup attempt failed", "attempt", attempt, "of", b.attempts, "key", key, "error", err)

		if attempt < b.attempts {
			if serr := sleep(ctx, b.backoff<<(attempt-1)); serr != nil {
				return errors.Join(err, serr)
			}
		}
	}

	if err != nil {
		return fmt.Errorf("backup %s: %w", key, err)
	}

	b.touchHealth()
	b.log.Info("backup complete", "key", key, "duration", time.Since(start).String())

	if perr := b.prune(ctx, key); perr != nil {
		b.log.Warn("retention pruning failed", "error", perr)
	}

	return nil
}

func (b *backupper) upload(ctx context.Context, key string) error {
	pr, pw := io.Pipe()
	cr := &countingReader{r: pr}
	errc := make(chan error, 1)

	go func() {
		err := b.dumpTo(ctx, pw)
		_ = pw.CloseWithError(err)
		errc <- err
	}()

	putErr := b.store.Put(ctx, key, cr)
	_ = pr.CloseWithError(putErr)
	dumpErr := <-errc

	if putErr != nil {
		return fmt.Errorf("upload: %w", putErr)
	}

	if dumpErr != nil {
		return fmt.Errorf("dump: %w", dumpErr)
	}

	size, err := b.store.Stat(ctx, key)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	if want := cr.n.Load(); size != want {
		_ = b.store.Delete(context.WithoutCancel(ctx), []string{key})

		return fmt.Errorf("verify: stored %d bytes, uploaded %d", size, want)
	}

	b.log.Info("upload verified", "key", key, "bytes", size)

	return nil
}

// dumpTo streams DUMP DATABASE through gzip and optional age encryption into w.
// An error leaves the stream unterminated so the uploader never completes a truncated object.
func (b *backupper) dumpTo(ctx context.Context, w io.Writer) error {
	enc, err := encryptWriter(w, b.recipient)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	gz := gzip.NewWriter(enc)
	stmts := 0

	err = b.db.Dump(ctx, func(stmt string) error {
		stmts++
		_, werr := io.WriteString(gz, stmt+"\n")

		return werr //nolint:wrapcheck // pipe write error
	})
	if err != nil {
		return err
	}

	if stmts == 0 {
		return errors.New("database dump is empty, refusing to upload")
	}

	if err = gz.Close(); err != nil {
		return fmt.Errorf("gzip: %w", err)
	}

	if err = enc.Close(); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	b.log.Info("dump streamed", "statements", stmts)

	return nil
}

// prune deletes backups beyond the retention count or older than the retention age, never the newest or keep.
func (b *backupper) prune(ctx context.Context, keep string) error {
	if b.retentionCount == 0 && b.retentionDays == 0 {
		return nil
	}

	all, err := b.backups(ctx)
	if err != nil {
		return err
	}

	cutoff := b.now().AddDate(0, 0, -b.retentionDays)

	var del []string

	for i, o := range all {
		if i == 0 || o.key == keep {
			continue
		}

		if (b.retentionCount > 0 && i >= b.retentionCount) || (b.retentionDays > 0 && o.modified.Before(cutoff)) {
			del = append(del, o.key)
		}
	}

	if len(del) == 0 {
		return nil
	}

	if err = b.store.Delete(ctx, del); err != nil {
		return fmt.Errorf("delete old backups: %w", err)
	}

	b.log.Info("pruned old backups", "deleted", len(del), "kept", len(all)-len(del))

	return nil
}
