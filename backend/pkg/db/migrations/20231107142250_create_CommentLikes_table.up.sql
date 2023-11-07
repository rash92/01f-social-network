CREATE TABLE IF NOT EXISTS  CommentLikes (
    UserId BLOB,
    CommentId BLOB,
    Liked BOOL,
    Disliked BOOL,
    PRIMARY KEY (UserId,  CommentId),
    FOREIGN KEY (CommentId) REFERENCES Comments(Id),
    FOREIGN KEY (UserID) REFERENCES Users(Id)
    );
