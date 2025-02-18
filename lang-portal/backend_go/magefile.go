//go:build mage

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

type DB mg.Namespace

// Define your tasks with mg.Namespace if needed
var (
	DB = mg.Namespace("db", "Database operations")
)

// InitDb initializes the database by creating the words.db file if it doesn't exist.
func (DB) Init() error {
	fmt.Println("Initializing database...")

	dbPath := filepath.Join(".", "words.db") // Database file in the project root

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("Creating database file: %s\n", dbPath)
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create database file: %w", err)
		}
		defer file.Close()
		fmt.Println("Database file created successfully.")

		// Basic verification (optional): try to open and close the database file
		// You'll need to import the database/sql and sqlite3 driver for this in the future
		// For now, we'll skip the verification to keep it simple.

	} else if err == nil {
		fmt.Println("Database file already exists.")
	} else {
		return fmt.Errorf("error checking database file existence: %w", err)
	}

	fmt.Println("Database initialization complete.")
	return nil
}

// Migrate runs database migrations.
func (DB) Migrate() error {
	fmt.Println("Running database migrations...")

	dbPath := filepath.Join(".", "words.db")
	migrationsDir := filepath.Join(".", "db", "migrations") // Path to migrations directory

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get already applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationFiles := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles) // Ensure migrations are run in order

	for _, filename := range migrationFiles {
		migrationName := strings.TrimSuffix(filename, ".sql")
		migrationNumberStr := strings.SplitN(migrationName, "_", 2)[0] // Extract number prefix (corrected split)

		// --- Debug logging ---
		fmt.Printf("Debug: Filename: '%s', migrationName: '%s', migrationNumberStr: '%s'\n", filename, migrationName, migrationNumberStr) // Corrected debug output
		// --- End debug logging ---

		migrationNumber, err := strconv.Atoi(migrationNumberStr)
		if err != nil {
			fmt.Printf("Skipping migration file with invalid name: %s\n", filename)
			continue // Skip files with invalid names
		}

		if _, applied := appliedMigrations[migrationNumber]; applied {
			fmt.Printf("Migration %s already applied, skipping.\n", filename)
			continue // Skip already applied migrations
		}

		migrationPath := filepath.Join(migrationsDir, filename)
		sqlBytes, err := ioutil.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}
		sqlStr := string(sqlBytes)

		fmt.Printf("Applying migration: %s\n", filename)
		_, err = db.Exec(sqlStr)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration as applied
		if err := recordMigration(db, migrationNumber); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}
		fmt.Printf("Migration %s applied successfully.\n", filename)
	}

	fmt.Println("Database migrations complete.")
	return nil
}

func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %w", err)
	}

	rows, err := db.Query("SELECT id FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan migration id: %w", err)
		}
		applied[id] = true
	}
	return applied, nil
}

func recordMigration(db *sql.DB, id int) error {
	_, err := db.Exec("INSERT INTO migrations (id) VALUES (?)", id)
	if err != nil {
		return fmt.Errorf("failed to insert migration record: %w", err)
	}
	return nil
}

// Seed seeds the database with initial data.
func (DB) Seed() error {
	fmt.Println("Seeding database with data...")

	dbPath := filepath.Join(".", "words.db")
	seedsDir := filepath.Join(".", "db", "seeds") // Path to seeds directory

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Read seed files
	files, err := ioutil.ReadDir(seedsDir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			seedFilePath := filepath.Join(seedsDir, file.Name())
			fmt.Printf("Seeding from file: %s\n", seedFilePath)

			// Determine table name based on filename (e.g., words.json -> words)
			tableName := strings.TrimSuffix(file.Name(), ".json")

			// Read JSON data from file
			jsonData, err := ioutil.ReadFile(seedFilePath)
			if err != nil {
				return fmt.Errorf("failed to read seed file %s: %w", seedFilePath, err)
			}

			// Unmarshal JSON data into a slice of maps
			var data []map[string]interface{}
			if err := json.Unmarshal(jsonData, &data); err != nil {
				return fmt.Errorf("failed to unmarshal JSON from %s: %w", seedFilePath, err)
			}

			// Prepare and execute SQL insert statements
			for _, item := range data {
				columns := make([]string, 0, len(item))
				values := make([]interface{}, 0, len(item))
				placeholders := make([]string, 0, len(item))

				for column, value := range item {
					columns = append(columns, column)
					// Handle time.Time values
					if column == "created_at" {
						timeStr, ok := value.(string)
						if !ok {
							return fmt.Errorf("created_at value is not a string: %v", value)
						}
						parsedTime, err := time.Parse(time.RFC3339, timeStr)
						if err != nil {
							return fmt.Errorf("failed to parse created_at time: %w", err)
						}
						values = append(values, parsedTime)
					} else {
						values = append(values, value)
					}
					placeholders = append(placeholders, "?")
				}

				// Construct SQL insert statement
				sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
					tableName,
					strings.Join(columns, ", "),
					strings.Join(placeholders, ", "))

				// Execute SQL insert statement
				_, err = db.Exec(sqlStr, values...)
				if err != nil {
					return fmt.Errorf("failed to insert data into %s: %w", tableName, err)
				}
			}

			fmt.Printf("Seeded %d records into table %s\n", len(data), tableName)
		}
	}

	fmt.Println("Database seeding complete.")
	return nil
}

// Install installs project dependencies. (Placeholder for now)
func Install() error {
	fmt.Println("Installing dependencies... (Not yet implemented)")
	// TODO: Implement dependency installation if needed (e.g., using go mod download)
	return nil
}

// Default task to run when you just type `mage`
var Default = DB.Seed

// Run all tests
func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run tests with coverage
func TestCover() error {
	fmt.Println("Running tests with coverage...")
	cmd := exec.Command("go", "test", "./...", "-coverprofile=coverage.out", "-covermode=atomic")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return exec.Command("go", "tool", "cover", "-html=coverage.out").Run()
}

/*
```
**Explanation of `magefile.go`:**

*   **`package main`**:  This is the main package for our Mage build file.
*   **`import (...)`**:  Imports necessary packages:
    *   `fmt`: For formatted printing.
    *   `os`: For operating system functionalities like file system operations.
    *   `path/filepath`: For working with file paths.
    *   `github.com/magefile/mage/mg`: The core Mage package.
*   **`InitDb() error` Function**:
    *   This function is our `initdb` task.
    *   It checks if `words.db` exists in the project root.
    *   If it doesn't exist, it creates the file.
    *   It prints messages to the console indicating the progress.
*   **`Migrate() error` and `Seed() error` Functions**:
    *   These are placeholder functions for the `migrate` and `seed` tasks.
    *   They currently just print "Not yet implemented". We will add the actual logic later.
*   **`Install() error` Function**:
    *   Placeholder for dependency installation.
*   **`var Default = InitDb`**:
    *   This line sets the `InitDb` function as the default task. So, when you run `mage` in the terminal without specifying a task, it will run `InitDb`.

**Step 5: Run the `initdb` task**

Now, in your terminal, still in the `lang-portal/backend_go` directory, run the `initdb` task using Mage:

```bash
mage initdb
```

**Expected Output:**

If this is the first time you run it, you should see output similar to:

```
Initializing database...
Creating database file: words.db
Database file created successfully.
Database initialization complete.
```

If you run it again, you should see:

```
Initializing database...
Database file already exists.
Database initialization complete.
```

You should now have a `words.db` file in your `lang-portal/backend_go` directory.

**Next Steps:**

1.  **Install Gin:**  We'll need to install the Gin framework to start building our API.
2.  **Create `main.go`:** We'll create a `main.go` file to set up the Gin server and define our API routes.
3.  **Database Interaction:** We'll start implementing the database layer to interact with `words.db`.

Let me know if you were able to run `mage initdb` successfully, and we can proceed to the next steps!
*/
