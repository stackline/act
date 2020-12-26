package main

// $ go get github.com/PuerkitoBio/goquery
import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	flag.Parse()
	command := flag.Arg(0)

	switch command {
	case "get":
		content := flag.Arg(1)
		get(content)
	case "test":
		taskID := flag.Arg(1)
		sampleID := flag.Arg(2)
		test(taskID, sampleID)
	default:
		help()
	}
}

// ref. https://github.com/PuerkitoBio/goquery#examples
func get(content string) {
	fmt.Println("url: " + content)

	// Request the HTML page
	resp, err := http.Get(content)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the sample data
	// ref. https://img.atcoder.jp/public/0c9d249/js/contest.js
	var samplesByTasks [][]string
	doc.Find("#task-statement").Each(func(_ int, s1 *goquery.Selection) {
		var samples []string
		s1.Find("span.lang-ja h3+pre").Each(func(_ int, s2 *goquery.Selection) {
			samples = append(samples, s2.Text())
		})
		samplesByTasks = append(samplesByTasks, samples)
	})

	// Make sample directory
	err = makeDirIfNeeded("sample")
	if err != nil {
		log.Fatal(err)
	}

	// Write the sample data to files
	for i, samples := range samplesByTasks {
		for j, sample := range samples {
			filename := buildFileName(i, j)
			data := []byte(sample)

			fmt.Println("write file: " + filename)
			err := ioutil.WriteFile(filename, data, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func makeDirIfNeeded(dir string) error {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) || !fi.IsDir() {
		fmt.Println("mkdir: " + dir)
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// format:
//   sample/[task id]-[sample id]-(in|out).txt
// example:
//   sample/a-01-in.txt
//   sample/b-02-out.txt
func buildFileName(i int, j int) string {
	taskID := toTaskID(i)
	ioName := toIoName(j)
	sampleID := j/2 + 1
	return fmt.Sprintf("sample/%s-%02d-%s.txt", taskID, sampleID, ioName)
}

func toTaskID(i int) string {
	return string(rune('a' + i))
}

func toIoName(i int) string {
	if i%2 == 0 {
		return "in"
	}
	return "out"
}

func test(taskID string, sampleID string) {
	err := checkArguments(taskID, sampleID)
	if err != nil {
		log.Fatal(err)
	}

	inputFileName := taskID + ".cc"
	checkSumFileName := fmt.Sprintf("./cache/%s.sha512sum.txt", inputFileName)
	sampleInputFileName := fmt.Sprintf("./sample/%s-%02s-%s.txt", taskID, sampleID, "in")
	sampleOutputFileName := fmt.Sprintf("./sample/%s-%02s-%s.txt", taskID, sampleID, "out")
	executableFileName := fmt.Sprintf("./cache/%s.out", taskID)

	files := []string{inputFileName, sampleInputFileName, sampleOutputFileName}
	for _, f := range files {
		_, err = os.Stat(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make cache directory
	err = makeDirIfNeeded("cache")
	if err != nil {
		log.Fatal(err)
	}

	// Checksum
	content, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	sha512 := sha512.Sum512(content)
	encodedStr := hex.EncodeToString(sha512[:])

	needToCompile := true

	// Check the exisitence of checksum
	_, err = os.Stat(checkSumFileName)
	if !os.IsNotExist(err) {
		existenceCheckSum, err := ioutil.ReadFile(checkSumFileName)
		if err != nil {
			log.Fatal(err)
		}

		if encodedStr == string(existenceCheckSum) {
			needToCompile = false
		}
	}

	if needToCompile {
		fmt.Println("compile " + inputFileName)
		// Make checksum
		err = ioutil.WriteFile(checkSumFileName, []byte(encodedStr), 0644)
		if err != nil {
			log.Fatal(err)
		}

		/// Make binary file
		// -std=gnu++17       Conform to the ISO 2017 C++ standard with GNU extensions.
		// -Wall              Enable most warning messages.
		// -Wextra            Print extra (possibly unwanted) warnings.
		// -O<number>         Set optimization level to <number>.
		// -D<macro>[=<val>]  Define a <macro> with <val> as its value.
		//                    If just <macro> is given, <val> is taken to be 1.
		//
		// ref. https://atcoder.jp/contests/language-test-202001
		_, err = exec.Command("g++-9", "-std=gnu++17", "-Wall", "-Wextra", "-O2",
			"-DONLINE_JUDGE", "-o", executableFileName, inputFileName).Output()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Execute
	cmd := exec.Command(executableFileName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	f, err := os.Open(sampleInputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(stdin, string(b))

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	sampleOutputFile, err := os.Open(sampleOutputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sampleOutput, err := ioutil.ReadAll(sampleOutputFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[output]")
	fmt.Println(string(out))
	fmt.Println("[expected]")
	fmt.Println(string(sampleOutput))
	fmt.Println("[diff]")
	fmt.Println(string(out) == string(sampleOutput))
	fmt.Println("")
}

func checkArguments(taskID string, sampleID string) error {
	if taskID == "" {
		return errors.New("taskID is empty")
	}
	if sampleID == "" {
		return errors.New("sampleID is empty")
	}
	return nil
}

func help() {
	fmt.Println(
		`Act is a tool for AtCoder.

Usage:

	act <command> [arguments]

The commands are:

	act get <content>
		get sample data
	act test <taskID> <sampleID>
		hoge
	act help
		show help
	`)
}
