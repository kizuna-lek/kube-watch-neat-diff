package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	neat "github.com/itaysk/kubectl-neat/cmd"
	"github.com/mohae/deepcopy"
	"github.com/r3labs/diff/v3"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var (
	resource      = kingpin.Arg("resource-type", "Resource to watch").Required().String()
	name          = kingpin.Arg("resource-name", "Name of the resource").Required().String()
	diffWithFirst = kingpin.Flag("diff-with-first", "Diff with first version, instead of previous version, default false").Short('f').Bool()
	noColor       = kingpin.Flag("no-color", "Disable colored output").Bool()
)

func colorize(text, color string) string {
	if *noColor {
		return text
	}
	return color + text + Reset
}

func formatDiffOutput(changelog diff.Changelog) string {
	if len(changelog) == 0 {
		return colorize("No changes detected", Yellow) + "\n"
	}

	var output strings.Builder
	separator := strings.Repeat("=", 60)
	output.WriteString(colorize(separator, Cyan) + "\n")
	output.WriteString(colorize(fmt.Sprintf("Found %d changes:", len(changelog)), Bold+White) + "\n")
	output.WriteString(colorize(separator, Cyan) + "\n")

	for i, change := range changelog {
		path := strings.Join(change.Path, ".")
		if path == "" {
			path = "root"
		}

		output.WriteString(colorize(fmt.Sprintf("%d. ", i+1), Bold+White))

		switch change.Type {
		case diff.CREATE:
			output.WriteString(colorize("+ CREATED: ", Bold+Green))
			output.WriteString(colorize(path, Green) + "\n")
			output.WriteString(colorize("  + Value: ", Green))
			output.WriteString(colorize(formatValue(change.To), Bold+Green) + "\n")

		case diff.UPDATE:
			output.WriteString(colorize("~ UPDATED: ", Bold+Yellow))
			output.WriteString(colorize(path, Yellow) + "\n")
			output.WriteString(colorize("  - Old: ", Red))
			output.WriteString(colorize(formatValue(change.From), Bold+Red) + "\n")
			output.WriteString(colorize("  + New: ", Green))
			output.WriteString(colorize(formatValue(change.To), Bold+Green) + "\n")

		case diff.DELETE:
			output.WriteString(colorize("- DELETED: ", Bold+Red))
			output.WriteString(colorize(path, Red) + "\n")
			output.WriteString(colorize("  - Value: ", Red))
			output.WriteString(colorize(formatValue(change.From), Bold+Red) + "\n")

		}
	}

	output.WriteString("\n" + colorize(strings.Repeat("=", 61), Cyan) + "\n")
	return output.String()
}

func formatValue(value interface{}) string {
	if value == nil {
		return "<nil>"
	}

	switch v := value.(type) {
	case string:
		if len(v) > 100 {
			return fmt.Sprintf("%.97s...", v)
		}
		return fmt.Sprintf("\"%s\"", v)
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		if len(keys) > 3 {
			return fmt.Sprintf("map[%s... (%d keys)]", strings.Join(keys[:3], ", "), len(keys))
		}
		return fmt.Sprintf("map[%s]", strings.Join(keys, ", "))
	case []interface{}:
		return fmt.Sprintf("array[%d items]", len(v))
	default:
		str := fmt.Sprintf("%v", v)
		if len(str) > 100 {
			return fmt.Sprintf("%.97s...", str)
		}
		return str
	}
}

func supportsColor() bool {
	// 检查是否为 TTY 终端
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	// 检查环境变量
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	return true
}

func init() {
	log.SetFlags(0)
	kingpin.HelpFlag.Short('\n')
	kingpin.Version("1.0.0")
	kingpin.Parse()

	if !supportsColor() {
		*noColor = true
	}
}

func main() {
	log.Println("Starting kube-watch-neat-diff")

	cmd := exec.Command("kubectl", "get", "-w", *resource, *name, "-o=json")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}

	var preObj, obj map[string]interface{}
	dec := json.NewDecoder(stdout)
	for dec.More() {
		var raw map[string]interface{}
		err = dec.Decode(&raw)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		var b []byte
		b, err = json.Marshal(raw)
		if err != nil {
			log.Println("Error marshaling object to JSON:", err)
			continue
		}

		b, err = neat.NeatYAMLOrJSON(b)
		if err != nil {
			log.Println("Error processing object:", err)
			continue
		}

		err = json.Unmarshal(b, &obj)
		if err != nil {
			log.Println("Error parsing previous object:", err)
			continue
		}

		if preObj == nil {
			preObj = deepcopy.Copy(obj).(map[string]interface{})
			log.Println("Watching resource, waiting for changes...")
		} else {
			var changelog diff.Changelog
			changelog, err = diff.Diff(preObj, obj)
			if err != nil {
				log.Println("Error computing diff:", err)
			} else {
				fmt.Print(formatDiffOutput(changelog))
			}

			if !*diffWithFirst {
				preObj = deepcopy.Copy(obj).(map[string]interface{})
			}
		}
	}

	if err = cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}
