// Task64 is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPostgreSQLConnectionString(t *testing.T) {
	t.Run("with schema", func(t *testing.T) {
		connStr := getPostgreSQLConnectionString("localhost:5432", "task64", "secret", "task64", "task64", "disable", "", "", "", "")
		assert.Equal(t, "postgres://task64:secret@localhost:5432/task64?sslmode=disable&sslcert=&sslkey=&sslrootcert=&search_path=%22task64%22%2Cpublic", connStr)
	})
	t.Run("without schema", func(t *testing.T) {
		connStr := getPostgreSQLConnectionString("localhost:5432", "task64", "secret", "task64", "", "disable", "", "", "", "")
		assert.Equal(t, "postgres://task64:secret@localhost:5432/task64?sslmode=disable&sslcert=&sslkey=&sslrootcert=", connStr)
	})
	t.Run("schema needing quoting", func(t *testing.T) {
		connStr := getPostgreSQLConnectionString("localhost:5432", "task64", "secret", "task64", "MySchema", "disable", "", "", "", "")
		assert.Equal(t, "postgres://task64:secret@localhost:5432/task64?sslmode=disable&sslcert=&sslkey=&sslrootcert=&search_path=%22MySchema%22%2Cpublic", connStr)
	})
	t.Run("query exec mode", func(t *testing.T) {
		connStr := getPostgreSQLConnectionString("localhost:5432", "task64", "secret", "task64", "", "disable", "", "", "", "exec")
		assert.Equal(t, "postgres://task64:secret@localhost:5432/task64?sslmode=disable&sslcert=&sslkey=&sslrootcert=&default_query_exec_mode=exec", connStr)
	})
	t.Run("unix socket", func(t *testing.T) {
		connStr := getPostgreSQLConnectionString("/var/run/postgresql", "task64", "secret", "task64", "public", "disable", "", "", "", "")
		assert.Equal(t, "postgres://task64:secret@:5432/task64?sslmode=disable&sslcert=&sslkey=&sslrootcert=&host=/var/run/postgresql&search_path=%22public%22", connStr)
	})
}

func TestSanitizePostgresConnectionError(t *testing.T) {
	err := errors.New(`parse "postgres://task64:secret@invalid host/task64": invalid IP-literal`)

	sanitized := sanitizePostgresConnectionError(err, "task64", "secret")

	assert.NotContains(t, sanitized.Error(), "task64:secret")
	assert.NotContains(t, sanitized.Error(), "secret")
	assert.Contains(t, sanitized.Error(), "postgres://<redacted>@")
	assert.Contains(t, sanitized.Error(), "invalid IP-literal")
}

func TestSanitizePostgresConnectionErrorWithShortPassword(t *testing.T) {
	err := errors.New(`parse "postgres://postgres:x@[object Object]:5432/railway": invalid IP-literal`)

	sanitized := sanitizePostgresConnectionError(err, "postgres", "x")

	assert.NotContains(t, sanitized.Error(), "postgres:x@")
	assert.Contains(t, sanitized.Error(), "postgres://<redacted>@")
	assert.Contains(t, sanitized.Error(), "invalid IP-literal")
}
