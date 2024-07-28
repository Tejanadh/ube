package main

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"
	"ube/src/cmd"
	"ube/src/language"
	"ube/src/tui"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type languageInfo struct {
	lines int
	files int
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		args := os.Args[1:] // Skip "ube"
		if contains(args, "--help") || contains(args, "-h") || contains(args, "--version") || contains(args, "-v") {
			os.Exit(0)
		}
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		log.Fatal("ube: no such file or directory")
		os.Exit(1)
	}

	path := os.Args[1]

	m := tui.Model{ElapsedTime: stopwatch.NewWithInterval(time.Millisecond), IsRunning: true}
	p := tea.NewProgram(m)

	go func() {
		msg := getMessage(path)
		p.Send(msg)
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getMessage(path string) tea.Msg {
	stats, err := countLinesOfCode(path)
	if err != nil {
		log.Error(err)
		return tui.CompletionResult{Err: err}
	}

	t := generateTable(stats)
	h := help.New()
	return tui.CompletionResult{Table: t, Help: h}
}

func isValidFile(fileName string, info fs.DirEntry) bool {
	if info.IsDir() || !info.Type().IsRegular() {
		return false
	}
	_, exists := getLanguageName(fileName)
	return exists
}

func getLanguageName(fileName string) (string, bool) {
	language, exists := language.Extensions[getLanguageExtension(fileName)]
	if !exists {
		return "", false
	}
	return language, true
}

func getLanguageExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func countLinesOfCode(path string) (map[string]languageInfo, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	lineCount := make(map[string]languageInfo)

	err := filepath.WalkDir(path, func(currPath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !isValidFile(currPath, info) {
			return nil
		}
		language, _ := getLanguageName(currPath)

		wg.Add(1)
		go func() {
			defer wg.Done()

			lines, err := countLinesOfFile(currPath)
			if err != nil {
				log.Error(err)
				return
			}

			mu.Lock()
			ld, exists := lineCount[language]
			if !exists {
				ld = languageInfo{}
			}
			ld.lines += lines
			ld.files++
			lineCount[language] = ld
			mu.Unlock()
		}()

		return nil
	})

	if err != nil {
		return lineCount, err
	}

	wg.Wait()
	return lineCount, nil
}

func countLinesOfFile(filePath string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	var count int
	var read int
	var target []byte = []byte("\n")

	buffer := make([]byte, 32*1024)

	for {
		read, err = file.Read(buffer)
		if err != nil {
			break
		}

		count += bytes.Count(buffer[:read], target)
	}

	if err == io.EOF {
		return count, nil
	}

	return count, err
}

func generateTable(lineCount map[string]languageInfo) table.Model {
	columns := []table.Column{
		{Title: "Language", Width: 16},
		{Title: "Lines", Width: 16},
		{Title: "Files", Width: 10},
	}

	rows := []table.Row{}
	lineTotal := 0
	fileTotal := 0
	for language, info := range lineCount {
		rows = append(rows, table.Row{language, strconv.Itoa(info.lines), strconv.Itoa(info.files)})
		lineTotal += info.lines
		fileTotal += info.files
	}

	// Sort by lines of code
	sort.Slice(rows, func(i, j int) bool {
		li1, _ := strconv.Atoi(rows[i][1])
		li2, _ := strconv.Atoi(rows[j][1])
		return li1 > li2
	})

	for i, row := range rows {
		rows[i][1] = FormatStringInteger(row[1])
		rows[i][2] = FormatStringInteger(row[2])
	}

	rows = append(rows, table.Row{"", "", ""})
	rows = append(rows, table.Row{"Total", FormatStringInteger(strconv.Itoa(lineTotal)), FormatStringInteger(strconv.Itoa(fileTotal))})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#5433ff")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func FormatStringInteger(n string) string {
	if len(n) < 4 {
		return n
	}

	var formatted string
	for i, r := range n {
		if i != 0 && (len(n)-i)%3 == 0 {
			formatted += ","
		}
		formatted += string(r)
	}

	return formatted
}
