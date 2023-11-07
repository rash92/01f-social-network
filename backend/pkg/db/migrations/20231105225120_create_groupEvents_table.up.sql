 CREATE TABLE IF NOT EXISTS GroupEvents (
        Id INTEGER PRIMARY KEY,
        GroupId INTEGER NOT NULL,
        Title TEXT NOT NULL,
        Description TEXT,
        CreatorId  BLOB NOT NULL,
        Time DATETIME NOT NULL,
        FOREIGN KEY(CreatorId) REFERENCES Users(Id),
        FOREIGN KEY ( GroupId) REFERENCES Groups (Id)
    )   