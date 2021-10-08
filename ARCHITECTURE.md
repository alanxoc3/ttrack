# ttrack architecture
`internal/` has most of the code that is called from `main.go`. This is named "internal", because that is a special golang directory name that means the folder content is only accessible within this application and not through the go package management system.

A small section for each module will be explained below.

## internal/cmds/
A public function for each of ttrack's commands is available here. There is also a public "state" struct definition available, which is the input to all functions.

## internal/config/
Currently, there is only a function that helps figure out where the config directory should be located. This function is used by main.go.

## internal/ttdb/
Functions for interfacing with the boltdb database are here.

## internal/ttfile/
Functions for interfacing with the the `.tt` data files are here.

## internal/types/
Standard types used throughout ttrack can be found here. This includes "dates", "seconds" (a subset of golang's duration), and groups (a subset of strings that enforces correct group names).

## main.go
This file has most of the cobra/cli specific logic. The purpose of this main file is to translate from cli arguments to internal/cmds module function calls.
