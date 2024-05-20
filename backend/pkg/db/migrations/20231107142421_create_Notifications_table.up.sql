CREATE TABLE IF NOT EXISTS Notifications(
   Id TEXT NOT NULL PRIMARY KEY,
   Body TEXT NOT NULL,
   Type TEXT NOT NULL,
   CreatedAt DATETIME NOT NULL,
   ReceiverId TEXT NOT NULL,
   SenderId TEXT NOT NULL,
   Seen BOOLEAN NOT NULL,
   FOREIGN KEY (ReceiverId) REFERENCES Users(Id) ON DELETE CASCADE,
   FOREIGN KEY (SenderId) REFERENCES Users(Id) ON DELETE CASCADE
);