   CREATE TABLE IF NOT EXISTS  GroupEventParticipants (
        EventId TEXT NOT NULL,
        UserId TEXT NOT NULL,
        GroupId TEXT NOT NULL,
        PRIMARY KEY (EventId , UserId),
        FOREIGN KEY (EventId) REFERENCES GroupEvents (Id) ON DELETE CASCADE,
        FOREIGN KEY (UserId) REFERENCES Users (Id) ON DELETE CASCADE,
        FOREIGN KEY (GroupId) REFERENCES Groups (Id) ON DELETE CASCADE
    )