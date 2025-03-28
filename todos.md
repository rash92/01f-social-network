handlefuncs/Handlecomments.go in unfinished state

redid convertbase64toimage, need to reorganize and get rid of it in dbfuncs and put it in handlefuncs. It no longer actually saves the file but just returns an image file in memory to be saved later.

- maybe save file in same function that adds it to database
- clean up file storage if db connection dies
- reorganize structs/ new functions made in handlecomments
- Purpose is to call 1 dbfunc from handlefunc, and that either worked or it didnt
- validate stuff in a func before sending to db, send stuff to db, do other stuff after seeing if that was successful.
- look at the saveImage in handlefuncs vs stuff in dbfuncs/ temp as both are used but shouldn't be.
- probably conversions in handlefuncs but actual saving in dbfuncs