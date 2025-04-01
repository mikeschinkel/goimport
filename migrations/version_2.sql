-- ASSUMES `imports` table from version 1 exists and `imports_v1` does not

CREATE TABLE IF NOT EXISTS version
(
   version INTEGER NOT NULL
)
;

CREATE TABLE repos
(
   id          INTEGER PRIMARY KEY AUTOINCREMENT,
   origin_url  TEXT,
   module_path TEXT,
   created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

ALTER TABLE repositories RENAME TO repositories_v1;

INSERT INTO repos(
   id,
   origin_url,
   module_path
) SELECT
   id,
   git_origin,
   module_path
FROM
   repositories_v1;


CREATE UNIQUE INDEX idx_repos_git_origin_url ON repos(origin_url)
;

CREATE TABLE dirs
(
   id            INTEGER PRIMARY KEY AUTOINCREMENT,
   repo_id INTEGER,
   local_path    TEXT      NOT NULL,
   created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (repo_id) REFERENCES repos (id),
   UNIQUE (local_path)
)
;

ALTER TABLE directories RENAME TO directories_v1;

INSERT INTO dirs(
   id,
   repo_id,
   local_path
) SELECT
   id,
   repository_id,
   local_path
FROM
   directories_v1
;


-- Rename v1 imports table first to avoid conflicts
ALTER TABLE imports RENAME TO imports_v1
;

-- Create new imports table
CREATE TABLE imports
(
   id          INTEGER PRIMARY KEY AUTOINCREMENT,
   import_path TEXT      NOT NULL UNIQUE,
   created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at  TIMESTAMP NOT NULL
)
;

-- Create new files table
CREATE TABLE files
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   dir_id INTEGER,
   filepath     TEXT      NOT NULL,
   created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at   TIMESTAMP NOT NULL,
   FOREIGN KEY (dir_id) REFERENCES dirs (id),
   UNIQUE (dir_id, filepath)
)
;

-- Create new file_imports table
CREATE TABLE file_imports
(
   id         INTEGER PRIMARY KEY AUTOINCREMENT,
   file_id    INTEGER,
   import_id  INTEGER,
   created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL,
   FOREIGN KEY (file_id) REFERENCES files (id),
   FOREIGN KEY (import_id) REFERENCES imports (id),
   UNIQUE (file_id, import_id)
)
;

-- Create indexes for new tables
CREATE INDEX idx_files_directory_id ON files (dir_id)
;

CREATE INDEX idx_file_imports_file_id ON file_imports (file_id)
;

CREATE INDEX idx_file_imports_import_id ON file_imports (import_id)
;

-- Migrate data from imports_v1 table to new tables
INSERT INTO imports (import_path, created_at, updated_at)
SELECT DISTINCT
   import_path,
   CURRENT_TIMESTAMP,
   CURRENT_TIMESTAMP
FROM
   imports_v1
;

INSERT INTO files (dir_id, filepath, created_at, updated_at)
SELECT DISTINCT
   directory_id,
   filepath,
   CURRENT_TIMESTAMP,
   CURRENT_TIMESTAMP
FROM
   imports_v1
;

INSERT INTO file_imports (file_id, import_id, created_at, updated_at)
SELECT
   f.id,
   i.id,
   CURRENT_TIMESTAMP,
   CURRENT_TIMESTAMP
FROM
   imports_v1 v1
      JOIN files f
         ON f.dir_id = v1.directory_id AND f.filepath = v1.filepath
      JOIN imports i
         ON i.import_path = v1.import_path
;

DELETE FROM version WHERE 1=1
;

INSERT INTO version (version) VALUES (2)
;

DROP VIEW IF EXISTS repos_named;
CREATE VIEW repos_named as
SELECT
   CASE WHEN IFNULL(origin_url,'')='' THEN module_path ELSE origin_url END AS name,
   *
FROM
   repos
;

DROP VIEW IF EXISTS dir_imports;
CREATE VIEW dir_imports as
SELECT
   d.local_path,
   d.repo_id,
   f.dir_id,
   f.filepath,
   fi.file_id,
   fi.import_id,
   i.import_path
FROM dirs d
        JOIN files f ON d.id = f.dir_id
        JOIN file_imports fi ON f.id = fi.file_id
        JOIN imports i ON fi.import_id = i.id
;
DROP VIEW IF EXISTS repo_imports;
CREATE VIEW repo_imports as
SELECT
   r.id,
   CASE WHEN IFNULL(origin_url,'')='' THEN module_path ELSE origin_url END AS repo_path,
   r.name AS repo_name,
   r.origin_url,
   r.module_path,
   dr.*
FROM repos_named r
  JOIN dir_imports dr ON r.id=dr.repo_id
;

SELECT * FROM repo_imports;