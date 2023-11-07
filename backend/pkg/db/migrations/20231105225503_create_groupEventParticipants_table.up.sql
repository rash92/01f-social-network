   CREATE TABLE IF NOT EXISTS  GroupEventParticipants (
        EventId INTEGER NOT NULL,
        UserId INTEGER NOT NULL,
        Choice TEXT NOT NULL, -- choice can be 'Going' or 'Not Going'
        PRIMARY KEY (EventId ,  UserId),
        FOREIGN KEY (EventId) REFERENCES GroupEvents(Id),
        FOREIGN KEY (UserId ) REFERENCES Users (Id)
    )