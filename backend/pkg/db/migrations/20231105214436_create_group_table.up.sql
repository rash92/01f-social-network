CREATE TABLE IF NOT EXISTS [Groups] (
        Id TEXT PRIMARY KEY,
        Title TEXT NOT NULL,
        Description TEXT NOT NULL,
        CreatorId TEXT NOT NULL,
        CreatedAt DATETIME NOT NULL,
        FOREIGN KEY (CreatorId) REFERENCES Users (Id) ON DELETE CASCADE
    );