package main

import (
	"go-error-handling/example"
)

func main() {
	example.BasicErrorExample()

	example.CustomErrorExample(-5)
	example.CustomErrorExample(150)

	example.FormattedErrorExample(-10)
	example.FormattedErrorExample(25)
	example.FormattedErrorExample(150)

	example.WrappingErrorExample("non_existent_file.txt")
	example.WrappingErrorExample("valid_file.txt")

	example.ComplexErrorExample()
	example.CustomErrorExample(999)
}
