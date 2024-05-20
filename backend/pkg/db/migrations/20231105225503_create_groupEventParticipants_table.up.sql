   CREATE TABLE IF NOT EXISTS  GroupEventParticipants (
        EventId TEXT NOT NULL,
        UserId TEXT NOT NULL,
        GroupId TEXT NOT NULL,
        PRIMARY KEY (EventId , UserId),
        FOREIGN KEY (EventId) REFERENCES GroupEvents (Id),
        FOREIGN KEY (UserId) REFERENCES Users (Id),
        FOREIGN KEY (GroupId) REFERENCES Groups (Id)
    )