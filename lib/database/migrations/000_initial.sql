-- up

CREATE TABLE file (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path text NOT NULL UNIQUE,
    
    size INTEGER,
    hash TEXT,

    ctime TEXT,
    mtime TEXT,

    created_at TEXT NOT NULL,
    updated_at TEXT,
    deleted_at TEXT
) STRICT;

CREATE INDEX file_path_idx ON file(path);

CREATE TABLE file_audio_tag (
    id INTEGER PRIMARY KEY,
    file_id INTEGER NOT NULL,

    name TEXT NOT NULL,
    value TEXT,
    
    created_at TEXT NOT NULL,
    updated_at TEXT,
    deleted_at TEXT,

    FOREIGN KEY(file_id) REFERENCES file(id)
) STRICT;

-- down

DROP TABLE file;
DROP INDEX file_path_idx;

DROP TABLE file_audio_tag;