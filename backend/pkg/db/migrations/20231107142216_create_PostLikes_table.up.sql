CREATE TABLE IF NOT EXISTS PostLikes (
    UserId TEXT,
    PostId TEXT,
    Liked BOOLEAN,
    Disliked BOOLEAN,
    PRIMARY KEY (UserId, PostId),
    FOREIGN KEY (PostId) REFERENCES Posts(Id) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(Id) ON DELETE CASCADE

    );