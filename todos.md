handlefuncs/Handlecomments.go in unfinished state

- reorganize structs/ new functions made in handlecomments
- validatedCommentRequest that we made still not in use.
- continue with testing db functions, figure out permissions for creating directories and files.

- db tests: test things that are easy to test, don't do super IO heavy stuff

- for handlefuncs ones, make fake requests and responseWriters to use in tests, and also mock databases to deal with other end.