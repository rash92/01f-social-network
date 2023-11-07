CREATE TABLE IF NOT EXISTS PostLikes (
    UserId BLOB,
    PostId BLOB,
    Liked BOOL,
    Disliked BOOL,
    PRIMARY KEY (UserId, PostId),
    FOREIGN KEY (PostId) REFERENCES Posts(Id),
    FOREIGN KEY (UserID) REFERENCES Users(Id)

    );