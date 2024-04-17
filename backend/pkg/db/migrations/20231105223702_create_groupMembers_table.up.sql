CREATE TABLE IF NOT EXISTS GroupMembers (
        GroupId TEXT NOT NULL,
        UserId TEXT NOT NULL,
        Status TEXT NOT NULL, -- status can be 'invited', 'joined', or 'request'
        PRIMARY KEY (GroupId, UserId),
        FOREIGN KEY (GroupId) REFERENCES Groups (Id) ON DELETE CASCADE,
        FOREIGN KEY (UserId) REFERENCES Users (Id) ON DELETE CASCADE
    );