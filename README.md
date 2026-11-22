# bidness
for bro

## Tooling

### Project (not provisional)

- Build tool -> [make](https://www.gnu.org/software/make/)

### Frontend (provisional)

- **Programming language:** [TypeScript](https://www.typescriptlang.org/)
- **JavaScript runtime:** [Node.js](https://nodejs.org/en)
- **Package manager:** [pnpm](https://pnpm.io/)
- **UI library:** [React](https://react.dev/)
- **React optimisation:** [React Compiler](https://react.dev/learn/react-compiler)
- **Build tool:** [Vite](https://vite.dev/)
- **Routing:** [React Router](https://reactrouter.com/)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **UI components:** [TBC]()
- **API client and data fetching:** [Axios](https://axios-http.com/)
- **Testing:** [Vitest](https://vitest.dev/) + [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- **Linting and formatting:** [Oxc](https://oxc.rs/) (Oxlint + Oxfmt)

### Backend (provisional)

- Programming language -> [go](https://go.dev/)
- Testing -> [testing](https://pkg.go.dev/testing@go1.27.3)
- HTTP server -> [net/http](https://pkg.go.dev/net/http)
- Router -> [chi](https://github.com/go-chi/chi)
- Database -> [Embedded Postgres](https://github.com/fergusstrange/embedded-postgres) (with the [pgx](https://github.com/jackc/pgx) driver)
- Migrations -> [goose](https://github.com/pressly/goose)
- SQL tooling -> [sqlc](https://sqlc.dev/?utm_source=chatgpt.com)
- Logging -> [slog](https://pkg.go.dev/log/slog)
- Profiling -> [pprof](https://pkg.go.dev/net/http/pprof)

## Monorepo

This is currently behaving as a monorepo as both apps are in their own folders under one git history.

To manage some common coordination from the root folder we can use `make`.

I've got some commands already in there -> run the api, build the api, run the api binary, run a vite dev server.

Just run `make <command name>` to try it out.

The two apps are currently not talking to each other but they do run fine independently.

The Go backend has its own Makefile in the `api/` directory. This is for running backend specific commands.

### Working on the backend

List of commands from the `api/` directory:

- `make help`	Lists available Make targets
- `make build`	Compiles the Go application into the api binary
- `make run`	Builds the binary and runs the application
- `make test`	Runs all Go tests
- `make test-verbose`	Runs tests with verbose output
- `make test-coverage`	Runs tests and reports coverage
- `make fmt`	Formats Go code
- `make vet`	Runs Go's static analysis checks
- `make lint`	Runs formatting and go vet
- `make deps`	Downloads Go module dependencies
- `make tidy`	Cleans up go.mod and go.sum dependencies
- `make clean` Removes the compiled application binary

## Embedded Postgres

Postgres is the database we're using, run via [embedded-postgres](https://github.com/fergusstrange/embedded-postgres). This means the Go API starts and stops its own Postgres server, so there's nothing to install or run separately.

### Embedded Postgres vs. a regular Postgres server

A regular Postgres setup works like this:

- Install or run a database server.
- Start the server as a separate process or container.
- Configure a host, port, database name, and credentials.
- Connect to it from your application.

Server database (Postgres)
`Frontend ──> HTTP ──> Go API ──> network connection ──> PostgreSQL server`

With embedded Postgres the Go API does those steps itself on startup:

Embedded Postgres
`Frontend ──> HTTP ──> Go API ──> starts ──> PostgreSQL server ──> data/database/postgres (local directory)`

On first run the library downloads the Postgres binaries into a user cache directory, then starts a real Postgres server as a child process on `localhost:5432`. The Go API connects to it like any other Postgres database (using the `pgx` driver). When the API stops, the Postgres server is shut down too.

### Where the database lives

The backend uses a local directory to keep source data separate from the database it creates:

api/
└── data/
    ├── seed/
    │   └── authors.json
    └── database/
        └── postgres/

`data/seed/authors.json` is the original dataset imported by the application.

`data/database/postgres/` is the Postgres data directory the embedded server creates and updates.

The database directory should generally be treated as generated local data rather than source code and is therefore ignored.

### How database initialization works

When setting up the database we follow this sequence of events.

1. Prepare the database directory. Create the directory if it doesn't exist.
2. Start embedded Postgres. Either reuse the existing data directory or initialise a new one if it doesn't exist, then connect to it.
3. Run migrations. Create or update tables and indexes to match the application's expected schema.
4. Seed initial data. If the database needs its initial dataset, import records from the JSON file.

### What happens on subsequent runs?

The seed logic checks whether the table already contains records and skips the initial import if it has already been populated.

The database can be reset by passing `shouldClear = true` in the setup.

Be careful when resetting: it removes the database directory, deleting all its contents.

### Using psql

`psql` is the Postgres CLI for inspecting the database interactively. It is separate from the Go driver used by the application.

The CLI is useful for when you want to inspect tables, check whether a seed import worked, or experiment with SQL without changing the Go code.

#### Check whether it is installed

`psql --version`

If the command isn't found, the CLI isn't installed or isn't available on your PATH.

See the official PostgreSQL [downloads](https://www.postgresql.org/download/) page for instructions (on macOS, `brew install libpq` gives you just the client tools).

#### Open the database

The embedded server only runs while our API is running, so start the API first. Then run:

`psql postgresql://postgres:postgres@localhost:5432/bidness`

#### Run a few queries

List the available tables:
```sql
\dt
```

Count the imported authors:
```sql
SELECT COUNT(*) FROM authors;
```

Inspect a few records:
```sql
SELECT key, name, revision
FROM authors
LIMIT 10;
```

Search for an author:
```sql
SELECT key, name
FROM authors
WHERE name ILIKE '%Pesgawan%';
```

Exit the CLI:

`\q`
