CREATE TABLE IF NOT EXISTS Follows (
    FollowerId TEXT NOT NULL,
    FollowingId TEXT NOT NULL,
    Status TEXT NOT NUll,
    PRIMARY KEY (FollowerId, FollowingId),
    FOREIGN KEY (FollowerId) REFERENCES Users (Id),
    FOREIGN KEY (FollowingId) REFERENCES Users (Id)
);