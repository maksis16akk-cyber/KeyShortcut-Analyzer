// shortcut_analyzer.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	reset  = "\033[0m"
	green  = "\033[92m"
	red    = "\033[91m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	cyan   = "\033[96m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

type Shortcut struct {
	Name     string `json:"name"`
	Shortcut string `json:"shortcut"`
	App      string `json:"app"`
	Category string `json:"category"`
}

type Data struct {
	Shortcuts []Shortcut `json:"shortcuts"`
}

type Analyzer struct {
	shortcuts []Shortcut
	apps      map[string][]Shortcut
	categories map[string][]Shortcut
	conflicts map[string][]Shortcut
}

func NewAnalyzer(data Data) *Analyzer {
	a := &Analyzer{
		shortcuts: data.Shortcuts,
		apps:      make(map[string][]Shortcut),
		categories: make(map[string][]Shortcut),
		conflicts: make(map[string][]Shortcut),
	}
	a.index()
	return a
}

func (a *Analyzer) index() {
	for _, s := range a.shortcuts {
		app := s.App
		if app == "" {
			app = "Без приложения"
		}
		a.apps[app] = append(a.apps[app], s)

		cat := s.Category
		if cat == "" {
			cat = "Без категории"
		}
		a.categories[cat] = append(a.categories[cat], s)

		if s.Shortcut != "" {
			a.conflicts[s.Shortcut] = append(a.conflicts[s.Shortcut], s)
		}
	}
}

func (a *Analyzer) Search(query string) []Shortcut {
	var results []Shortcut
	q := strings.ToLower(query)
	for _, s := range a.shortcuts {
		if strings.Contains(strings.ToLower(s.Name), q) || strings.Contains(strings.ToLower(s.Shortcut), q) {
			results = append(results, s)
		}
	}
	return results
}

func (a *Analyzer) GetConflicts() map[string][]Shortcut {
	conf := make(map[string][]Shortcut)
	for k, v := range a.conflicts {
		if len(v) > 1 {
			conf[k] = v
		}
	}
	return conf
}

func (a *Analyzer) GroupByApp() map[string][]Shortcut {
	return a.apps
}

func (a *Analyzer) ExportCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{"name", "shortcut", "app", "category"}
	writer.Write(headers)
	for _, s := range a.shortcuts {
		writer.Write([]string{s.Name, s.Shortcut, s.App, s.Category})
	}
	fmt.Println(colorize("Экспортировано в "+filename, green))
	return nil
}

func (a *Analyzer) ExportJSON(filename string) error {
	data := Data{Shortcuts: a.shortcuts}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonData, 0644)
}

func (a *Analyzer) Display(shortcuts []Shortcut, group, conflictsOnly, verbose bool) {
	data := shortcuts
	if data == nil {
		data = a.shortcuts
	}
	if len(data) == 0 {
		fmt.Println(colorize("Нет данных для отображения.", yellow))
		return
	}
	if conflictsOnly {
		conf := a.GetConflicts()
		if len(conf) == 0 {
			fmt.Println(colorize("Конфликтов не найдено.", green))
			return
		}
		fmt.Println(colorize("🔍 Конфликты (одинаковые комбинации):", bold))
		for combo, items := range conf {
			fmt.Printf("%s\n", colorize("  "+combo+":", red))
			for _, item := range items {
				fmt.Printf("    - %s (%s)\n", item.Name, item.App)
			}
		}
		return
	}
	if group {
		grouped := a.GroupByApp()
		for app, items := range grouped {
			fmt.Printf("%s\n", colorize(fmt.Sprintf("\n📁 %s (%d):", app, len(items)), blue))
			for _, item := range items {
				a.printItem(item, verbose)
			}
		}
	} else {
		for _, item := range data {
			a.printItem(item, verbose)
		}
	}
}

func (a *Analyzer) printItem(item Shortcut, verbose bool) {
	if verbose {
		fmt.Printf("  %s → %s  (приложение: %s, категория: %s)\n",
			colorize(item.Name, bold),
			colorize(item.Shortcut, cyan),
			item.App,
			item.Category)
	} else {
		fmt.Printf("  %s → %s\n", item.Name, item.Shortcut)
	}
}

func loadJSON(filename string) (Data, error) {
	var data Data
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bytes, &data)
	return data, err
}

func main() {
	var (
		file        string
		search      string
		group       bool
		conflicts   bool
		exportFile  string
		verbose     bool
	)
	flag.StringVar(&file, "f", "", "Файл конфигурации (JSON)")
	flag.StringVar(&search, "s", "", "Поиск по названию или комбинации")
	flag.BoolVar(&group, "g", false, "Группировать по приложениям")
	flag.BoolVar(&conflicts, "c", false, "Показать конфликты")
	flag.StringVar(&exportFile, "e", "", "Экспортировать в файл (CSV или JSON)")
	flag.BoolVar(&verbose, "v", false, "Подробный вывод")
	flag.Parse()

	if file == "" && flag.NArg() > 0 {
		file = flag.Arg(0)
	}
	if file == "" {
		fmt.Println(colorize("Укажите файл конфигурации.", red))
		flag.Usage()
		os.Exit(1)
	}

	data, err := loadJSON(file)
	if err != nil {
		fmt.Println(colorize("Ошибка загрузки файла: "+err.Error(), red))
		os.Exit(1)
	}

	analyzer := NewAnalyzer(data)

	if search != "" {
		results := analyzer.Search(search)
		if len(results) > 0 {
			fmt.Printf("%s\n", colorize(fmt.Sprintf("🔍 Найдено %d совпадений:", len(results)), bold))
			analyzer.Display(results, group, false, verbose)
		} else {
			fmt.Println(colorize("Ничего не найдено.", yellow))
		}
		return
	}

	if conflicts {
		analyzer.Display(nil, false, true, false)
		return
	}

	if exportFile != "" {
		var err error
		if strings.HasSuffix(exportFile, ".csv") {
			err = analyzer.ExportCSV(exportFile)
		} else if strings.HasSuffix(exportFile, ".json") {
			err = analyzer.ExportJSON(exportFile)
		} else {
			fmt.Println(colorize("Неизвестный формат экспорта. Используйте .csv или .json", red))
			return
		}
		if err != nil {
			fmt.Println(colorize("Ошибка экспорта: "+err.Error(), red))
		}
		return
	}

	analyzer.Display(nil, group, false, verbose)
}
