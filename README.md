The project is a work in progress and mainly aims to organize files and provide recovery/rollback. It doesn't support symlinks and might lead to unexpected results. 

Current Issues:
- [ ] I have realized that I made a pretty significant design error, which is to let only one move happen instead of letting mulitple moves, honestly its so dumb I am now wondering how in the world did I even proceed with that and thats a very serious refactor right there.
- [ ] Executor isn't concurrent yet and the error handling could be better or atleast I feel like that.
- [ ] Rollback feature is yet to be implemented, which essentially is reading the log and executing the operations in reverse.
- [ ] Cleanup Logs are made so that say a move command is done but its cross device leading to EXDER and then we do make a temp file in dest, copyy data to temp and then fsync and then rename the temp and then delete the source and if said delete is failed, the data still exists in the destination, things are safe overall and a cleanUp function can run after execution to well clean things up.
- [ ] Planning funciton use ripgrep and that should be changed to use something better and more appropriate for the situation. (another refactor here ig).
- [ ] Last and final issue is to style the TUI for dryrun and patch up the control keys, try to make it more user friendly at the least.

Current Progress:
- An parser + validator of the rules.yaml exist
- The planning function exists to get all the files and a dry run function to see what will be executed, conflicts and u get a json file after that to see how that works.
- Executor exists, but isn't tested yet. tbh most of the project isn't tested so proceed with caution.

How to run the project:
- clone the project
- go build main.go
- ./main dryrun rules_dest working_dest for dryrun
- ./main run rules_dest working_dest for dryrun [not recommended yet]
