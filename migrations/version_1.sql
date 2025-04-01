CREATE TABLE version
(
   version INTEGER NOT NULL
)
;

CREATE TABLE repositories
(
   id          INTEGER PRIMARY KEY AUTOINCREMENT,
   git_origin  TEXT,
   module_path TEXT,
   created_at  TIMESTAMP NOT NULL,
   updated_at  TIMESTAMP NOT NULL
)
;

CREATE TABLE directories
(
   id            INTEGER PRIMARY KEY AUTOINCREMENT,
   repository_id INTEGER,
   local_path    TEXT      NOT NULL,
   created_at    TIMESTAMP NOT NULL,
   updated_at    TIMESTAMP NOT NULL,
   FOREIGN KEY (repository_id) REFERENCES repos (id),
   UNIQUE (local_path)
)
;

CREATE TABLE imports
(
   id          INTEGER PRIMARY KEY AUTOINCREMENT,
   import_path TEXT      NOT NULL UNIQUE,
   created_at  TIMESTAMP NOT NULL,
   updated_at  TIMESTAMP NOT NULL
)
;

CREATE TABLE files
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   directory_id INTEGER,
   filepath     TEXT      NOT NULL,
   created_at   TIMESTAMP NOT NULL,
   updated_at   TIMESTAMP NOT NULL,
   FOREIGN KEY (directory_id) REFERENCES dirs (id),
   UNIQUE (directory_id, filepath)
)
;

CREATE TABLE file_imports
(
   id         INTEGER PRIMARY KEY AUTOINCREMENT,
   file_id    INTEGER,
   import_id  INTEGER,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL,
   FOREIGN KEY (file_id) REFERENCES files (id),
   FOREIGN KEY (import_id) REFERENCES imports (id),
   UNIQUE (file_id, import_id)
)
;

CREATE INDEX IF NOT EXISTS idx_directories_repository_id ON dirs (repo_id)
;

CREATE INDEX IF NOT EXISTS idx_files_directory_id ON files (directory_id)
;

CREATE INDEX IF NOT EXISTS idx_file_imports_file_id ON file_imports (file_id)
;

CREATE INDEX IF NOT EXISTS idx_file_imports_import_id ON file_imports (import_id)
;

