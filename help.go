package main

import "fmt"

func help() {
	fmt.Println(
		`Act is a tool for AtCoder.

Usage:

	act <command> [arguments]

The commands are:

	act get <content>
		get sample data
	act test <taskID> <sampleID>
		test the task
	act help
		show help
	`)
}
