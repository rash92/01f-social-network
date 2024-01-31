CREATE TABLE IF NOT EXISTS Sessions (
  Id TEXT NOT NULL PRIMARY KEY,
  Expires DATETIME NOT NULL,
  UserId TEXT,
  UNIQUE (userId),
  FOREIGN KEY (userId) REFERENCES Users (Id)
);