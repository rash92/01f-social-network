-- CREATE TABLE IF NOT EXISTS Posts (
-- Id TEXT NOT NULL PRIMARY KEY,
-- Title TEXT NOT NULL,
-- Body TEXT NOT NULL ,
-- CreatorId TEXT NOT NULL,
-- GroupId  TEXT,
-- CreatedAt DATETIME NOT NULL,
-- Image TEXT,
-- PrivacyLevel TEXT NOT NULL,
-- FOREIGN KEY (CreatorId) REFERENCES Users(Id),
-- FOREIGN KEY (GroupId) REFERENCES [Groups](Id)
-- );