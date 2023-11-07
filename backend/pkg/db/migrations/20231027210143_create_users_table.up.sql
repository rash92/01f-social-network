CREATE TABLE IF NOT EXISTS Users (
    Id BLOB NOT NULL PRIMARY KEY,
    Nickname TEXT NOT NULL UNIQUE,
    FirstName TEXT,
    LastName TEXT,
    Age INTEGER,
    Gender TEXT,
    Email TEXT NOT NULL UNIQUE,
    Password BLOB NOT NULL,
    Profile TEXT,
    AboutMe TEXT,
     privacy_setting TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL
);

