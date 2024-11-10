# vgrep
`vgrep` is a dynamic or layer 2 terminal which sets above the shell, and manages command execution through an 
application that dynamically renders the results of commands in the shell, allowing rapid feedback and tuning of 
commands that are complex or need to be used repetitively.

The intention behind vgrep is to make search based commands more visual, adding a feedback layer to the shell while
protecting against commands which could damage system state. Grep was chosen for the name as this is what I use it for most :)

## Shortcuts
- "ctrl+c" exit the vgrep dynamic terminal.
- "ctrl+q" copy the results from the view directly to the clipboard.
