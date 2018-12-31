package main

import (
	"testing"
)

func Test_LoadingFile_With_False_Path(t *testing.T) {
	//_, available := loadFile("C:\\Users\\User\\Documents\\DHBW Studium\\DHBW\\Meine Kurse\\Semester 5\\Go")
	//assert.True(t, available != true, "Read File is going wrong")
	//Output:
	//<Error>[Read File is going wrong]
	//false
}

func Test_ThatFilePathExist(t *testing.T) {
	//_, valide := loadEmail("alsdfjölaksfdj")
	//assert.True(t, valide == false, "File doesnt exist")
}

func Test_EntryNumberOfRandomMails(t *testing.T) {

	/*
		if you entry not a Number
		print: "Entry is no Number!"
		otherwise, you get your number back
	*/
	entryNumberOfRandomMails()
}
