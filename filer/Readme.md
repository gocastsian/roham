# File Handling Microservices

This project provides two file server (upload and download) for efficient file handling:

- **Upload Server**: Handles resumable file uploads using the `tusd` protocol, storing files in a MinIO object storage backend.
- **File Server**: Manages file downloads and metadata retrieval from Filesystem/S3.

## Features

- Resumable file uploads with `tusd` for handling large files and network interruptions.
- Scalable file storage with S3 compatible object storage (minio).
- RESTful API for file downloads and metadata retrieval.


## Run upload server

1.**Start Services with Docker Compose**

   ```bash
   
   make start-filer-dev
   
   #start upload server
   go run cmd/filer/main.go upload
  
   #run upload test
   go test -v -count=1 ./filer/upload/e2e
   ```

3. **Access MinIO UI**

    - URL: `http://localhost:9000`
    - Credentials: `minioadmin/minioadmin`
    - Create a bucket (e.g., `temp-storage`) for storing uploads.


## API Endpoints

### Uploader Service (`http://localhost:5006`)

- `POST /files`: Initiate a file upload (tus protocol).
- `PATCH /uploads/<file-id>`: Upload chunks for an existing upload.
- `HEAD /uploads/<file-id>`: Check upload status.

### Downloader Service (`http://localhost:5005`)

- `GET /api/v1/files/:key/download`: Direct download.
- `GET /api/v1/files/:key/download`: Download using pre-signed-url.

## Futures

- [x] Upload files using tusd handler.(filesystem store)
- [x] Direct upload to s3/minio.
- [x] Download a file by ID.
- [x] Implementing a repository and store metadata in a relational database.
- [ ] Validating uploads base on specific upload type (avatars, layers and ...).


