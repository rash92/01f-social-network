CREATE TABLE IF NOT EXISTS GroupMessages (
    Id BLOB NOT NULL PRIMARY KEY,
    SenderId BLOB NOT NULL,
    GroupId BLOB NOT NULL,
    Message TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    FOREIGN KEY (SenderId) REFERENCES Users(Id),
    FOREIGN KEY (GroupId) REFERENCES Groups(Id)
);