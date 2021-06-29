## Invocation

`$ go run main.go input.txt`

`$ cat input.txt | go run main.go`

## Notes

The central recurring theme is a focus on architecture for future extensibility -- putting structures into
place that make it easy for future developers (and encourage them) to Do The Right Thing (tm).

### Directory layout

Having the `eventhandlers` intermediate directory allows for adding common shared helpers/utilities there in the future; each eventhandler getting its own subdirectory in there leads to stronger ACL isolation across different `eventhandlers`, with the option of completely different test patterns and build strategies, if so desired.