package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
)

var errNotEmpty = errors.New("database is not empty, use --force to wipe it before restoring")

func (b *backupper) resolveKey(ctx context.Context, target string) (string, error) {
	if target != "latest" {
		return target, nil
	}

	all, err := b.backups(ctx)
	if err != nil {
		return "", err
	}

	if len(all) == 0 {
		return "", fmt.Errorf("no backups found under prefix %q", b.prefix)
	}

	return all[0].key, nil
}

// open returns the decrypted, decompressed dump stream of the object and a function that closes it.
func (b *backupper) open(ctx context.Context, key string) (io.Reader, func(), error) {
	rc, err := b.store.Get(ctx, key)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s: %w", key, err)
	}

	var r io.Reader = rc

	if strings.HasSuffix(key, ageSuffix) {
		id, ierr := b.identity()
		if ierr != nil {
			_ = rc.Close()

			return nil, nil, ierr
		}

		if r, err = age.Decrypt(r, id); err != nil {
			_ = rc.Close()

			return nil, nil, fmt.Errorf("decrypt %s: %w", key, err)
		}
	}

	gz, err := gzip.NewReader(r)
	if err != nil {
		_ = rc.Close()

		return nil, nil, fmt.Errorf("gunzip %s: %w", key, err)
	}

	return gz, func() { _ = gz.Close(); _ = rc.Close() }, nil
}

// restore replays the dump of target (an object key or "latest") into an empty database.
func (b *backupper) restore(ctx context.Context, target string, force bool) error {
	empty, err := b.db.IsEmpty(ctx)
	if err != nil {
		return fmt.Errorf("check database: %w", err)
	}

	if !empty && !force {
		return errNotEmpty
	}

	key, err := b.resolveKey(ctx, target)
	if err != nil {
		return err
	}

	r, closeFn, err := b.open(ctx, key)
	if err != nil {
		return err
	}
	defer closeFn()

	if !empty {
		b.log.Warn("wiping non-empty database before restore")

		if err = b.db.Wipe(ctx); err != nil {
			return fmt.Errorf("wipe database: %w", err)
		}
	}

	stmts, err := b.replay(ctx, r)
	if err != nil {
		return fmt.Errorf("restore %s failed after %d statements (database is partially restored, re-run with --force): %w", key, stmts, err)
	}

	b.log.Info("restore complete", "key", key, "statements", stmts)

	return nil
}

func (b *backupper) replay(ctx context.Context, r io.Reader) (int, error) {
	br := bufio.NewReader(r)
	stmts := 0

	for {
		line, err := br.ReadString('\n')
		if stmt := strings.TrimSuffix(strings.TrimSpace(line), ";"); stmt != "" {
			if xerr := b.db.Exec(ctx, stmt); xerr != nil {
				return stmts, fmt.Errorf("statement %d: %w", stmts+1, xerr)
			}

			stmts++
		}

		if errors.Is(err, io.EOF) {
			return stmts, nil
		}

		if err != nil {
			return stmts, fmt.Errorf("read dump: %w", err)
		}
	}
}
