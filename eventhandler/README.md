## Adding Support For New `EventType`s

1. Add a new subdirectory in `./eventhandlers` following the lead of the existing packages in there.
2. Remember to call `eventhandler.GlobalRegistry().RegisterEventHandler()` in your new package's `init`() function.
3. Remember to actually load your new package into the system by importing your package anonymously in `../eventprocessor/eventprocessor.go`.