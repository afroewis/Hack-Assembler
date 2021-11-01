package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var variables = make(map[string]int)
var lines = make([]string, 0)
var output = make([]string, 0)
var labels = make(map[string]int)
var predefinedSymbols = make(map[string]int)
var variablesMemoryStart = 16
func main() {
	predefinedSymbols = map[string]int{
		"SO": 0,
		"LCL": 1,
		"ARG": 2,
		"THIS": 3,
		"THAT": 4,
		"R0": 0x00,
		"R1": 0x01,
		"R2": 0x02,
		"R3": 0x03,
		"R4": 0x04,
		"R5": 0x05,
		"R6": 0x06,
		"R7": 0x07,
		"R8": 0x08,
		"R9": 0x09,
		"R10": 0x0a,
		"R11": 0x0b,
		"R12": 0x0c,
		"R13": 0x0d,
		"R14": 0x0e,
		"R15": 0x0f,
		"SCREEN": 0x4000,
		"KBD": 0x6000,
	}

	file, err := os.Open("assembly/pong/Pong.asm")
	//file, err := os.Open("assembly/add/Add.asm")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line = scanner.Text()

		var commentIndex = strings.Index(line, "//")
		if commentIndex != -1 {
			line = line[:commentIndex]
		}

		line = strings.Trim(line, " ")

		if len(line) == 0 {
			continue
		}

		lines = append(lines, line)
	}

	var instructionNumber = 0
	for _, line := range lines {
		if line[0] == '(' {
			var end = strings.IndexRune(line, ')')
			var symbol = line[1:end]
			labels[symbol] = instructionNumber
			continue
		}
		instructionNumber++
	}

	for _, line := range lines {
		if line[0] == '(' {
			continue
		}

		if line[0] == '@' {
			var value = line[1:]

			if num, err := strconv.Atoi(value); err == nil {
				str := fmt.Sprintf("%016b", num)
				output = append(output, str)
			} else {
				var addressOrValue int

				if val, found := predefinedSymbols[value]; found {
					addressOrValue = val
				} else if variables[value] >= 16 {
					addressOrValue = variables[value]
				} else if labels[value] > 0 {
					addressOrValue = labels[value]
				} else {
					addressOrValue = variablesMemoryStart + len(variables);
					variables[value] = addressOrValue
				}

				str := fmt.Sprintf("%016b", addressOrValue)
				output = append(output, str)
			}

			continue
		}

		var res = parseCInstruction(line)
		var str = fmt.Sprintf("%016b", res)
		output = append(output, str)
	}

	f, err := os.Create("bin/Pong.hack")
	defer f.Close()

	for _, str := range output{
		f.WriteString(str + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseCInstruction(line string) int {
	var equalIndex = strings.IndexRune(line, '=')
	var semicolonIndex = strings.IndexRune(line, ';')

	out := 0b1110_0000_0000_0000

	var isStore = equalIndex != -1
	var isJump = semicolonIndex != -1

	var computation string
	if isStore {
		computation = line[equalIndex+1:]
	} else {
		computation = line[:semicolonIndex]
	}

	switch computation {
	case "0":
		out |= 0b0000_1010_1000_0000
		break;
	case "1":
		out |= 0b0000_1111_1100_0000
		break;
	case "-1":
		out |= 0b0000_1110_1000_0000
		break;
	case "D":
		out |= 0b0000_0011_0000_0000
		break;
	case "A":
		out |= 0b0000_1100_0000_0000
		break;
	case "!D":
		out |= 0b0000_0011_0100_0000
		break;
	case "!A":
		out |= 0b0000_1100_0100_0000
		break;
	case "-D":
		out |= 0b0000_0011_1100_0000
		break;
	case "-A":
		out |= 0b0000_1100_1100_0000
		break;
	case "D+1":
		out |= 0b0000_0111_1100_0000
		break;
	case "A+1":
		out |= 0b0000_1101_1100_0000
		break;
	case "D-1":
		out |= 0b0000_0011_1000_0000
		break;
	case "A-1":
		out |= 0b0000_1100_1000_0000
		break;
	case "D+A":
		out |= 0b0000_0000_1000_0000
		break;
	case "D-A":
		out |= 0b0000_0100_1100_0000
		break;
	case "A-D":
		out |= 0b0000_0001_1100_0000
		break;
	case "D&A":
		out |= 0b0000_0000_0000_0000
		break;
	case "D|A":
		out |= 0b0000_0101_0100_0000
		break;

	// "a" flag set to 1
	case "M":
		out |= 0b0001_1100_0000_0000
		break;
	case "!M":
		out |= 0b0001_1100_0100_0000
		break;
	case "-M":
		out |= 0b0001_1100_1100_0000
		break;
	case "M+1":
		out |= 0b0001_1101_1100_0000
		break;
	case "M-1":
		out |= 0b0001_1100_1000_0000
		break;
	case "D+M":
		out |= 0b0001_0000_1000_0000
		break;
	case "D-M":
		out |= 0b0001_0100_1100_0000
		break;
	case "M-D":
		out |= 0b0001_0001_1100_0000
		break;
	case "D&M":
		out |= 0b0001_0000_0000_0000
		break;
	case "D|M":
		out |= 0b0001_0101_0100_0000
		break;
	}

	if isStore {
		var dest = line[:equalIndex]
		switch dest {
		case "M":
			out |= 0b0000_0000_0000_1000
			break
		case "D":
			out |= 0b0000_0000_0001_0000
			break
		case "MD":
			out |= 0b0000_0000_0001_1000
			break
		case "A":
			out |= 0b0000_0000_0010_0000
			break
		case "AM":
			out |= 0b0000_0000_0010_1000
			break
		case "AD":
			out |= 0b0000_0000_0011_0000
			break
		case "AMD":
			out |= 0b0000_0000_0011_1000
			break
		}
	}

	if isJump {
		var jump = line[semicolonIndex+1:]

		switch jump {
		case "JGT":
			out |= 0b0000_0000_0000_0001
			break
		case "JEQ":
			out |= 0b0000_0000_0000_0010
			break
		case "JGE":
			out |= 0b0000_0000_0000_0011
			break
		case "JLT":
			out |= 0b0000_0000_0000_0100
		case "JNE":
			out |= 0b0000_0000_0000_0101
			break
		case "JLE":
			out |= 0b0000_0000_0000_0110
			break
		case "JMP":
			out |= 0b0000_0000_0000_0111
			break
		}
	}

	return out
}