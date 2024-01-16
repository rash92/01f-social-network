CREATE TABLE IF NOT EXISTS Users (
    Id BLOB NOT NULL PRIMARY KEY,
    Nickname TEXT NOT NULL UNIQUE,
    FirstName TEXT,
    LastName TEXT,
    Email TEXT NOT NULL UNIQUE,
    Password BLOB NOT NULL,
    Profile TEXT,
    AboutMe TEXT,
    Privacy_setting TEXT NOT NULL,
    DOB  TEXT  NOT NULL,
    CreatedAt DATETIME NOT NULL
);

