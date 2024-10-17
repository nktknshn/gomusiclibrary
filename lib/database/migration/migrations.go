package migration

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"slices"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
)

var ErrInvalidFolderContent = errors.New("Migrations folder must only contain migration files")

type migrationScripts struct {
	up   string
	down string
}

func (m migrationScripts) String() string {
	return fmt.Sprintf("migration(up: `%s`, down: `%s`)", m.up, m.down)
}

type migration struct {
	id        uint
	name      string
	migration migrationScripts
	filename  string
}

func parseMigration(r io.Reader) (*migrationScripts, error) {
	/*
		-- up
		up script;

		-- down
		down script;
	*/

	s := bufio.NewScanner(r)

	var (
		up, down string
		inDown   = false
	)

	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "-- down") {
			inDown = true
			continue
		}
		if strings.HasPrefix(t, "--") {
			continue
		}

		if inDown {
			down += t + "\n"
			continue
		}
		up += t + "\n"
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("parse migration: %w", err)
	}

	m := &migrationScripts{up, down}

	return m, nil
}

var regexpMigrationFile = regexp.MustCompile(`^[0-9]+_\w+\.sql$`)

func readMigrationsFromFolder(folder fs.FS) ([]migration, error) {

	list, err := fs.ReadDir(folder, ".")

	if err != nil {
		return nil, fmt.Errorf("migration: list dir: %w", err)
	}

	var migrations []migration

	for _, e := range list {
		if e.IsDir() {
			panic(ErrInvalidFolderContent)
		}

		fname := e.Name()

		if !regexpMigrationFile.MatchString(fname) {
			return nil, fmt.Errorf("migration: invalid item name: %s", fname)
		}

		var (
			id   uint
			name string
		)

		_, err := fmt.Sscanf(strings.TrimSuffix(fname, ".sql"), "%03d_%s", &id, &name)

		if err != nil {
			return nil, fmt.Errorf("migration: invalid item name: %s", fname)
		}

		f, err := folder.Open(fname)

		if err != nil {
			return nil, fmt.Errorf("migration: open file: %s", fname)
		}

		m, err := parseMigration(f)

		if err != nil {
			return nil, fmt.Errorf("migration: parse migration: %w", err)
		}

		migrations = append(migrations, migration{id, name, *m, fname})

	}

	return migrations, nil
}

func readSchemaFromFolder(folder fs.FS) (*sqlitemigration.Schema, error) {
	ms, err := readMigrationsFromFolder(folder)

	if err != nil {
		return nil, err
	}

	ups := make([]string, len(ms))

	slices.SortFunc(ms, func(m1, m2 migration) int {
		return int(m1.id) - int(m2.id)
	})

	for idx, m := range ms {
		ups[idx] = m.migration.up
	}

	return &sqlitemigration.Schema{Migrations: ups}, nil
}

func MigrateConn(ctx context.Context, conn *sqlite.Conn, folder fs.FS) error {
	schema, err := readSchemaFromFolder(folder)

	if err != nil {
		return err
	}

	return sqlitemigration.Migrate(ctx, conn, *schema)
}

func MigrateFile(ctx context.Context, dbfile string, folder fs.FS) error {

	conn, err := sqlite.OpenConn(dbfile, sqlite.OpenReadWrite|sqlite.OpenCreate)

	if err != nil {
		return fmt.Errorf("migrate: open database: %w", err)
	}

	defer conn.Close()

	return MigrateConn(ctx, conn, folder)
}
