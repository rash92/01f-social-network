CREATE TABLE IF NOT EXISTS PrivateMessages (
    Id TEXT NOT NULL PRIMARY KEY,
    SenderId TEXT NOT NULL,
    RecipientId TEXT NOT NULL,
    Message TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    FOREIGN KEY (SenderId) REFERENCES Users(Id),
    FOREIGN KEY (RecipientId) REFERENCES Users(Id)
);