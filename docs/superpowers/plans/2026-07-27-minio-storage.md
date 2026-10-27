# Plan: MinIO Storage + Konversi WebP

**Spec:** [2026-10-27-minio-storage-design.md](../specs/2026-10-27-minio-storage-design.md)
**Date:** 2026-10-27
**Branch:** `feat/minio-storage` (dari `main`)

## Architecture

Seam storage diganti dari `common/gcs` (GCS, berbayar) ke `common/storage` (MinIO
self-hosted), interface method `UploadFile` tetap sehingga service logic gak berubah.
Konversi gambar (resize max 1600px → webp lossless) dipasang di service sebelum upload.

Dependency baru (semua pure Go, `CGO_ENABLED=0` tetap aman):
- `github.com/minio/minio-go/v7`
- `github.com/HugoSmits86/nativewebp`
- `golang.org/x/image` (draw + webp decode)

## Execution order

**1 → 2 → 4 → 3 → 5 → 6 → 7 → 8** (nomor = nomor task di bawah). Tiap task satu
commit, dan tiap commit harus `go build ./...` EXIT=0:

- Task 4 (config, **additive** — tambah `MinIOConfig` saja) harus sebelum task 3
  (storage client yang pakai struct itu).
- Task 5 (rename caller + hapus `common/gcs` + **hapus** field `GCS*`) jadi titik
  di mana config lama ikut mati — di situlah config boleh berkurang, karena
  semua caller-nya sudah pindah di commit yang sama.

## Global Constraints

- Bahasa: Indonesia untuk pesan commit; kode + identifier bahasa Inggris.
- Konvensi repo: interface `I`-prefix (`IStorageClient`), constructor `New*`
  mengembalikan interface (lihat `NewGCSClient` lama).
- `go build ./...` dan `go vet ./...` harus EXIT=0 sebelum commit task apa pun.
- Dockerfile **tidak boleh diubah** (tetap `CGO_ENABLED=0`) — karena itu semua
  dep harus pure Go.
- Style commit dari plan ini jangan disingkat.
- PowerShell: `\` bukan line continuation; keluaran `git push` "Exited with code 1"
  itu noise stderr — cek teks outputnya, bukan exit code.

## Tasks

### Task 1: Tambah dependencies

**Files:**
- Modify: `go.mod`, `go.sum`

**Interfaces:**
- (none — hanya `go get`)

**Steps:**

1. Jalankan:
   ```powershell
   go get github.com/minio/minio-go/v7@latest
   go get github.com/HugoSmits86/nativewebp@latest
   go get golang.org/x/image@latest
   ```
2. Verifikasi: `go mod tidy` EXIT=0, `go build ./...` EXIT=0.

**Test:** Tidak ada test baru (task murni deps).
**Commit message:** `feat: tambah deps minio-go, nativewebp, golang.org/x/image`

---

### Task 2: Konversi gambar — `services/field/image.go`

**Files:**
- Create: `services/field/image.go`
- Create: `services/field/image_test.go`

**Interfaces:**
- `func convertToWebP(data []byte) ([]byte, error)` — unexported helper:
  1. `image.Decode` (register format via blank import `_ "image/jpeg"`, `_ "image/png"`,
     `_ "golang.org/x/image/webp"` di file yang sama).
  2. Kalau sisi terpanjang > 1600px → resize proporsional pakai `draw.CatmullRom.Scale`
     ke `1600×(h*1600/w)` atau `(w*1600/h)×1600` — sesuaikan orientasi.
  3. `nativewebp.Encode` → keluar `[]byte` webp (lossless).
  4. Error decode/encode → return error apa adanya.

**Steps (TDD):**

1. **FAIL dulu** — tulis `services/field/image_test.go`:
   - Bikin JPEG sintetis 2000×1500 (`image.NewRGBA` + fill warna + `jpeg.Encode`).
   - Panggil `convertToWebP(got)`.
   - Assert: `err == nil`, magic bytes `RIFF....WEBP` (`bytes.HasPrefix(data[:12], []byte("RIFF"))` + `string(data[8:12]) == "WEBP"`), dan sisi terpanjang ≤ 1600px
     (decode ulang via `image.Decode`).
   - `go test ./services/field/` → **FAIL** (file belum ada / belum compile).
2. **PASS** — tulis `services/field/image.go` sesuai kontrak di atas.
3. `go test ./services/field/` → PASS.

**Commit message:** `feat: konversi gambar resize max 1600px + encode webp lossless`

---

### Task 3: `common/storage` — interface + MinIO client

**Files:**
- Create: `common/storage/storage.go` — interface `IStorageClient`
- Create: `common/storage/minio.go` — `MinIOClient`
- Create: `common/storage/minio_test.go`

**Interfaces:**
- `type IStorageClient interface { UploadFile(ctx context.Context, filename string, data []byte) (string, error) }`
  — signature **persis** sama dengan `IGCSClient` lama (biar rename di task 5 saja).
- `func NewMinIOClient(cfg config.MinIOConfig) (IStorageClient, error)`:
  1. `minio.New(cfg.Endpoint, &minio.Options{Cred: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""), Secure: cfg.UseSSL})`.
  2. `BucketExists` → kalau tidak ada: `MakeBucket` → `SetBucketPolicy` dengan policy JSON:
     `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::field-images/*"]}]}`
     (Resource pakai `cfg.BucketName`).
  3. Error di step mana pun → return error.
- `UploadFile`: `PutObject` (size `int64(len(data))`, contentType `http.DetectContentType(data)`)
  → return `fmt.Sprintf("%s/%s/%s", cfg.PublicURL, cfg.BucketName, filename)`.

**Steps (TDD):**

1. Pastikan **task 4 (config)** sudah dikerjakan lebih dulu — struct
   `config.MinIOConfig` + key `minio` di config example harus sudah ada.
2. **FAIL dulu** — tulis `common/storage/minio_test.go`:
   - `TestUploadFile`: skip kalau `localhost:9000` tidak terhubung (dial timeout 1s
     via `net.DialTimeout`). Kalau nyala: `UploadFile(ctx, "images/test-<ts>.webp", []byte("RIFF....WEBP"))`,
     lalu `http.Get(url)` → assert status 200 + body identik.
3. **PASS** — tulis `storage.go` + `minio.go` sesuai kontrak.
4. `go test ./common/storage/` → PASS (atau SKIP kalau MinIO mati — skip itu expected).

**Commit message:** `feat: seam storage IStorageClient + client MinIO auto-bucket public`

---

### Task 4: Config — tambah `minio` (additive saja)

**Files:**
- Modify: `config/config.go`
- Modify: `config.json.example`
- Modify: `config.docker.example.json`

**Interfaces:**
- `type MinIOConfig struct { Endpoint string \`json:"endpoint"\`; PublicURL string \`json:"publicURL"\`; AccessKey string \`json:"accessKey"\`; SecretKey string \`json:"secretKey"\`; UseSSL bool \`json:"useSSL"\`; BucketName string \`json:"bucketName"\` }`
- Field `AppConfig`: **tambah saja** `MinIO MinIOConfig \`json:"minio"\``.
  Penghapusan field `GCS*` + key `gcs*` di file example **bukan di task ini** —
  pindah ke task 5 biar tiap commit tetap compile (caller `initGCS` masih hidup
  sampai task 5).

**Values:**

`config.json.example` (local) — tambahkan blok:
```json
"minio": {
  "endpoint": "localhost:9000",
  "publicURL": "http://localhost:9000",
  "accessKey": "minioadmin",
  "secretKey": "minioadmin",
  "useSSL": false,
  "bucketName": "field-images"
}
```

`config.docker.example.json`: sama, tapi `endpoint` = `minio:9000`,
`publicURL` = `http://localhost:9000` (di-sed Jenkins → `http://<HOST>:9000`).

**Steps:**

1. Edit `config/config.go`: tambah struct `MinIOConfig` + field `MinIO`.
2. Tambah blok `minio` di kedua file example.
3. `go build ./...` EXIT=0 (additive — gak ada yang rusak).

**Commit message:** `feat: config minio (endpoint/publicURL/keys/bucket)`

---

### Task 5: Hapus `common/gcs` + config `GCS*`, rename caller → `storage`

**Files:**
- Delete: `common/gcs/gcs.go`
- Modify: `config/config.go` (hapus field `GCS*`)
- Modify: `config.json.example`, `config.docker.example.json` (hapus key `gcs*` + `bucketName` lama GCS)
- Modify: `services/field/field.go` (import + tipe field `gcs` → `storage`)
- Modify: `services/registry.go` (`gcs.IGCSClient` → `storage.IStorageClient`)
- Modify: `cmd/main.go` (`initGCS` → `initStorage`, pakai `storage.NewMinIOClient`)

**Interfaces:**
- `cmd/main.go`: ganti isi `initGCS()` → `initStorage()` yang memanggil
  `storage.NewMinIOClient(config.Config.MinIO)`; kalau error → `logrus.Fatal`/panic
  (konsisten dengan gaya lama `initGCS`).

**Steps:**

1. Grep dulu pemakaian `config.Config.GCS*` + `config.Config.BucketName` (GCS) di
   seluruh repo — semua harus ikut mati di commit ini (biasanya cuma dipakai
   `initGCS`).
2. Delete `common/gcs/gcs.go`, hapus field `GCS*` + key example, update 3 file
   caller (import/type/constructor).
3. Grep sisa referensi `gcs`/`GCS` di `*.go` → harus nol (kecuali komentar yang
   memang relevan).
4. `go build ./...` + `go vet ./...` EXIT=0.

**Commit message:** `refactor: seam storage IStorageClient, hapus common/gcs`

---

### Task 6: Wire konversi webp di `processAndUploadImages`

**Files:**
- Modify: `services/field/field.go`

**Interfaces:**
- (memakai `convertToWebP` dari task 2)

**Steps:**

1. Di loop `processAndUploadImages`, setelah `fileHeader.Open()` + baca `fileData`,
   sebelum `storage.UploadFile`:
   ```go
   webpData, err := convertToWebP(fileData)
   if err != nil { return nil, err }
   ```
   lalu `UploadFile(..., webpData)`.
2. Pola filename ganti dari `images/<ts>-<nama+ext>-<ext>` (double ext) menjadi:
   ```go
   ext := ".webp"
   base := strings.TrimSuffix(path.Base(fileHeader.Filename), path.Ext(fileHeader.Filename))
   fileName := "images/" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + base + ext
   ```
   (hapus perhitungan `randomInt`/`ext` lama yang tidak terpakai).
3. `go build ./...` EXIT=0.

**Commit message:** `feat: upload field gambar dikonversi webp + nama file rapi`

---

### Task 7: Compose + Jenkins + config docker

**Files:**
- Modify: `docker-compose.yml`
- Modify: `Jenkinsfile`

**Interfaces:**
- (infra, bukan Go)

**Steps:**

1. `docker-compose.yml`:
   - Tambah service:
     ```yaml
     minio:
       image: minio/minio:latest
       container_name: minio
       command: server /data --console-address :9001
       ports:
         - "9000:9000"
         - "9001:9001"
       environment:
         MINIO_ROOT_USER: minioadmin
         MINIO_ROOT_PASSWORD: minioadmin
       volumes:
         - miniodata:/data
     ```
   - `volumes:` root tambah `miniodata:`.
   - Service `field-service` tambah `depends_on: [minio]` (dan `postgres` kalau ada).
2. `Jenkinsfile` deploy stage — tambah sebelum `docker compose up -d`:
   ```bash
   [ -f config.docker.json ] || cp config.docker.example.json config.docker.json
   sed -i 's|"publicURL": "http://localhost:9000"|"publicURL": "http://${HOST}:9000"|' config.docker.json
   ```
3. Verifikasi: `docker compose config` EXIT=0.

**Commit message:** `feat: service minio di compose + siapkan config.docker.json di deploy`

---

### Task 8: Verifikasi akhir — full build + test + push

**Files:**
- (verifikasi saja)

**Steps:**

1. `go build ./...` EXIT=0.
2. `go vet ./...` EXIT=0.
3. `go test ./...` — `services/field` (webp) PASS; `common/storage` SKIP/PASS
   (MinIO lokal mati = skip, nyala = pass).
4. `docker compose config` EXIT=0.
5. Merge/PR ke `main` sesuai alur branch (`feat/minio-storage` → `main`).

**Commit message:** (tidak ada — verifikasi; kalau ada perbaikan kecil: `fix: ...`)
