CREATE TABLE IF NOT EXISTS GroupMessages (
    Id TEXT NOT NULL PRIMARY KEY,
    SenderId TEXT NOT NULL,
    GroupId TEXT NOT NULL,
    Message TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL,
    FOREIGN KEY (SenderId) REFERENCES Users(Id) ON DELETE CASCADE,
    FOREIGN KEY (GroupId) REFERENCES Groups(Id) ON DELETE CASCADE
);