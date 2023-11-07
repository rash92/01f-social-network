CREATE TABLE IF NOT EXISTS [Groups] (
        Id INTEGER PRIMARY KEY,
        Title TEXT NOT NULL,
        Description TEXT NOT NULL,
        CreatorId INTEGER NOT NULL,
        CreatedAt DATETIME NOT NULL,
        FOREIGN KEY (CreatorId) REFERENCES Users (Id)
    );