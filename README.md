# ShareX Server

## Installing

### Running as a single binary

```bash
git clone https://github.com/shinomineko/sharex-server.git
cd sharex-server
go build

export SHAREX_UPLOAD_KEY="<your-secret-key>"
./sharex-server
```

The server runs on port 3939 by default

### Running as a container

Example `compose.yml`:

```yaml
name: sharex-server
services:
  sharex-server:
    image: docker.io/shinomineko/sharex-server:main
    restart: unless-stopped
    ports:
      - 3939:3939
    volumes:
      - ./uploads:/app/uploads
    environment:
      SHAREX_UPLOAD_KEY: <your-secret-key>
```

## ShareX configuration

1. Go to "Destinations" -> "Custom uploader settings"
2. Create a new uploader with these settings:
   - Request URL: `http://<your-server>:3939/upload`
   - Request method: POST
   - Body: Form data
   - File form name: `file`
   - Headers:
     - Name: `Authorization`
     - Value: `Bearer <your-secret-key>`
