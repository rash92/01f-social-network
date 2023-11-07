CREATE TABLE  IF NOT EXISTS PostChosenFollowers (
    PostId INTEGER NOT NULL,
    FollowerId INTEGER NOT NULL,
    PRIMARY KEY (PostId, FollowerId)
    FOREIGN KEY (postId) REFERENCES Posts(Id),
    FOREIGN KEY (followerId) REFERENCES Users(Id)
);
