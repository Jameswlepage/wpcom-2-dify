# Dify-WP-Sync

This project integrates WordPress.com sites with a [Dify](https://dify.ai) dataset. It syncs posts from a WordPress site into a Dify dataset, allowing you to keep your textual content up-to-date for use with Dify-based applications.

## Prerequisites

1. **Disable Automatic Proxy (if applicable)**  
   Ensure that any automatic proxy (autoproxy) settings are turned off.  
   If you're behind a corporate proxy or using tools that auto-configure a proxy, please disable them or add exceptions for the domains used (like `boc.local`) to avoid connection issues.

2. **Add `boc.local` to Your Hosts File**  
   The OAuth flow uses `http://boc.local:8080` as the callback URL. You need to map `boc.local` to `localhost` on your machine.

   - On **Linux/macOS**, edit `/etc/hosts`:

     ```bash
     sudo nano /etc/hosts
     ```

     Add the following line:

     ```
     127.0.0.1   boc.local
     ```

     Save and close.

   - On **Windows**, edit `C:\Windows\System32\Drivers\etc\hosts` with an Administrator-level editor and add:
     ```
     127.0.0.1   boc.local
     ```

   After editing your hosts file, `boc.local` should resolve to `127.0.0.1`.

3. **Docker and Docker Compose**  
   Make sure you have [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed.

4. **Go (Optional)**  
   If you want to run the CLI tool or server locally without Docker, install [Go 1.20+](https://go.dev/dl/).

   Otherwise, running via Docker Compose will handle building and running the binary for you.

## Configuration

1. **Environment Variables**  
   The app uses a `.env` file to load configuration. There is a `.env.example` file provided that you can use as a template.

   Steps:

   - Copy `.env.example` to `.env`
     ```bash
     cp .env.example .env
     ```
   - Edit the `.env` file and set:
     - `WPCOM_CLIENT_ID` and `WPCOM_CLIENT_SECRET` to your WordPress.com OAuth app credentials.
     - `WPCOM_REDIRECT_URI` should remain `http://boc.local:8080/oauth/callback`.
     - `DIFY_API_KEY` should be set to your Dify API key.
     - `DIFY_BASE_URL` can be updated if you have a custom endpoint, otherwise leave as provided.

   Make sure **NOT** to commit your `.env` file since it contains sensitive credentials. It is listed in `.gitignore` to prevent accidental commits.

2. **Redirect URI Setup in WordPress OAuth App**  
   In your WordPress.com OAuth app settings, ensure the redirect URI matches:
   ```
   http://boc.local:8080/oauth/callback
   ```
   This must be exactly as configured or the OAuth flow will fail.

## Running the Project

### Using Docker Compose

1. Build and start the services (Redis and the Go server):

   ```bash
   docker-compose up --build
   ```

   This command:

   - Pulls and runs a Redis instance.
   - Builds and runs the Go server from the provided `Dockerfile`.

2. After the server starts, it will be available at:

   ```
   http://boc.local:8080
   ```

   - `GET /` will return a simple system status: `System status: OK`.
   - The OAuth callback endpoint is at `GET /oauth/callback`.

3. To authorize a new WordPress site, you'll need to open the OAuth authorization URL in your browser. You can obtain this URL by running the CLI tool or by constructing it manually:

   ```bash
   docker-compose run --rm app ./cli open-oauth
   ```

   Or manually (replace `your_client_id` and `your_redirect_uri` accordingly):

   ```
   https://public-api.wordpress.com/oauth2/authorize?client_id=YOUR_WPCOM_CLIENT_ID&redirect_uri=http%3A%2F%2Fboc.local%3A8080%2Foauth%2Fcallback&response_type=code
   ```

   Open the URL in your browser, authorize the application, and after successful authorization, the site will be registered in the system.

### Using the CLI

The CLI binary (named `cli`) is built inside the Docker image. You can run it in the container:

```bash
docker-compose run --rm app ./cli list-sites
docker-compose run --rm app ./cli sync-site <site_id>
docker-compose run --rm app ./cli sync-all-sites
```

**Commands:**

- `list-sites`: Lists all registered WordPress sites.
- `sync-site <site_id>`: Syncs a single site by ID.
- `sync-all-sites`: Syncs all registered sites.
- `open-oauth`: Prints out the OAuth authorization URL.

### Running Locally Without Docker

If you have Go and Redis installed locally:

1. Start Redis locally:
   ```bash
   redis-server
   ```
2. Build and run the server:

   ```bash
   go build -o server ./cmd/server
   ./server
   ```

   The server will run at `:8080`. Make sure `boc.local` resolves to `127.0.0.1` as described above.

3. Run the CLI locally:
   ```bash
   go build -o cli ./cmd/cli
   ./cli list-sites
   ```

## Data Storage

- **Redis** is used to store site configurations and mapping between WordPress posts and Dify documents.
- The default configuration stores data in Docker volumes as defined in `docker-compose.yml`. For a persistent store, you can mount a volume or configure Redis to store data more reliably.

## Updating the Code

To update dependencies or modules, run:

```bash
go mod tidy
```

Then rebuild the Docker image if you're using Docker:

```bash
docker compose build
```

## Troubleshooting

- If you cannot access `http://boc.local:8080/`, ensure your hosts file is correctly set and that autoproxy is disabled.
- Check Docker logs:
  ```bash
  docker-compose logs -f
  ```
- Ensure your `.env` file is correct and that the `WPCOM_CLIENT_ID`, `WPCOM_CLIENT_SECRET`, and `DIFY_API_KEY` are valid.
- Make sure the OAuth callback URI configured in your WordPress.com application matches `http://boc.local:8080/oauth/callback`.

---

**By following the steps above, you should have a running environment capable of syncing WordPress posts to a Dify dataset.**
