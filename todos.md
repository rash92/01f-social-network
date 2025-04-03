handlefuncs/Handlecomments.go in unfinished state

redid convertbase64toimage, need to reorganize and get rid of it in dbfuncs and put it in handlefuncs. It no longer actually saves the file but just returns an image file in memory to be saved later. 

- convertbase64toimage from temp is only used in handlecomments and handlewebsockets, might have correct one to replace it with or move from temp
- make image saving happen in all relevant dbfuncs as used by handlewebsockets

- clean up file storage if db connection dies
- reorganize structs/ new functions made in handlecomments
- Purpose is to call 1 dbfunc from handlefunc, and that either worked or it didnt
- validate stuff in a func before sending to db, send stuff to db, do other stuff after seeing if that was successful.
- look at the saveImage in handlefuncs vs stuff in dbfuncs/ temp as both are used but shouldn't be.


- validatedCommentRequest that we made still not in use.

- don't know about approach used to saving images and cleaning up in dbfuncs/comments.go and dbfuncs/posts.go, but seems to work. seek feedback and possibly redo? In any case, old convertbase64toimage now commented out and using new one.

erroring when making posts sometimes, to do with using unusable objects? may be preexisting. May be frontend thing.