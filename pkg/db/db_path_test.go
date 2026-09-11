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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/OrnilioNeto/Task64/pkg/config"
	"github.com/OrnilioNeto/Task64/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.InitLogger()
	os.Exit(m.Run())
}

func Test_resolveDatabasePath(t *testing.T) {
	mockUserDataDir := func(path string) func() (string, error) {
		return func() (string, error) {
			return path, nil
		}
	}

	mockUserDataDirError := func() (string, error) {
		return "", fmt.Errorf("no home directory")
	}

	tests := []struct {
		name        string
		cfg         DatabasePathConfig
		userDataDir func() (string, error)
		expected    string
		expectError bool
	}{
		{
			name: "memory database",
			cfg: DatabasePathConfig{
				ConfiguredPath: "memory",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "memory",
		},

		{
			name: "absolute path should be used as-is",
			cfg: DatabasePathConfig{
				ConfiguredPath: "/var/lib/task64/task64.db",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/var/lib/task64/task64.db",
		},
		{
			name: "absolute path with different rootpath still used as-is",
			cfg: DatabasePathConfig{
				ConfiguredPath: "/data/mydb.db",
				RootPath:       "/custom/path",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/data/mydb.db",
		},

		{
			name: "relative path with explicit rootpath",
			cfg: DatabasePathConfig{
				ConfiguredPath: "task64.db",
				RootPath:       "/var/lib/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/var/lib/task64/task64.db",
		},
		{
			name: "relative subdirectory path with explicit rootpath",
			cfg: DatabasePathConfig{
				ConfiguredPath: "data/task64.db",
				RootPath:       "/var/lib/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/var/lib/task64/data/task64.db",
		},

		{
			name: "relative path with default rootpath uses user data dir",
			cfg: DatabasePathConfig{
				ConfiguredPath: "task64.db",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/home/user/.local/share/task64/task64.db",
		},

		{
			name: "os.Executable failure falls back to user data dir",
			cfg: DatabasePathConfig{
				ConfiguredPath: "task64.db",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/home/user/.local/share/task64/task64.db",
		},

		{
			name: "falls back to rootpath when userDataDir fails",
			cfg: DatabasePathConfig{
				ConfiguredPath: "task64.db",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDirError,
			expected:    "/opt/task64/task64.db",
		},

		{
			name: "empty configured path with explicit rootpath",
			cfg: DatabasePathConfig{
				ConfiguredPath: "",
				RootPath:       "/var/lib/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/var/lib/task64",
		},
		{
			name: "empty configured path with default rootpath",
			cfg: DatabasePathConfig{
				ConfiguredPath: "",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/home/user/.local/share/task64",
		},
		{
			name: "path with dots normalized",
			cfg: DatabasePathConfig{
				ConfiguredPath: "/var/lib/task64/../task64/./db.db",
				RootPath:       "/opt/task64",
				ExecutablePath: "/opt/task64",
			},
			userDataDir: mockUserDataDir("/home/user/.local/share/task64"),
			expected:    "/var/lib/task64/db.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveDatabasePath(tt.cfg, tt.userDataDir)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_resolveDatabasePath_Integration(t *testing.T) {
	t.Run("with explicitly configured rootpath", func(t *testing.T) {
		cfg := DatabasePathConfig{
			ConfiguredPath: "task64.db",
			RootPath:       "/custom/path",
			ExecutablePath: "/opt/task64",
		}

		result, err := resolveDatabasePath(cfg, resolveUserDataDir)
		require.NoError(t, err)

		expected := filepath.Join("/custom/path", "task64.db")
		assert.Equal(t, expected, result)
	})

	t.Run("with default rootpath uses user data directory", func(t *testing.T) {
		execPath, err := os.Executable()
		require.NoError(t, err)
		execDir := filepath.Dir(execPath)

		cfg := DatabasePathConfig{
			ConfiguredPath: "task64.db",
			RootPath:       execDir,
			ExecutablePath: execDir,
		}

		result, err := resolveDatabasePath(cfg, resolveUserDataDir)
		require.NoError(t, err)

		assert.NotEqual(t, filepath.Join(execDir, "task64.db"), result)
		assert.Contains(t, result, "task64.db")

		switch runtime.GOOS {
		case "windows":
			assert.Contains(t, result, "Task64")
		case "darwin":
			assert.Contains(t, result, "Library")
			assert.Contains(t, result, "Application Support")
		default:
			assert.NotEqual(t,
				filepath.Dir(result),
				execDir,
				"Database should not be in executable directory",
			)
		}
	})

	t.Run("with subdirectory path", func(t *testing.T) {
		cfg := DatabasePathConfig{
			ConfiguredPath: "data/task64.db",
			RootPath:       "/custom/path",
			ExecutablePath: "/opt/task64",
		}

		result, err := resolveDatabasePath(cfg, resolveUserDataDir)
		require.NoError(t, err)

		expected := filepath.Join("/custom/path", "data", "task64.db")
		assert.Equal(t, expected, result)
	})
}

func Test_resolveDatabasePath_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	mockUserDataDir := func(path string) func() (string, error) {
		return func() (string, error) {
			return path, nil
		}
	}

	tests := []struct {
		name        string
		cfg         DatabasePathConfig
		userDataDir func() (string, error)
		expected    string
	}{
		{
			name: "windows absolute path",
			cfg: DatabasePathConfig{
				ConfiguredPath: "C:\\ProgramData\\Task64\\task64.db",
				RootPath:       "C:\\Program Files\\Task64",
				ExecutablePath: "C:\\Program Files\\Task64",
			},
			userDataDir: mockUserDataDir("C:\\Users\\test\\AppData\\Local\\Task64"),
			expected:    "C:\\ProgramData\\Task64\\task64.db",
		},
		{
			name: "windows relative path with explicit rootpath",
			cfg: DatabasePathConfig{
				ConfiguredPath: "task64.db",
				RootPath:       "C:\\ProgramData\\Task64",
				ExecutablePath: "C:\\Program Files\\Task64",
			},
			userDataDir: mockUserDataDir("C:\\Users\\test\\AppData\\Local\\Task64"),
			expected:    "C:\\ProgramData\\Task64\\task64.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveDatabasePath(tt.cfg, tt.userDataDir)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolveUserDataDir(t *testing.T) {
	test := func() string {
		dataDir, err := resolveUserDataDir()
		require.NoError(t, err)
		assert.NotEmpty(t, dataDir)

		_, err = os.Stat(dataDir)
		require.ErrorIs(t, err, os.ErrNotExist, "resolveUserDataDir must not create the directory")

		return dataDir
	}

	// Verify platform-specific paths. Env vars point at fresh t.TempDir()s so the
	// "not created" assertion isn't confused by a directory left over from real usage.
	switch runtime.GOOS {
	case "windows":
		t.Setenv("LOCALAPPDATA", filepath.Join(t.TempDir(), "AppData", "Local"))
		dataDir := test()
		assert.Contains(t, dataDir, "Task64")
	case "darwin":
		t.Setenv("HOME", t.TempDir())
		dataDir := test()
		assert.Contains(t, dataDir, "Library")
		assert.Contains(t, dataDir, "Application Support")
		assert.Contains(t, dataDir, "Task64")
	default:
		t.Run("with XDG_DATA_HOME", func(t *testing.T) {
			xdgDataHome := t.TempDir()
			t.Setenv("XDG_DATA_HOME", xdgDataHome)
			dataDir := test()
			assert.Equal(t, filepath.Join(xdgDataHome, "task64"), dataDir)
		})

		t.Run("without XDG_DATA_HOME", func(t *testing.T) {
			t.Setenv("XDG_DATA_HOME", "")
			home := t.TempDir()
			t.Setenv("HOME", home)
			dataDir := test()
			assert.Equal(t, filepath.Join(home, ".local", "share", "task64"), dataDir)
		})
	}
}

func setUserDataDirForTest(t *testing.T, basePath string) string {
	switch runtime.GOOS {
	case "windows":
		t.Setenv("LOCALAPPDATA", basePath)
		return filepath.Join(basePath, "Task64")
	case "darwin":
		t.Setenv("HOME", basePath)
		return filepath.Join(basePath, "Library", "Application Support", "Task64")
	default:
		t.Setenv("XDG_DATA_HOME", basePath)
		return filepath.Join(basePath, "task64")
	}
}

func setupDatabasePathTest(t *testing.T, configuredPath string) (dataDir string) {
	t.Cleanup(config.ResetForTests)

	dataDir = setUserDataDirForTest(t, filepath.Join(t.TempDir(), "data"))

	execPath, err := os.Executable()
	require.NoError(t, err)

	config.DatabasePath.Set(configuredPath)
	// Matching the executable dir is what makes resolution use the data dir.
	config.ServiceRootpath.Set(filepath.Dir(execPath))

	return dataDir
}

func TestResolvedDatabasePath_DoesNotCreateDataDir(t *testing.T) {
	dataDir := setupDatabasePathTest(t, "task64.db")

	path, err := ResolvedDatabasePath()
	require.NoError(t, err)
	assert.Contains(t, path, "task64.db")

	_, err = os.Stat(dataDir)
	require.ErrorIs(t, err, os.ErrNotExist, "ResolvedDatabasePath must not create the data directory")
}

func TestResolvedDatabasePath_MatchesEnsureDatabasePath(t *testing.T) {
	dataDir := setupDatabasePathTest(t, "task64.db")

	// Order matters: until the data directory exists, resolution reports the rootpath
	// fallback instead - see TestResolvedDatabasePath_FallsBackWhileDataDirMissing.
	ensured, err := ensureDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dataDir, "task64.db"), ensured)

	info, err := os.Stat(filepath.Dir(ensured))
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "ensureDatabasePath must create the directory it returns")

	resolved, err := ResolvedDatabasePath()
	require.NoError(t, err)

	assert.Equal(t, ensured, resolved)
}

// TestResolvedDatabasePath_FallsBackWhileDataDirMissing pins the one case where doctor
// and the server disagree: a data directory nobody created yet.
func TestResolvedDatabasePath_FallsBackWhileDataDirMissing(t *testing.T) {
	setupDatabasePathTest(t, "task64.db")

	resolved, err := ResolvedDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(config.ServiceRootpath.GetString(), "task64.db"), resolved)
}

func TestEnsureDatabasePath_DoesNotCreateConfiguredDirectories(t *testing.T) {
	t.Run("absolute path", func(t *testing.T) {
		missingParent := filepath.Join(t.TempDir(), "deep", "nested")
		setupDatabasePathTest(t, filepath.Join(missingParent, "task64.db"))

		path, err := ensureDatabasePath()
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(missingParent, "task64.db"), path)

		_, err = os.Stat(missingParent)
		require.ErrorIs(t, err, os.ErrNotExist, "ensureDatabasePath must not create a configured path's parents")
	})

	t.Run("rootpath relative path", func(t *testing.T) {
		setupDatabasePathTest(t, filepath.Join("sub", "dir", "task64.db"))
		rootPath := filepath.Join(t.TempDir(), "install")
		config.ServiceRootpath.Set(rootPath)

		path, err := ensureDatabasePath()
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(rootPath, "sub", "dir", "task64.db"), path)

		_, err = os.Stat(rootPath)
		require.ErrorIs(t, err, os.ErrNotExist, "ensureDatabasePath must not create the rootpath")
	})
}

func TestEnsureDatabasePath_CreatesUserDataDir(t *testing.T) {
	dataDir := setupDatabasePathTest(t, "task64.db")

	path, err := ensureDatabasePath()
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(dataDir, "task64.db"), path)

	info, err := os.Stat(dataDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestEnsureDatabasePath_FallsBackWhenDataDirIsNotCreatable(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("directory permissions do not apply here")
	}

	setupDatabasePathTest(t, "task64.db")

	readOnly := t.TempDir()
	setUserDataDirForTest(t, filepath.Join(readOnly, "data"))
	require.NoError(t, os.Chmod(readOnly, 0o500))
	t.Cleanup(func() { _ = os.Chmod(readOnly, 0o700) })

	path, err := ensureDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(config.ServiceRootpath.GetString(), "task64.db"), path)

	resolved, err := ResolvedDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, path, resolved, "doctor must report the path the server falls back to")
}

func TestIsSystemDirectory(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Windows system directories
		{
			name:     "Windows System32",
			path:     "C:\\Windows\\System32\\task64.db",
			expected: runtime.GOOS == "windows",
		},
		{
			name:     "Windows SysWOW64",
			path:     "C:\\Windows\\SysWOW64\\task64.db",
			expected: runtime.GOOS == "windows",
		},
		{
			name:     "Windows root",
			path:     "C:\\Windows\\task64.db",
			expected: runtime.GOOS == "windows",
		},
		{
			name:     "Windows System32 lowercase",
			path:     "c:\\windows\\system32\\task64.db",
			expected: runtime.GOOS == "windows",
		},
		// Unix-like system directories
		{
			name:     "/bin",
			path:     "/bin/task64.db",
			expected: runtime.GOOS != "windows",
		},
		{
			name:     "/sbin",
			path:     "/sbin/task64.db",
			expected: runtime.GOOS != "windows",
		},
		{
			name:     "/usr/bin",
			path:     "/usr/bin/task64.db",
			expected: runtime.GOOS != "windows",
		},
		{
			name:     "/etc",
			path:     "/etc/task64.db",
			expected: runtime.GOOS != "windows",
		},
		// Non-system directories
		{
			name:     "user home directory (Unix)",
			path:     "/home/user/task64.db",
			expected: false,
		},
		{
			name:     "user profile directory (Windows)",
			path:     "C:\\Users\\user\\task64.db",
			expected: false,
		},
		{
			name:     "custom directory",
			path:     "/opt/task64/task64.db",
			expected: false,
		},
		{
			name:     "relative path",
			path:     "./task64.db",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSystemDirectory(tt.path)
			assert.Equal(t, tt.expected, result,
				"Expected isSystemDirectory(%s) to be %v on %s",
				tt.path, tt.expected, runtime.GOOS)
		})
	}
}

func TestIsSystemDirectory_EdgeCases(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Run("false positives - paths containing 'windows' but not system directories", func(t *testing.T) {
			tests := []struct {
				name string
				path string
			}{
				{
					name: "custom app with windows in path",
					path: "C:\\myapp\\windows\\data\\task64.db",
				},
				{
					name: "windows directory on non-C drive",
					path: "D:\\windows\\task64.db",
				},
				{
					name: "user directory named windows",
					path: "C:\\Users\\windows\\task64.db",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					assert.False(t, isSystemDirectory(tt.path))
				})
			}
		})

		t.Run("safe Windows subdirectories", func(t *testing.T) {
			assert.False(t, isSystemDirectory("C:\\Windows\\Temp\\task64.db"))
		})

		t.Run("actual Windows system directories", func(t *testing.T) {
			tests := []struct {
				name string
				path string
			}{
				{
					name: "Windows root",
					path: "C:\\Windows\\task64.db",
				},
				{
					name: "Windows root lowercase",
					path: "c:\\windows\\task64.db",
				},
				{
					name: "System32",
					path: "C:\\Windows\\System32\\task64.db",
				},
				{
					name: "System32 uppercase",
					path: "C:\\WINDOWS\\SYSTEM32\\task64.db",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					assert.True(t, isSystemDirectory(tt.path))
				})
			}
		})
	} else {
		t.Run("false positives - paths containing system dir names", func(t *testing.T) {
			tests := []struct {
				name string
				path string
			}{
				{
					name: "/home/bin not same as /bin",
					path: "/home/bin/task64.db",
				},
				{
					name: "/opt/sbin not same as /sbin",
					path: "/opt/sbin/task64.db",
				},
				{
					name: "/usr/local/bin is safe",
					path: "/usr/local/bin/task64.db",
				},
				{
					name: "/binaries not same as /bin",
					path: "/binaries/task64.db",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					assert.False(t, isSystemDirectory(tt.path))
				})
			}
		})

		t.Run("actual Unix system directories", func(t *testing.T) {
			tests := []struct {
				name string
				path string
			}{
				{
					name: "/bin",
					path: "/bin/task64.db",
				},
				{
					name: "/etc",
					path: "/etc/task64.db",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					assert.True(t, isSystemDirectory(tt.path))
				})
			}
		})
	}
}
