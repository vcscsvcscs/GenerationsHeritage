# Generational Heritage - Generációról Generációra: Digitális Családfa és Intellektuális Örökségek Rögzítése
Open Source family tree software that has the objective of keeping digital heritage for future generations.

Nyílt forráskódú családfa-szoftver, amelynek célja a digitális örökség megőrzése a jövő nemzedékei számára.

## Hungarian thesis disclaimer - Szakdolgozat leírása magyarul:
A szakdolgozat célja egy webalkalmazás kifejlesztése, mely lehetővé teszi a digitális családfa közösségi alapú építését és a családtagok intellektuális örökségeinek
rögzítését. Az alkalmazás elérhetővé válna mindenféle eszközről, biztosítva a felhasználók számára a családfájuk személyre szabott kezelését és a családtagok életének
széleskörű dokumentálását.

A felhasználók regisztrálhatnak, és a rendszer lehetőséget biztosít saját profiljuk szerkesztésére, mely része a családfa struktúrájának. A családfában nem csupán a
nevek és születési dátumok lennének elérhetők, hanem további információk is, mint például iskolák, lakhelyek, munkahelyek, életbölcsességek, fontos tudás és fotók.
Az alkalmazás ezenkívül védelmi intézkedéseket alkalmazna, így csak a vérrokonságban állók érnének el egymás adatait.

A feladat magas komplexitással jár, hiszen nemcsak a felhasználói felületet és a családfa struktúrát kell kialakítani, de a biztonsági rétegeket is megfelelően implementálni.
Az adatbázis rendszer, a felhőalapú szerver és a CI/CD rendszer kialakítása további kihívásokat rejt. Az alkalmazásnak a különféle eszközökön és kijelzőméreteken
történő optimális megjelenést kell biztosítania, ami további fejlesztési és tervezési készségeket igényel. A szakdolgozat részleteiben kifejti, hogy a projekthez kapcsolódó
specifikus kihívások és megoldások milyen mértékben járulnak hozzá a szoftver sikeréhez és funkcionalitásához.

## English thesis disclaimer - Szakdolgozat leírása angolul:
The purpose of the thesis is to develop a web application that enables the community-based construction of a digital family tree and the recording of the intellectual heritages of family members. The application would become accessible from all kinds of devices, ensuring personalized management of users' family trees and extensive documentation of family members' lives.

Users could register and the system would provide the ability to edit their own profiles, which are part of the family tree structure. The family tree would contain not only names and birth dates but also additional information such as schools, residences, workplaces, life wisdom, important knowledge, and photos. Furthermore, the application would employ protective measures, so that only those related by blood could access each other's data.

The task is highly complex, as it involves not only designing the user interface and the family tree structure but also properly implementing security layers. The development of the database system, the cloud-based server, and the CI/CD system present further challenges. The application must ensure optimal display on various devices and screen sizes, which requires additional development and design skills. The thesis details the extent to which specific challenges and solutions related to the project contribute to the success and functionality of the software.

## Deployment
To deploy all micro services use:

```bash:
kubectl apply --server-side -k .
```

## Backups
The `backup` service in `compose.yaml` (source: `apps/backup`) dumps Memgraph on a schedule and uploads it to an S3-compatible bucket (Cloudflare R2, MinIO, AWS S3, Backblaze B2).

**Why a logical dump:** Memgraph runs with `ON_DISK_TRANSACTIONAL`, where durability is handled by RocksDB; the snapshot/WAL files (`--storage-snapshot-*`, `CREATE SNAPSHOT`) belong to the in-memory storage modes and are not a reliable backup source here. The service therefore runs `DUMP DATABASE` over Bolt, which works in every storage mode, and streams the Cypher statements through gzip (and optional age encryption) straight into a multipart upload, without buffering the dump in memory or mounting the Memgraph data volume. Each upload is verified by comparing the stored size with the bytes sent. Objects are named `<S3_PREFIX>YYYY/MM/DD/gheritage-<UTC timestamp>.cypher.gz[.age]`.

**Setup (Cloudflare R2):**
1. Create a private bucket (R2 > Create bucket). Never enable public access or a public domain: the dump contains family PII.
2. R2 > Manage API tokens > Create API token, permission `Object Read & Write`, scoped to this bucket only. Copy the access key id, secret and the account endpoint `https://<ACCOUNT_ID>.r2.cloudflarestorage.com`.
3. Fill in the `backup` placeholders in `compose.yaml`, then `docker compose up -d backup`. Check `docker logs backup` and run `docker compose run --rm backup run-once` once to confirm.

| Variable | Default | Meaning |
|---|---|---|
| `MEMGRAPH_URI` / `MEMGRAPH_USER` / `MEMGRAPH_PASSWORD` | `bolt://memgraph:7687` / - / - | Bolt connection |
| `S3_ENDPOINT` | AWS S3 | Endpoint URL of the S3-compatible service |
| `S3_BUCKET`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY` | required | Bucket and credentials |
| `S3_REGION` | `auto` | `auto` for R2; the real region for AWS/Backblaze |
| `S3_PREFIX` | `memgraph/` | Key prefix for backups |
| `S3_PATH_STYLE` | `false` | Force path-style addressing; custom endpoints (R2, MinIO) already use it by default, AWS S3 uses virtual-hosted style |
| `BACKUP_SCHEDULE` | `0 3 * * *` | Standard 5-field cron expression, evaluated in `TZ` |
| `BACKUP_ON_START` | `false` | Also back up at container start |
| `BACKUP_RETENTION_COUNT` | `30` | Keep the newest N backups (`0` disables) |
| `BACKUP_RETENTION_DAYS` | `0` | Also delete backups older than N days (`0` disables) |
| `BACKUP_AGE_RECIPIENT` | - | age public key (`age1...`); encrypts backups client-side. Recommended |
| `BACKUP_ENCRYPTION_PASSPHRASE` | - | Alternative to the public key (symmetric). Mutually exclusive with `BACKUP_AGE_RECIPIENT` |
| `BACKUP_AGE_IDENTITY` | - | age secret key (`AGE-SECRET-KEY-...`), only needed to restore encrypted backups |

Behaviour: failed runs are retried 3 times with exponential backoff, then logged (JSON) while the service keeps running; `run-once` exits non-zero on failure. An empty dump is rejected so a wiped database can never replace good backups. Retention pruning runs only after a successful, verified upload, never deletes the newest backup and ignores objects that do not match the backup naming. The Docker healthcheck (`backup healthcheck`) turns unhealthy when there was no successful backup for two schedule intervals plus one hour.

**Encryption keys:** generate with `age-keygen`; put only the public key on the server and keep the secret key offline. Lost key means unreadable backups.

**Restore** (into an empty Memgraph; the restore refuses a non-empty database, `--force` wipes it first):
```bash
docker compose run --rm backup restore latest              # newest backup
docker compose run --rm backup restore memgraph/2025/06/07/gheritage-20250607T030005Z.cypher.gz.age
docker compose run --rm -e BACKUP_AGE_IDENTITY=AGE-SECRET-KEY-... backup restore --force latest
```
Statements are replayed one by one over Bolt, so large graphs take a while. If a restore fails midway, re-run it with `--force`.

**Operations:**
- Add a bucket lifecycle rule (R2/S3) as a safety net, e.g. expire objects under `memgraph/` after 90 days; it complements, not replaces, `BACKUP_RETENTION_*`.
- Test a restore periodically, e.g. into a scratch Memgraph with `compose.backup-test.yaml`; `apps/backup/e2e/run.sh` runs the full seed, backup, restore and compare cycle against MinIO.
- Unhealthy status is only a Docker flag; wire it to an alert (e.g. Uptime Kuma, autoheal) to be notified.