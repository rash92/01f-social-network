CREATE TABLE IF NOT EXISTS Follows (
    FollowerId INTEGER NOT NULL,
    FollowingId INTEGER NOT NULL,
    Status  TEXT not NULL,
    PRIMARY KEY (FollowerId, FollowingId),
    FOREIGN KEY (FollowerId) REFERENCES Users(Id),
    FOREIGN KEY (FollowingId) REFERENCES Users(Id)

);