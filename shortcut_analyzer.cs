// shortcut_analyzer.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;

class Shortcut
{
    public string name { get; set; }
    public string shortcut { get; set; }
    public string app { get; set; }
    public string category { get; set; }
}

class Data
{
    public List<Shortcut> shortcuts { get; set; }
}

class Analyzer
{
    private List<Shortcut> shortcuts;
    private Dictionary<string, List<Shortcut>> apps = new();
    private Dictionary<string, List<Shortcut>> categories = new();
    private Dictionary<string, List<Shortcut>> conflicts = new();

    public Analyzer(Data data)
    {
        shortcuts = data.shortcuts ?? new List<Shortcut>();
        Index();
    }

    private void Index()
    {
        foreach (var s in shortcuts)
        {
            string app = string.IsNullOrEmpty(s.app) ? "Без приложения" : s.app;
            if (!apps.ContainsKey(app)) apps[app] = new List<Shortcut>();
            apps[app].Add(s);

            string cat = string.IsNullOrEmpty(s.category) ? "Без категории" : s.category;
            if (!categories.ContainsKey(cat)) categories[cat] = new List<Shortcut>();
            categories[cat].Add(s);

            if (!string.IsNullOrEmpty(s.shortcut))
            {
                if (!conflicts.ContainsKey(s.shortcut)) conflicts[s.shortcut] = new List<Shortcut>();
                conflicts[s.shortcut].Add(s);
            }
        }
    }

    public List<Shortcut> Search(string query)
    {
        string q = query.ToLower();
        return shortcuts.Where(s =>
            s.name.ToLower().Contains(q) ||
            s.shortcut.ToLower().Contains(q)
        ).ToList();
    }

    public Dictionary<string, List<Shortcut>> GetConflicts()
    {
        var result = new Dictionary<string, List<Shortcut>>();
        foreach (var kv in conflicts)
        {
            if (kv.Value.Count > 1) result[kv.Key] = kv.Value;
        }
        return result;
    }

    public Dictionary<string, List<Shortcut>> GroupByApp() => apps;

    public void ExportCSV(string filename)
    {
        using var writer = new StreamWriter(filename);
        writer.WriteLine("name,shortcut,app,category");
        foreach (var s in shortcuts)
            writer.WriteLine($"{s.name},{s.shortcut},{s.app},{s.category}");
        Console.WriteLine(Colorize($"Экспортировано в {filename}", "green"));
    }

    public void ExportJSON(string filename)
    {
        var data = new Data { shortcuts = shortcuts };
        var json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine(Colorize($"Экспортировано в {filename}", "green"));
    }

    public void Display(List<Shortcut> items, bool group, bool conflictsOnly, bool verbose)
    {
        var data = items ?? shortcuts;
        if (data.Count == 0)
        {
            Console.WriteLine(Colorize("Нет данных для отображения.", "yellow"));
            return;
        }
        if (conflictsOnly)
        {
            var conf = GetConflicts();
            if (conf.Count == 0)
            {
                Console.WriteLine(Colorize("Конфликтов не найдено.", "green"));
                return;
            }
            Console.WriteLine(Colorize("🔍 Конфликты (одинаковые комбинации):", "bold"));
            foreach (var kv in conf)
            {
                Console.WriteLine(Colorize($"  {kv.Key}:", "red"));
                foreach (var s in kv.Value)
                    Console.WriteLine($"    - {s.name} ({s.app})");
            }
            return;
        }
        if (group)
        {
            var grouped = GroupByApp();
            foreach (var kv in grouped)
            {
                Console.WriteLine(Colorize($"\n📁 {kv.Key} ({kv.Value.Count}):", "blue"));
                foreach (var s in kv.Value)
                    PrintItem(s, verbose);
            }
        }
        else
        {
            foreach (var s in data)
                PrintItem(s, verbose);
        }
    }

    private void PrintItem(Shortcut s, bool verbose)
    {
        if (verbose)
            Console.WriteLine($"  {Colorize(s.name, "bold")} → {Colorize(s.shortcut, "cyan")}  (приложение: {s.app}, категория: {s.category})");
        else
            Console.WriteLine($"  {s.name} → {s.shortcut}");
    }

    private static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "green" => "\x1b[92m",
            "red" => "\x1b[91m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "cyan" => "\x1b[96m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    static void Main(string[] args)
    {
        string file = null, search = null, exportFile = null;
        bool group = false, conflictsOnly = false, verbose = false;

        for (int i = 0; i < args.Length; i++)
        {
            string arg = args[i];
            if (arg == "-s" && i+1 < args.Length) search = args[++i];
            else if (arg == "-g") group = true;
            else if (arg == "-c") conflictsOnly = true;
            else if (arg == "-e" && i+1 < args.Length) exportFile = args[++i];
            else if (arg == "-v") verbose = true;
            else if (arg == "-h" || arg == "--help")
            {
                Console.WriteLine("Usage: shortcut_analyzer <file> [-s search] [-g] [-c] [-e file] [-v]");
                return;
            }
            else if (file == null) file = arg;
        }
        if (file == null)
        {
            Console.WriteLine(Colorize("Укажите файл конфигурации.", "red"));
            return;
        }

        Data data;
        try
        {
            string json = File.ReadAllText(file);
            data = JsonSerializer.Deserialize<Data>(json);
        }
        catch (Exception e)
        {
            Console.WriteLine(Colorize($"Ошибка загрузки файла: {e.Message}", "red"));
            return;
        }

        var analyzer = new Analyzer(data);

        if (search != null)
        {
            var results = analyzer.Search(search);
            if (results.Count > 0)
            {
                Console.WriteLine(Colorize($"🔍 Найдено {results.Count} совпадений:", "bold"));
                analyzer.Display(results, group, false, verbose);
            }
            else Console.WriteLine(Colorize("Ничего не найдено.", "yellow"));
            return;
        }

        if (conflictsOnly)
        {
            analyzer.Display(null, false, true, false);
            return;
        }

        if (exportFile != null)
        {
            if (exportFile.EndsWith(".csv")) analyzer.ExportCSV(exportFile);
            else if (exportFile.EndsWith(".json")) analyzer.ExportJSON(exportFile);
            else Console.WriteLine(Colorize("Неизвестный формат экспорта. Используйте .csv или .json", "red"));
            return;
        }

        analyzer.Display(null, group, false, verbose);
    }
}
