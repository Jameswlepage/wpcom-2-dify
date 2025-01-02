# Dify-WP-Sync

This project integrates WordPress.com sites with a [Dify](https://dify.ai) dataset. It automatically syncs WordPress posts into a Dify dataset, allowing you to keep textual content up-to-date for use with Dify-based applications.

## Intended Audience

This guide is for developers who have experience with Docker, Go, and basic OAuth flows. If you are unfamiliar with these technologies, it’s recommended to review their documentation before proceeding.

## Prerequisites

1. **Disable Autoproxxy (if applicable)**  
   Some environments or corporate networks may use automatic proxy configurations that interfere with local callback URLs or API requests. Disable these auto-configured proxies or add exceptions for the `boc.local` domain.

2. **Add `boc.local` to Your Hosts File**  
   The OAuth flow requires `http://boc.local:8080` as the callback URL. Set `boc.local` to point to `127.0.0.1` on your machine.

   - **Linux/macOS**:

     ```bash
     sudo nano /etc/hosts
     ```

     Add:

     ```
     127.0.0.1   boc.local
     ```

     Save and close.

   - **Windows**:  
     Edit `C:\Windows\System32\Drivers\etc\hosts` as Administrator:
     ```
     127.0.0.1   boc.local
     ```
     After that, `boc.local` should resolve to `127.0.0.1`.

3. **Docker & Docker Compose (v2)**  
   Make sure you have [Docker](https://www.docker.com/) and Docker Compose v2 installed.

4. **Docker and Docker Compose**  
   Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

5. **Go (Optional)**  
   If you want to run the CLI or server locally without Docker, install [Go 1.20+](https://go.dev/dl/).
   Otherwise, using Docker Compose will handle building and running the binaries.

## Configuration

1. **Environment Variables**  
   The application uses a `.env` file for configuration. A template `.env.example` is provided.

   ```bash
   cp .env.example .env
   ```

   ```bash
   cp .env.example .env
   ```

   Edit `.env`:

   - `WPCOM_CLIENT_ID` and `WPCOM_CLIENT_SECRET`: Your WordPress.com OAuth credentials.
   - `WPCOM_REDIRECT_URI`: Should remain `http://boc.local:8080/oauth/callback`.
   - `DIFY_API_KEY`: Your Dify API key.
   - `DIFY_BASE_URL`: The Dify endpoint (default is `https://api.dify.ai/v1`).

   **Do not commit the `.env` file** as it contains sensitive information.

   **Example `.env` configuration:**

   ```env
   WPCOM_CLIENT_ID=your_wpcom_client_id
   WPCOM_CLIENT_SECRET=your_wpcom_client_secret
   WPCOM_REDIRECT_URI=http://boc.local:8080/oauth/callback
   DIFY_API_KEY=your_dify_api_key
   DIFY_BASE_URL=https://api.dify.ai/v1
   REDIS_ADDR=redis:6379
   REDIS_DB=0
   REDIS_PASSWORD=
   PORT=8080
   ```

2. **WordPress OAuth App Setup**  
   In your WordPress.com OAuth app settings, ensure:
   ```
   Redirect URI: http://boc.local:8080/oauth/callback
   ```
   This must match exactly or the OAuth flow will fail.

## Running the Project

### Using Docker Compose (v2)

1. Build and start the services:

   ```bash
   docker compose up --build
   ```

   This will:

   - Run Redis in a container.
   - Build and run the Go server.

2. After startup, the server listens on:

   ```
   http://boc.local:8080
   ```

   - `GET /` returns a simple status message.
   - `GET /oauth/callback` handles the OAuth process.

3. Authorize a new WordPress site by opening the OAuth authorization URL. The easiest way is:
   ```bash
   docker compose run --rm app ./cli open-oauth
   ```
   Open the displayed URL in your browser, authorize the app, and your site will be registered.

---

The CLI (named `cli`) is built inside the Docker image. Run it via Docker:

```bash
docker compose run --rm app ./cli list-sites
docker compose run --rm app ./cli sync-site <site_id>
docker compose run --rm app ./cli sync-all-sites
```

**CLI Commands:**

- `list-sites`: Lists all registered WordPress sites.
- `sync-site <site_id>`: Syncs a single site by ID.
- `sync-all-sites`: Syncs all registered sites.
- `open-oauth`: Prints the OAuth authorization URL.

  ```bash
  docker compose run --rm app ./cli list-sites
  ```

- **`sync-site <site_id>`**  
  Syncs a single site by ID.

  ```bash
  docker compose run --rm app ./cli sync-site 123456789
  ```

- **`sync-all-sites`**  
  Syncs all registered sites.

  ```bash
  docker compose run --rm app ./cli sync-all-sites
  ```

- **`open-oauth`**  
  Prints out the OAuth authorization URL so you can copy/paste it into a browser.

  ```bash
  docker compose run --rm app ./cli open-oauth
  ```

- **`force-sync-site <site_id>`**  
  Resets the site’s mapping so that **all** posts will be recreated in Dify upon the next sync.

  ```bash
  docker compose run --rm app ./cli force-sync-site 123456789
  ```

- **`force-sync-doc <site_id> <post_id>`**  
  Removes a single post’s document mapping so it can be recreated.

  ```bash
  docker compose run --rm app ./cli force-sync-doc 123456789 42
  ```

- **`set-post-types <site_id> <post_types_comma_separated>`**  
  Sets which post types will be synced for a site. Defaults to `post` if unset.
  ```bash
  docker compose run --rm app ./cli set-post-types 123456789 post,page
  ```

---

## Running Locally (Without Docker)

If you have Go and Redis installed locally:

1. Start Redis locally:

   ```bash
   redis-server
   ```

2. Install Go dependencies:

   ```bash
   go mod tidy
   ```

3. Build and run the server:

   ```bash
   go build -o server ./cmd/server
   ./server
   ```

   The server runs at `:8080`. Ensure `boc.local` resolves to `127.0.0.1` as described above.

4. Run the CLI locally:
   ```bash
   go build -o cli ./cmd/cli
   ./cli list-sites
   ```
   You can use all the same commands shown above.

---

## Data Storage

- **Redis** stores site configurations and the mapping between WordPress posts and Dify documents.
- By default, data is stored in Docker volumes as defined in `docker compose.yml`.

## API Endpoints

- **GET /**: Returns a "System status: OK" message for health checks.
- **GET /oauth/callback**: Handles the OAuth callback from WordPress.com. On successful authorization, the site is registered and a corresponding Dify dataset is created.

## Updating the Code

To update dependencies or modules:

```bash
go mod tidy
```

If using Docker:

```bash
docker compose build
```

---

## Troubleshooting

- **Accessing `http://boc.local:8080/`**:  
  Ensure your hosts file is correct and autoproxy is disabled.

- **Check Docker logs**:

  ```bash
  docker compose logs -f
  ```

- **OAuth Errors**:
  - Confirm `WPCOM_CLIENT_ID` and `WPCOM_CLIENT_SECRET` match your WordPress OAuth app settings.
  - Ensure the redirect URI matches exactly across `.env` and the WordPress OAuth app.
  - Check for detailed error messages in the browser during authorization.

## Security Note

For production, do not store credentials in plain `.env` files. Use a secure secrets manager (e.g., AWS Secrets Manager, HashiCorp Vault, Docker secrets) to store sensitive information.

---

By following the steps above, you should have a running environment capable of syncing WordPress posts into a Dify dataset.
