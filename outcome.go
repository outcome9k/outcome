package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Color and style ANSI codes
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
)

// Tool represents one downloadable tool
type Tool struct {
	Name string
	Desc string
	URL  string
}

var tools = map[string]Tool{
	"1": {"black2", "Black2 Obfuscator", "https://raw.githubusercontent.com/outcome9k/test/main/black2"},
	"2": {"fix", "Fix Utility", "https://raw.githubusercontent.com/outcome9k/test/main/fix"},
	"3": {"hyperion", "Hyperion Obfuscator", "https://raw.githubusercontent.com/outcome9k/test/main/hyperion"},
	"4": {"kramer", "KRAMER Obfuscator", "https://raw.githubusercontent.com/outcome9k/test/main/kramer"},
	"5": {"pymor", "PYMOR Obfuscator", "https://raw.githubusercontent.com/outcome9k/test/main/pymor"},
	"6": {"emo", "EMO Tool", "https://raw.githubusercontent.com/outcome9k/test/main/emo"},
	"7": {"error", "Error Utility", "https://raw.githubusercontent.com/outcome9k/test/main/error"},
	"8": {"user", "User Tool", "https://raw.githubusercontent.com/outcome9k/test/main/user"},
}

// clearScreen clears terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// printMenu shows the menu with color
func printMenu() {
	fmt.Println(Cyan + Bold + "╔══════════════════════════╗" + Reset)
	fmt.Println("║      9k CLI Toolkit      ║")
	fmt.Println("╚══════════════════════════╝")

	// Print two columns if possible
	keys := []string{}
	for k := range tools {
		keys = append(keys, k)
	}
	half := len(tools)/2 + len(tools)%2

	for i := 0; i < half; i++ {
		left := tools[keys[i]]
		right := Tool{}
		if i+half < len(keys) {
			right = tools[keys[i+half]]
		}

		leftText := fmt.Sprintf("%s%s%s. %s%s", Green, Bold, keys[i], left.Desc, Reset)
		rightText := ""
		if right.Name != "" {
			rightText = fmt.Sprintf("%s%s%s. %s%s", Green, Bold, keys[i+half], right.Desc, Reset)
		}

		fmt.Printf("%-30s %s\n", leftText, rightText)
	}

	fmt.Println()
	fmt.Println(Green + Bold + "a. Run ALL tools" + Reset)
	fmt.Println(Red + Bold + "0. Exit" + Reset)
	fmt.Println(strings.Repeat("═", 30))
}

// downloadTool downloads the tool content from URL
func downloadTool(tool Tool) (string, error) {
	fmt.Printf("%sDownloading %s...%s\n", Cyan, tool.Name, Reset)
	resp, err := http.Get(tool.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("%sDownloaded %s%s\n", Green, tool.Name, Reset)
	return string(body), nil
}

// runToolWithArgs prompts for input/output file paths and runs the tool with proper args
func runToolWithArgs(name, content string) error {
	tmpFile, err := os.CreateTemp("", name+"-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return err
	}
	tmpFile.Close()

	isPython := strings.Contains(content, "import ") || strings.Contains(content, "def ")

	reader := bufio.NewReader(os.Stdin)

	// Prompt user for input file path
	var inputFile string
	for {
		fmt.Printf("%sEnter the path of the Python file to obfuscate (input file): %s", Yellow, Reset)
		inputFile, _ = reader.ReadString('\n')
		inputFile = strings.TrimSpace(inputFile)
		if inputFile == "" {
			fmt.Println(Red + "Input file path cannot be empty!" + Reset)
			continue
		}
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Println(Red + "Input file does not exist. Please enter a valid path." + Reset)
			continue
		}
		break
	}

	// Prompt for optional output file path
	fmt.Printf("%sEnter output file path (optional, press Enter to skip): %s", Yellow, Reset)
	outputFile, _ := reader.ReadString('\n')
	outputFile = strings.TrimSpace(outputFile)

	// Build argument slice with --input/-i and optional --output/-o
	args := []string{}

	// Add --input or -i
	args = append(args, "--input", inputFile)

	// Add output if provided
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}

	var cmd *exec.Cmd
	if isPython {
		cmd = exec.Command("python3", append([]string{tmpFile.Name()}, args...)...)
	} else {
		cmd = exec.Command("bash", append([]string{tmpFile.Name()}, args...)...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf(Red+"Execution failed: %v\n"+Reset, err)
	}
	return err
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printMenu()
		fmt.Print(Green + Bold + "Select an option: " + Reset)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "0" {
			fmt.Println(Red + Bold + "Exiting..." + Reset)
			break
		} else if choice == "a" {
			for _, tool := range tools {
				fmt.Printf("%sRunning %s...%s\n", Yellow, tool.Name, Reset)
				content, err := downloadTool(tool)
				if err != nil {
					fmt.Printf("%sFailed to download %s: %v%s\n", Red, tool.Name, err, Reset)
					continue
				}
				if err := runToolWithArgs(tool.Name, content); err != nil {
					fmt.Printf("%sFailed to run %s: %v%s\n", Red, tool.Name, err, Reset)
				}
				fmt.Println(strings.Repeat("-", 30))
				fmt.Print("Press Enter to continue...")
				reader.ReadString('\n')
			}
		} else if tool, ok := tools[choice]; ok {
			content, err := downloadTool(tool)
			if err != nil {
				fmt.Printf("%sFailed to download %s: %v%s\n", Red, tool.Name, err, Reset)
				fmt.Print("Press Enter to continue...")
				reader.ReadString('\n')
				continue
			}
			if err := runToolWithArgs(tool.Name, content); err != nil {
				fmt.Printf("%sFailed to run %s: %v%s\n", Red, tool.Name, err, Reset)
			}
			fmt.Print("Press Enter to continue...")
			reader.ReadString('\n')
		} else {
			fmt.Println(Red + "Invalid selection!" + Reset)
			fmt.Print("Press Enter to try again...")
			reader.ReadString('\n')
		}
	}
}
