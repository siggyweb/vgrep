# vgrep
`vgrep` is a dynamic or layer 2 terminal which sets between the shell and user, to run and dislplay commands in real-time. 


The intention behind `vgrep` is to make search based commands more visual and interactive, adding a feedback layer to the shell while
protecting against commands which could damage system state. `grep` was chosen for the name as this is what I use it for most :)

## Features
- Shell commands run after a debounced timer based on user input.
- Command history using up/down arrow keys in the familiar shell fashion.
- Integrated clipboard support to immediately copy results to the clipboard (this will copy all results including cutoff values).
- Writes local log files to `app.log` off user's `$HOME` directory, this is **overwritten** every time `vgrep` is run.


## Shortcuts & Controls
- `ctrl+c` will exit the `vgrep` dynamic terminal.
- `ctrl+q` will copy the results from the view directly to the clipboard.
