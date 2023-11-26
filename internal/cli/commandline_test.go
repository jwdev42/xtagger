package cli

import (
	"testing"
)

func TestParseSizeStatement(t *testing.T) {
	//test definitions
	tests := make(map[string]uint64)
	tests["0"] = 0
	tests["10"] = 10
	tests["010"] = 10
	tests["666"] = 666
	tests["64K"] = 65536
	tests["1M"] = 1 * 1024 * 1024
	tests["69M"] = 69 * 1024 * 1024
	tests["5G"] = 5 * 1024 * 1024 * 1024
	tests["23T"] = 23 * 1024 * 1024 * 1024 * 1024
	//run tests
	cmd := new(CommandLine)
	for input, expectedOutput := range tests {
		if err := cmd.parseSizeStatement(input); err != nil {
			t.Errorf("Error for input \"%s\": %s", input, err)
		}
		if expectedOutput != cmd.sizeLimit {
			t.Errorf("Error for input \"%s\": Expected size limit is %d, but received size limit is %d", input, expectedOutput, cmd.sizeLimit)
		}
	}
}

func TestParseSizeStatementNegative(t *testing.T) {
	//negative test definitions
	tests := []string{
		"",
		" ",
		" 54",
		"54 ",
		"-1",
		"+1",
		"2+3",
		"0x12",
		"-23G",
		"5 T",
		"92\tM",
		"abc",
		"5k",
		"6m",
		"7g",
		"8t",
		"2.3",
		"2,3",
		"62,78K",
		"77.2M",
		"18446744073709551616", //uint64 overflow
	}
	//run negative tests
	cmd := new(CommandLine)
	for _, input := range tests {
		if err := cmd.parseSizeStatement(input); err == nil {
			t.Errorf("Expected an error for input \"%s\"", input)
		}
		if cmd.sizeLimit != 0 {
			t.Errorf("Test input \"%s\" did modify the size limit despite having an error", input)
		}
	}
}
