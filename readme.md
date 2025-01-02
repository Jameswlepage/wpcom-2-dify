# Dify-WP-Sync

This project integrates WordPress.com sites with a [Dify](https://dify.ai) dataset. It syncs posts (or any other chosen post types) from a WordPress site into a Dify dataset, keeping your textual content up-to-date for use with Dify-based applications.

---

## Prerequisites

1. **Disable Automatic Proxy (if applicable)**  
   Ensure any automatic proxy (autoproxy) settings are turned off or that `boc.local` is excluded from proxy rules.

2. **Add `boc.local` to Your Hosts File**  
   The OAuth flow uses `http://boc.local:8080` as the callback URL. You need to map `boc.local` to `localhost` on your machine.

   - **Linux/macOS**: Edit `/etc/hosts`:

     ```bash
     sudo nano /etc/hosts
     ```

     Add:

     ```
     127.0.0.1   boc.local
     ```

     Save and close.

   - **Windows**: Edit `C:\Windows\System32\Drivers\etc\hosts` with Administrator privileges and add:
     ```
     127.0.0.1   boc.local
     ```
     After that, `boc.local` should resolve to `127.0.0.1`.

3. **Docker & Docker Compose (v2)**  
   Make sure you have [Docker](https://www.docker.com/) and Docker Compose v2 installed.

   > For Docker Compose v2, the command is `docker compose ...` rather than `docker-compose ...`.

4. **Go (Optional)**  
   If you want to run the CLI tool or server locally outside of Docker, install [Go 1.20+](https://go.dev/dl/).  
   Otherwise, running via Docker Compose will handle building and running the binaries for you.

---

## Configuration

1. **Environment Variables**  
   The app uses a `.env` file to load configuration. A `.env.example` file is provided:

   ```bash
   cp .env.example .env
   ```

   Then edit `.env` to set:

   - `WPCOM_CLIENT_ID` and `WPCOM_CLIENT_SECRET`: your WordPress.com OAuth credentials.
   - `WPCOM_REDIRECT_URI`: should remain `http://boc.local:8080/oauth/callback`.
   - `DIFY_API_KEY`: your Dify API key.
   - `DIFY_BASE_URL`: the Dify endpoint (defaults to `https://api.dify.ai/v1`).

   **Never commit** your `.env` file since it contains sensitive credentials (it's in `.gitignore`).

2. **Redirect URI Setup in WordPress OAuth App**  
   In your WordPress.com OAuth app settings, ensure the callback URI is:
   ```
   http://boc.local:8080/oauth/callback
   ```
   It must match exactly or the OAuth flow will fail.

---

## Running the Project

### Using Docker Compose (v2)

1. **Build and start** the services:

   ```bash
   docker compose up --build
   ```

   This will:

   - Pull and run a Redis instance.
   - Build and run the Go server from the provided `Dockerfile`.
   - Mount the project files for hot reloading via Air.

2. After startup, the server listens on:

   ```
   http://boc.local:8080
   ```

   - `GET /` returns `System status: OK`.
   - The OAuth callback endpoint is `GET /oauth/callback`.

3. **Authorize a new WordPress site**  
   Open the OAuth authorization URL in your browser. You can obtain the URL using the CLI:

   ```bash
   docker compose run --rm app ./cli open-oauth
   ```

   Or construct it manually (replace `YOUR_WPCOM_CLIENT_ID` and `YOUR_REDIRECT_URI`):

   ```
   https://public-api.wordpress.com/oauth2/authorize?client_id=YOUR_WPCOM_CLIENT_ID
     &redirect_uri=YOUR_REDIRECT_URI
     &response_type=code
   ```

   After authorizing, your site will be registered in the system.

---

## CLI Usage

The CLI (`cli`) is built inside the Docker image. You can run it in the container with:

```bash
docker compose run --rm app ./cli <command> [args...]
```

### Commands

- **`list-sites`**  
  Lists all registered WordPress sites.

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

1. **Start Redis**:
   ```bash
   redis-server
   ```
2. **Build and run** the server:

   ```bash
   go build -o server ./cmd/server
   ./server
   ```

   The server will run at `:8080`.  
   Ensure `boc.local` → `127.0.0.1` resolution as described earlier.

3. **Run the CLI locally**:
   ```bash
   go build -o cli ./cmd/cli
   ./cli list-sites
   ```
   You can use all the same commands shown above.

---

## Data Storage

- **Redis** is used to store site configurations and the mapping of WordPress posts to Dify documents.
- Default Docker setup stores data in a volume defined in `docker-compose.yml`.  
  For production or long-term storage, consider configuring Redis persistence or an external volume.

---

## Updating the Code

To update dependencies or modules, run:

```bash
go mod tidy
```

Then rebuild the Docker image if using Docker:

```bash
docker compose build
```

---

## Troubleshooting

- If you can’t reach `http://boc.local:8080/`, check your `/etc/hosts` or `\Windows\System32\Drivers\etc\hosts` file and confirm you have `127.0.0.1 boc.local`.
- Confirm any autoproxy settings are disabled or have exceptions for `boc.local`.
- Inspect logs:
  ```bash
  docker compose logs -f
  ```
- Verify your `.env` file is correct (proper `WPCOM_CLIENT_ID`, `WPCOM_CLIENT_SECRET`, `DIFY_API_KEY`).
- Ensure your WordPress.com app’s callback URI matches `http://boc.local:8080/oauth/callback`.

---

**By following the above steps, you should have a working environment capable of syncing WordPress posts—and any specified post types—to a Dify dataset.**

Enjoy building with Dify!
