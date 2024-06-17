# Todo

- Check `HandleGetGroupMessages`.
- Check ws `comment` function.

- Add ws `likes` function.

- fronted: fix likes and comments on group posts.

- Look out for rare glitch where a group message is only sent to other users, not the creator. Could it be due to how React batches changes to the UI, similar to why `notifyClientOfError` was interfering with sending another message if the two ws messages arrived in close succession?



- Go through audit questions.
- Dockerize.

- edited sqlite.go migrations to run down then up
- followers doubling when following back - if 1 public 1 private - needs more testing? fixes when refreshed
- group posts don't let comments/likes/dislike show in group also can't add comment, can only see stuff on dash
