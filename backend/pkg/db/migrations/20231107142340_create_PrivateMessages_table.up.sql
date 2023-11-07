CREATE TABLE IF NOT EXISTS PrivateMessages (
    Id BLOB NOT NULL PRIMARY KEY,
    SenderId BLOB NOT NULL,
    RecipientId BLOB NOT NULL,
    Message TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    FOREIGN KEY (SenderId) REFERENCES Users(Id),
    FOREIGN KEY (RecipientId) REFERENCES Users(Id)
);