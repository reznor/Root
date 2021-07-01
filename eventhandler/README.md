# Directory Layout

Having the [eventhandlers/](eventhandlers/) intermediate directory provides a placeholder for adding common shared helpers/utilities at that level in the future.

Each handler getting its own dedicated subdirectory in [eventhandlers/](eventhandlers/) leads to stronger ACL isolation that allows multiple teams
to own and develop their custom handlers in a central location (which in turn helps with discoverability, auditing, and framework-level
refactoring) while still retaining the option of completely different test patterns, build strategies, etc.

# Adding Support For New `EventType`s

1. Add a new subdirectory in [eventhandlers/](eventhandlers/) following the lead of the existing packages in there.

2. Remember to call `eventhandler.GlobalRegistry().RegisterEventHandler()` in your new package's `init`() function.

3. Remember to actually load your new package into the system by importing your package anonymously in [eventprocessor.go](../eventprocessor/eventprocessor.go)].