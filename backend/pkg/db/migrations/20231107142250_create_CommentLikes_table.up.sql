CREATE TABLE IF NOT EXISTS  CommentLikes (
    UserId TEXT,
    CommentId TEXT,
    Liked BOOLEAN,
    Disliked BOOLEAN,
    PRIMARY KEY (UserId,  CommentId),
    FOREIGN KEY (CommentId) REFERENCES Comments(Id),
    FOREIGN KEY (UserID) REFERENCES Users(Id)
    );
