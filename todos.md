handlefuncs/Handlecomments.go in unfinished state

- reorganize structs/ new functions made in handlecomments
- validatedCommentRequest that we made still not in use.
- continue with testing db functions, figure out permissions for creating directories and files.

- db tests: test things that are easy to test, don't do super IO heavy stuff

- for handlefuncs ones, make fake requests and responseWriters to use in tests, and also mock databases to deal with other end.

- move repository struct in handlecomments to structs for use for other handlefuncs. Add tests in handlefuncs and dbfuncs that don't touch the db a ton.

- update repository and default repository in handlefuncs/ structs.go for any dbfuncs functions used in the handlefuncs package so custom ones can later be used in testing.

- test/fix docker compose/ dockerfile stuff

- cleanup all the old unused commented functions

- rewrite websockets?