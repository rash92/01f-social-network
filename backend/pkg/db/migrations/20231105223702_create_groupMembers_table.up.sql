CREATE TABLE IF NOT EXISTS GroupMembers (
        GroupId INTEGER NOT NULL,
        UserId INTEGER NOT NULL,
        Status TEXT NOT NULL, -- status can be 'invited', 'joined', or 'request'
        PRIMARY KEY ( GroupId ,  UserId ),
        FOREIGN KEY ( GroupId) REFERENCES Groups ( Id),
        FOREIGN KEY ( UserId) REFERENCES Users (Id)
    );