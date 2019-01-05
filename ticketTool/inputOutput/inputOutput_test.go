package inputOutput

/*
	Example to print a message.
 */
func ExampleDefaultInputOutput_Print() {
	io := DefaultInputOutput{}
	io.Print("TestString")
	// Output:
	// TestString
}