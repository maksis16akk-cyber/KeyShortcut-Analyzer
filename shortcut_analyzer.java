// shortcut_analyzer.java
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.stream.*;
import com.google.gson.*; // install gson

public class shortcut_analyzer {
    private static final String RESET = "\u001B[0m";
    private static final String GREEN = "\u001B[92m";
    private static final String RED = "\u001B[91m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String CYAN = "\u001B[96m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    static class Shortcut {
        String name;
        String shortcut;
        String app;
        String category;
    }

    static class Data {
        List<Shortcut> shortcuts = new ArrayList<>();
    }

    static class Analyzer {
        List<Shortcut> shortcuts;
        Map<String, List<Shortcut>> apps = new HashMap<>();
        Map<String, List<Shortcut>> categories = new HashMap<>();
        Map<String, List<Shortcut>> conflicts = new HashMap<>();

        Analyzer(Data data) {
            shortcuts = data.shortcuts;
            index();
        }

        void index() {
            for (Shortcut s : shortcuts) {
                String app = s.app == null || s.app.isEmpty() ? "Без приложения" : s.app;
                apps.computeIfAbsent(app, k -> new ArrayList<>()).add(s);

                String cat = s.category == null || s.category.isEmpty() ? "Без категории" : s.category;
                categories.computeIfAbsent(cat, k -> new ArrayList<>()).add(s);

                if (s.shortcut != null && !s.shortcut.isEmpty()) {
                    conflicts.computeIfAbsent(s.shortcut, k -> new ArrayList<>()).add(s);
                }
            }
        }

        List<Shortcut> search(String query) {
            String q = query.toLowerCase();
            return shortcuts.stream()
                .filter(s -> s.name.toLowerCase().contains(q) || s.shortcut.toLowerCase().contains(q))
                .collect(Collectors.toList());
        }

        Map<String, List<Shortcut>> getConflicts() {
            Map<String, List<Shortcut>> result = new HashMap<>();
            for (Map.Entry<String, List<Shortcut>> e : conflicts.entrySet()) {
                if (e.getValue().size() > 1) result.put(e.getKey(), e.getValue());
            }
            return result;
        }

        Map<String, List<Shortcut>> groupByApp() {
            return apps;
        }

        void exportCSV(String filename) throws IOException {
            try (PrintWriter pw = new PrintWriter(filename)) {
                pw.println("name,shortcut,app,category");
                for (Shortcut s : shortcuts) {
                    pw.printf("%s,%s,%s,%s\n", s.name, s.shortcut, s.app, s.category);
                }
            }
            System.out.println(colorize("Экспортировано в " + filename, GREEN));
        }

        void exportJSON(String filename) throws IOException {
            Gson gson = new GsonBuilder().setPrettyPrinting().create();
            Map<String, Object> root = new HashMap<>();
            root.put("shortcuts", shortcuts);
            String json = gson.toJson(root);
            Files.write(Paths.get(filename), json.getBytes());
            System.out.println(colorize("Экспортировано в " + filename, GREEN));
        }

        void display(List<Shortcut> items, boolean group, boolean conflictsOnly, boolean verbose) {
            List<Shortcut> data = items != null ? items : shortcuts;
            if (data.isEmpty()) {
                System.out.println(colorize("Нет данных для отображения.", YELLOW));
                return;
            }
            if (conflictsOnly) {
                Map<String, List<Shortcut>> conf = getConflicts();
                if (conf.isEmpty()) {
                    System.out.println(colorize("Конфликтов не найдено.", GREEN));
                    return;
                }
                System.out.println(colorize("🔍 Конфликты (одинаковые комбинации):", BOLD));
                for (Map.Entry<String, List<Shortcut>> e : conf.entrySet()) {
                    System.out.println(colorize("  " + e.getKey() + ":", RED));
                    for (Shortcut s : e.getValue()) {
                        System.out.println("    - " + s.name + " (" + s.app + ")");
                    }
                }
                return;
            }
            if (group) {
                Map<String, List<Shortcut>> grouped = groupByApp();
                for (Map.Entry<String, List<Shortcut>> e : grouped.entrySet()) {
                    System.out.println(colorize("\n📁 " + e.getKey() + " (" + e.getValue().size() + "):", BLUE));
                    for (Shortcut s : e.getValue()) {
                        printItem(s, verbose);
                    }
                }
            } else {
                for (Shortcut s : data) {
                    printItem(s, verbose);
                }
            }
        }

        void printItem(Shortcut s, boolean verbose) {
            if (verbose) {
                System.out.printf("  %s → %s  (приложение: %s, категория: %s)\n",
                    colorize(s.name, BOLD),
                    colorize(s.shortcut, CYAN),
                    s.app, s.category);
            } else {
                System.out.println("  " + s.name + " → " + s.shortcut);
            }
        }
    }

    public static void main(String[] args) throws IOException {
        String file = null, search = null, exportFile = null;
        boolean group = false, conflictsOnly = false, verbose = false;

        for (int i = 0; i < args.length; i++) {
            String arg = args[i];
            if (arg.equals("-s") && i+1 < args.length) search = args[++i];
            else if (arg.equals("-g")) group = true;
            else if (arg.equals("-c")) conflictsOnly = true;
            else if (arg.equals("-e") && i+1 < args.length) exportFile = args[++i];
            else if (arg.equals("-v")) verbose = true;
            else if (arg.equals("-h") || arg.equals("--help")) {
                System.out.println("Usage: java shortcut_analyzer <file> [-s search] [-g] [-c] [-e file] [-v]");
                return;
            } else if (file == null) file = arg;
        }
        if (file == null) {
            System.out.println(colorize("Укажите файл конфигурации.", RED));
            return;
        }

        Gson gson = new Gson();
        Data data;
        try {
            String json = new String(Files.readAllBytes(Paths.get(file)));
            data = gson.fromJson(json, Data.class);
        } catch (Exception e) {
            System.out.println(colorize("Ошибка загрузки файла: " + e.getMessage(), RED));
            return;
        }

        Analyzer analyzer = new Analyzer(data);

        if (search != null) {
            List<Shortcut> results = analyzer.search(search);
            if (!results.isEmpty()) {
                System.out.println(colorize("🔍 Найдено " + results.size() + " совпадений:", BOLD));
                analyzer.display(results, group, false, verbose);
            } else {
                System.out.println(colorize("Ничего не найдено.", YELLOW));
            }
            return;
        }

        if (conflictsOnly) {
            analyzer.display(null, false, true, false);
            return;
        }

        if (exportFile != null) {
            if (exportFile.endsWith(".csv")) analyzer.exportCSV(exportFile);
            else if (exportFile.endsWith(".json")) analyzer.exportJSON(exportFile);
            else System.out.println(colorize("Неизвестный формат экспорта. Используйте .csv или .json", RED));
            return;
        }

        analyzer.display(null, group, false, verbose);
    }
}
