package main

func help() string {
	message := `Act is a tool for AtCoder.

Usage:

	act <command> [arguments]

The commands are:

	act get <content>
		get sample data
	act test <taskID> <sampleID>
		test the task
	act help
		show help
	`
	return message
}
