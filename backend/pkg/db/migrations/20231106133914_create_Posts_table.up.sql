CREATE TABLE IF NOT EXISTS Posts (
Id BLOB NOT NULL PRIMARY KEY,
Title TEXT NOT NULL,
Body TEXT NOT NULL ,
CreatorId BLOB NOT NULL,
GroupId  BLOB,
CreatedAt DATETIME NOT NULL,
Image TEXT,
PrivacyLevel TEXT NOT NULl,
FOREIGN KEY (CreatorId) REFERENCES Users(Id),
FOREIGN KEY (GroupId) REFERENCES [Groups](Id)
);