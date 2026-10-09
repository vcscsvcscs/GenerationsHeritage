package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const stmtA = "CREATE (:A);"

var errBoom = errors.New("boom")

type fakeDB struct {
	stmts    []string
	dumpErr  error // returned after emitting all statements
	executed []string
	nonEmpty bool
	wiped    bool
}

func (f *fakeDB) Dump(_ context.Context, emit func(string) error) error {
	for _, s := range f.stmts {
		if err := emit(s); err != nil {
			return err
		}
	}

	return f.dumpErr
}

func (f *fakeDB) Exec(_ context.Context, stmt string) error {
	f.executed = append(f.executed, stmt)

	return nil
}

func (f *fakeDB) IsEmpty(context.Context) (bool, error) { return !f.nonEmpty, nil }

func (f *fakeDB) Wipe(context.Context) error {
	f.wiped = true

	return nil
}

type fakeStore struct {
	objs     map[string][]byte
	modified map[string]time.Time
	deleted  []string
	putFails int // number of Put calls to fail before succeeding
	puts     int
	statSize int64 // when non-zero, reported instead of the real size
}

func newFakeStore() *fakeStore {
	return &fakeStore{objs: map[string][]byte{}, modified: map[string]time.Time{}}
}

// Put mimics a multipart uploader: the object only exists if the reader reached EOF without error.
func (f *fakeStore) Put(_ context.Context, key string, r io.Reader) error {
	f.puts++

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if f.putFails > 0 {
		f.putFails--

		return errBoom
	}

	f.objs[key] = data

	return nil
}

func (f *fakeStore) Stat(_ context.Context, key string) (int64, error) {
	if f.statSize != 0 {
		return f.statSize, nil
	}

	d, ok := f.objs[key]
	if !ok {
		return 0, errors.New("not found")
	}

	return int64(len(d)), nil
}

func (f *fakeStore) List(_ context.Context, prefix string) ([]object, error) {
	var out []object

	for k, d := range f.objs {
		if strings.HasPrefix(k, prefix) {
			out = append(out, object{key: k, size: int64(len(d)), modified: f.modified[k]})
		}
	}

	return out, nil
}

func (f *fakeStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	d, ok := f.objs[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return io.NopCloser(bytes.NewReader(d)), nil
}

func (f *fakeStore) Delete(_ context.Context, keys []string) error {
	for _, k := range keys {
		delete(f.objs, k)
		f.deleted = append(f.deleted, k)
	}

	return nil
}

var fixedNow = time.Date(2025, 6, 7, 3, 0, 5, 0, time.UTC)

func newTestBackupper(t *testing.T, db *fakeDB, st *fakeStore) *backupper {
	t.Helper()

	return &backupper{
		db:         db,
		store:      st,
		log:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		prefix:     "memgraph/",
		attempts:   3,
		backoff:    time.Millisecond,
		healthFile: filepath.Join(t.TempDir(), "last_success"),
		now:        func() time.Time { return fixedNow },
		identity:   func() (age.Identity, error) { return nil, errors.New("no identity") },
	}
}

func gunzip(t *testing.T, data []byte) string {
	t.Helper()

	zr, err := gzip.NewReader(bytes.NewReader(data))
	require.NoError(t, err)

	out, err := io.ReadAll(zr)
	require.NoError(t, err)

	return string(out)
}

const wantKey = "memgraph/2025/06/07/gheritage-20250607T030005Z.cypher.gz"

func TestRunOnceUploadsGzippedDump(t *testing.T) {
	db := &fakeDB{stmts: []string{"CREATE (:A {n: 1});", "CREATE (:B);"}}
	st := newFakeStore()
	b := newTestBackupper(t, db, st)

	require.NoError(t, b.runOnce(t.Context()))

	require.Contains(t, st.objs, wantKey)
	assert.Equal(t, "CREATE (:A {n: 1});\nCREATE (:B);\n", gunzip(t, st.objs[wantKey]))
	assert.FileExists(t, b.healthFile)
}

func TestRunOnceEncryptedRoundTrip(t *testing.T) {
	id, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	db := &fakeDB{stmts: []string{"CREATE (:Secret);"}}
	st := newFakeStore()
	b := newTestBackupper(t, db, st)
	b.recipient = id.Recipient()
	b.identity = func() (age.Identity, error) { return id, nil }

	require.NoError(t, b.runOnce(t.Context()))

	key := wantKey + ".age"
	require.Contains(t, st.objs, key)
	assert.NotContains(t, string(st.objs[key]), "Secret")

	target := &fakeDB{}
	b.db = target
	require.NoError(t, b.restore(t.Context(), "latest", false))
	assert.Equal(t, []string{"CREATE (:Secret)"}, target.executed)
}

func TestRunOnceEncryptedWithoutIdentityFailsRestore(t *testing.T) {
	r, err := age.NewScryptRecipient("pw")
	require.NoError(t, err)
	r.SetWorkFactor(10)

	st := newFakeStore()
	b := newTestBackupper(t, &fakeDB{stmts: []string{stmtA}}, st)
	b.recipient = r
	require.NoError(t, b.runOnce(t.Context()))

	require.ErrorContains(t, b.restore(t.Context(), "latest", false), "no identity")
}

func TestDumpFailureUploadsNothing(t *testing.T) {
	db := &fakeDB{stmts: []string{stmtA}, dumpErr: errBoom}
	st := newFakeStore()
	b := newTestBackupper(t, db, st)

	err := b.runOnce(t.Context())

	require.ErrorIs(t, err, errBoom)
	assert.Empty(t, st.objs)
	assert.Equal(t, 3, st.puts, "every attempt retried")
	assert.NoFileExists(t, b.healthFile)
}

func TestEmptyDumpIsRejected(t *testing.T) {
	st := newFakeStore()
	b := newTestBackupper(t, &fakeDB{}, st)

	require.ErrorContains(t, b.runOnce(t.Context()), "empty")
	assert.Empty(t, st.objs)
}

func TestUploadRetriesTransientFailure(t *testing.T) {
	st := newFakeStore()
	st.putFails = 2
	b := newTestBackupper(t, &fakeDB{stmts: []string{stmtA}}, st)

	require.NoError(t, b.runOnce(t.Context()))

	assert.Equal(t, 3, st.puts)
	assert.Contains(t, st.objs, wantKey)
}

func TestSizeMismatchFailsAndDeletesObject(t *testing.T) {
	st := newFakeStore()
	st.statSize = 1
	b := newTestBackupper(t, &fakeDB{stmts: []string{stmtA}}, st)

	require.ErrorContains(t, b.runOnce(t.Context()), "verify")
	assert.Empty(t, st.objs)
}

func TestPruneKeepsNewestAndIgnoresForeignObjects(t *testing.T) {
	st := newFakeStore()
	b := newTestBackupper(t, &fakeDB{}, st)
	b.retentionCount = 2

	for _, d := range []string{"01", "02", "03", "04"} {
		st.objs["memgraph/2025/06/"+d+"/gheritage-2025060"+d[1:]+"T000000Z.cypher.gz"] = []byte("x")
	}

	st.objs["memgraph/notes.txt"] = []byte("x")

	require.NoError(t, b.prune(t.Context(), ""))

	assert.ElementsMatch(t, []string{
		"memgraph/2025/06/01/gheritage-20250601T000000Z.cypher.gz",
		"memgraph/2025/06/02/gheritage-20250602T000000Z.cypher.gz",
	}, st.deleted)
	assert.Contains(t, st.objs, "memgraph/notes.txt")
}

func TestPruneByAgeKeepsAtLeastNewest(t *testing.T) {
	st := newFakeStore()
	b := newTestBackupper(t, &fakeDB{}, st)
	b.retentionDays = 7

	old := fixedNow.AddDate(0, 0, -30)
	for _, k := range []string{
		"memgraph/2025/05/01/gheritage-20250501T000000Z.cypher.gz",
		"memgraph/2025/05/02/gheritage-20250502T000000Z.cypher.gz.age",
	} {
		st.objs[k] = []byte("x")
		st.modified[k] = old
	}

	require.NoError(t, b.prune(t.Context(), ""))

	assert.Equal(t, []string{"memgraph/2025/05/01/gheritage-20250501T000000Z.cypher.gz"}, st.deleted)
	assert.Len(t, st.objs, 1)
}

func TestRestoreRefusesNonEmptyWithoutForce(t *testing.T) {
	st := newFakeStore()
	b := newTestBackupper(t, &fakeDB{stmts: []string{stmtA}}, st)
	require.NoError(t, b.runOnce(t.Context()))

	target := &fakeDB{nonEmpty: true}
	b.db = target

	require.ErrorIs(t, b.restore(t.Context(), "latest", false), errNotEmpty)
	assert.Empty(t, target.executed)
	assert.False(t, target.wiped)

	require.NoError(t, b.restore(t.Context(), wantKey, true))
	assert.True(t, target.wiped)
	assert.Equal(t, []string{strings.TrimSuffix(stmtA, ";")}, target.executed)
}

func TestRestoreMissingKeyDoesNotWipe(t *testing.T) {
	target := &fakeDB{nonEmpty: true}
	b := newTestBackupper(t, target, newFakeStore())

	require.Error(t, b.restore(t.Context(), "memgraph/nope.cypher.gz", true))
	assert.False(t, target.wiped)
}

func TestObjectKey(t *testing.T) {
	assert.Equal(t, wantKey, objectKey("memgraph/", fixedNow.In(time.FixedZone("x", 7200)), false))
	assert.Equal(t, "2025/06/07/gheritage-20250607T030005Z.cypher.gz.age", objectKey("", fixedNow, true))
}

func TestLoadConfig(t *testing.T) {
	env := map[string]string{"S3_BUCKET": "b", "S3_ACCESS_KEY_ID": "k", "S3_SECRET_ACCESS_KEY": "s", "S3_PREFIX": "/x"}
	get := func(k string) string { return env[k] }

	c, err := loadConfig(get)
	require.NoError(t, err)
	assert.Equal(t, "x/", c.s3Prefix)
	assert.Equal(t, "auto", c.s3Region)
	assert.Equal(t, defaultRetention, c.retentionCount)

	env["BACKUP_SCHEDULE"] = "nonsense"
	env["S3_BUCKET"] = ""
	env["BACKUP_RETENTION_DAYS"] = "-1"
	_, err = loadConfig(get)
	require.ErrorContains(t, err, "S3_BUCKET is required")
	require.ErrorContains(t, err, "BACKUP_SCHEDULE")
	require.ErrorContains(t, err, "BACKUP_RETENTION_DAYS")
}

func TestHealthcheck(t *testing.T) {
	file := filepath.Join(t.TempDir(), "last_success")
	get := func(k string) string { return map[string]string{"BACKUP_HEALTH_FILE": file}[k] }

	require.Error(t, healthcheck(get), "missing file is unhealthy")

	require.NoError(t, os.WriteFile(file, nil, 0o600))
	require.NoError(t, healthcheck(get))

	stale := time.Now().Add(-72 * time.Hour)
	require.NoError(t, os.Chtimes(file, stale, stale))
	require.ErrorContains(t, healthcheck(get), "last successful backup")
}
