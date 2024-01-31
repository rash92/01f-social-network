CREATE TABLE IF NOT EXISTS Follows (
    FollowerId TEXT NOT NULL,
    FollowingId TEXT NOT NULL,
    PRIMARY KEY (FollowerId, FollowingId),
    FOREIGN KEY (FollowerId) REFERENCES Users (Id),
    FOREIGN KEY (FollowingId) REFERENCES Users (Id)
);