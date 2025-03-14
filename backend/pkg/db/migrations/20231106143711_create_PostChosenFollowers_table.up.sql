CREATE TABLE  IF NOT EXISTS PostChosenFollowers (
    PostId TEXT NOT NULL,
    FollowerId TEXT NOT NULL,
    PRIMARY KEY (PostId, FollowerId)
    FOREIGN KEY (postId) REFERENCES Posts(Id) ON DELETE CASCADE,
    FOREIGN KEY (followerId) REFERENCES Users(Id) ON DELETE CASCADE
);
