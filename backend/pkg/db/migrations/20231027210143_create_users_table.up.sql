CREATE TABLE IF NOT EXISTS Users (
    Id TEXT NOT NULL PRIMARY KEY,
    Nickname TEXT,
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Password BLOB NOT NULL,
    Avatar TEXT,
    AboutMe TEXT,
    PrivacySetting TEXT NOT NULL,
    git clone --filter=blob:none <repository_url> --path <folder_path>
    CreatedAt DATETIME NOT NULL
);
