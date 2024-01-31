CREATE TABLE IF NOT EXISTS  CommentLikes (
    UserId TEXT,
    CommentId TEXT,
    Liked BOOL,
    Disliked BOOL,
    PRIMARY KEY (UserId,  CommentId),
    FOREIGN KEY (CommentId) REFERENCES Comments(Id),
    FOREIGN KEY (UserID) REFERENCES Users(Id)
    );
