 CREATE TABLE IF NOT EXISTS GroupEvents (
        Id TEXT PRIMARY KEY,
        GroupId TEXT NOT NULL,
        Title TEXT NOT NULL,
        Description TEXT,
        CreatorId  TEXT NOT NULL,
        Time DATETIME NOT NULL,
        FOREIGN KEY (CreatorId) REFERENCES Users (Id),
        FOREIGN KEY (GroupId) REFERENCES Groups (Id)
    )   