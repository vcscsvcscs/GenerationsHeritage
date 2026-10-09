#!/usr/bin/env bash
# Seed memgraph-src, back up to MinIO, restore into memgraph-dst and compare counts (plain and encrypted).
set -euo pipefail

cd "$(dirname "$0")/../../.."
dc() { docker compose -f compose.backup-test.yaml "$@"; }
trap 'dc down -v' EXIT

mg() { echo "$2" | dc exec -T "$1" mgconsole --username=memgraph --password=memgraph --output-format=csv | tail -n +2; }
counts() { echo "$(mg "$1" 'MATCH (n) RETURN count(n);') $(mg "$1" 'MATCH ()-[r]->() RETURN count(r);')"; }

dc up -d --build --wait memgraph-src memgraph-dst minio
dc run --rm mc mb --ignore-existing m/gheritage-backups

mg memgraph-src 'CREATE INDEX ON :Person(name);'
mg memgraph-src 'CREATE (a:Person {name: "Anna", note: "line1\nline2 \"quoted\" ; end"}), (b:Person {name: "Bela"}), (c:Person {name: "Cili"}),
  (a)-[:PARENT_OF {since: 1990}]->(b), (b)-[:PARENT_OF]->(c), (a)-[:GRANDPARENT_OF]->(c);'
want=$(counts memgraph-src)
echo "source counts (nodes rels): $want"

dc run --rm backup run-once
dc run --rm mc ls --recursive m/gheritage-backups | tee /dev/stderr | grep -q 'gheritage-.*\.cypher\.gz$'

dc run --rm -e MEMGRAPH_URI=bolt://memgraph-dst:7687 backup restore latest
got=$(counts memgraph-dst)
[ "$got" = "$want" ] || { echo "FAIL plain restore: got '$got', want '$want'"; exit 1; }

dc run --rm -e MEMGRAPH_URI=bolt://memgraph-dst:7687 backup restore latest && { echo "FAIL: restore into non-empty db must be refused"; exit 1; }

sleep 1
dc run --rm -e BACKUP_ENCRYPTION_PASSPHRASE=e2e-secret backup run-once
dc run --rm mc ls --recursive m/gheritage-backups | grep -q '\.cypher\.gz\.age$'
dc run --rm -e MEMGRAPH_URI=bolt://memgraph-dst:7687 -e BACKUP_ENCRYPTION_PASSPHRASE=e2e-secret backup restore --force latest
got=$(counts memgraph-dst)
[ "$got" = "$want" ] || { echo "FAIL encrypted restore: got '$got', want '$want'"; exit 1; }

echo "e2e OK: $got"
