CREATE TABLE IF NOT EXISTS  CommentLikes (
    UserId TEXT,
    CommentId TEXT,
    Liked BOOLEAN,
    Disliked BOOLEAN,
    PRIMARY KEY (UserId,  CommentId),
    FOREIGN KEY (CommentId) REFERENCES Comments(Id) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(Id) ON DELETE CASCADE
    );
