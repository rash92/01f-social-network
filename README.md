# Todo

- Check new handler function to send group messages to avoid "temporal dead zone" where messages aren't seen if sent after the user has gone to the group profile and before they open the chat.
- Check new sw comment function.

- Add ws `likes` function.

- fronted: fix likes and comments on group posts.

- Look out for rare glitch where a group message is only sent to other users, not the creator. Could it be due to how React batches changes to the UI, similar to why the notifyClientOfError function was interfering with sending another message if the two ws messages arrived in close succession?

- Go through audit questions.
- Dockerize.
