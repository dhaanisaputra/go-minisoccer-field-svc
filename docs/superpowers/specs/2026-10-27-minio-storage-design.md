# Design: Ganti storage gambar GCS → MinIO

Tanggal: 2026-10-27
Status: menunggu approval

## Intent

Foto lapangan (upload di `FieldService`) disimpan ke MinIO yang di-host sendiri —
gratis, jalan di dev lokal, dan tetap kebaca dari luar setelah deploy ke server.
GCS (berbayar) ditinggalkan.

## Goals

- Upload lewat MinIO, URL publik yang di-return `UploadFile` disimpan ke DB `field.images`.
- Dua mode dalam satu config: dev lokal (`localhost:9000`) dan deploy
  (`endpoint=minio:9000` di dalam network compose, `publicURL=http://<server>:9000`
  buat browser).
- Bucket dibuat + di-set public (policy read-only) otomatis pas app start.
- Service/controller logic tidak berubah — cuma rename import/type.

## Non-goals / anti-goals

- Tidak ada migrasi gambar GCS (tabel field masih kosong, ga ada data).
- Tidak ada saklar `storage.type` GCS/MinIO (opsi 3 ditolak — YAGNI; balik ke GCS
  = git revert kapan aja).
- Tidak pakai presigned URL — public bucket sudah dipilih.
- Pola penamaan file di `services/field` (double ekstensi `xxx.png-.png`) tidak
  disentuh — di luar scope.
- Frontend/URL lama: n/a, belum ada data.

## Constraints

- Repo pakai konvensi interface `I`-prefix (`IFieldService`, `IGCSClient`) → nama
  baru ikut konvensi itu.
- `config.docker.json` sudah di-gitignore; Jenkins deploy harus menyediakan file
  itu di server (sekarang belum).
- Port 9000 server harus dibuka di firewall (catatan ops, di luar repo).

## Approaches (dipertimbangkan)

1. **Tukar badan di tempat** — isi ulang `common/gcs` jadi MinIO, nol ubah caller.
   Ditolak: package `gcs` isinya MinIO = bingungin.
2. **Rename seam → `common/storage`** ✅ **dipilih** — nama jujur, blast radius
   kecil (4 tempat, karena rename doang), service logic tak tersentuh.
3. **Saklar config GCS|MinIO** — dua implementasi hidup. Ditolak: YAGNI.

SDK: **`github.com/minio/minio-go/v7`** (client resmi; endpoint custom +
MakeBucket + SetBucketPolicy simpel). Call saya.

## How (keputusan)

- Paket baru `common/storage`:
  - `storage.go`: interface `IStorageClient` — method tetap
    `UploadFile(ctx, filename string, data []byte) (string, error)`.
  - `minio.go`: `MinIOClient` (minio-go). Constructor terima config, bikin client,
    **ensure bucket**: `BucketExists` → `MakeBucket` → `SetBucketPolicy` (policy
    `s3:GetObject` untuk `*` = publik). Error → return error, `cmd/main.go` panic
    (konsisten sama gaya `initGCS` sekarang).
  - `UploadFile`: `PutObject` + `ContentType: http.DetectContentType(data)`,
    return `"{publicURL}/{bucket}/{filename}"`.
- Hapus `common/gcs/gcs.go` + field `GCS*` di `config.AppConfig` + key `gcs*` di
  kedua file example (mati semua; kalau nanti balik GCS, git revert).
- Rename di caller (cuma nama): `services/field` (field + param `gcs`→`storage`),
  `services/registry.go`, `cmd/main.go` (`initGCS`→`initStorage`).
- Config baru (AppConfig `json:"minio"`): `endpoint`, `publicURL`, `accessKey`,
  `secretKey`, `useSSL`, `bucketName`.
  - Local: `localhost:9000` / `http://localhost:9000` / `minioadmin` /
    `minioadmin` / `false` / `field-images`.
  - Docker (config.docker.example.json): `endpoint=minio:9000`,
    `publicURL=http://localhost:9000` (di-sed Jenkins → `http://<HOST>:9000`).
- `docker-compose.yml`: service `minio` (image `minio/minio`, `server /data
  --console-address :9001`, port `9000:9000` + `9001:9001`,
  `MINIO_ROOT_USER/PASSWORD=minioadmin`, volume `miniodata`),
  `field-service.depends_on += minio`.
- `Jenkinsfile` deploy stage tambah 2 baris (sebelum `docker compose up`):
  `[ -f config.docker.json ] || cp config.docker.example.json config.docker.json`
  dan sed `publicURL` → `http://${HOST}:9000`.

## Verification

1. `go build ./...` + `go vet ./...` bersih.
2. `docker compose config` valid.
3. Test kecil `storage/minio_test.go`: skip kalau minio ga nyala; kalau nyala →
   `UploadFile` file dummy, assert URL ke-`GET` balas 200 + isi sesuai.
   (Curl POST end-to-end tetap butuh user-service + JWT — di luar test ini.)

## What's left to builder

Eksekusi persis keputusan di atas, urutan bebas; semua titik sentuh sudah
terdaftar (storage pkg, 3 file rename, config ×3, compose, Jenkins, test).
