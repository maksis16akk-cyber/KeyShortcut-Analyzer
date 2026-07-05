// shortcut_analyzer.cpp
#include <iostream>
#include <vector>
#include <string>
#include <map>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <cctype>
#include <json/json.h> // sudo apt-get install libjsoncpp-dev

using namespace std;

const string RESET = "\033[0m";
const string GREEN = "\033[92m";
const string RED = "\033[91m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string CYAN = "\033[96m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

struct Shortcut {
    string name;
    string shortcut;
    string app;
    string category;
};

class Analyzer {
public:
    vector<Shortcut> shortcuts;
    map<string, vector<Shortcut>> apps;
    map<string, vector<Shortcut>> categories;
    map<string, vector<Shortcut>> conflicts;

    Analyzer(const Json::Value& root) {
        auto& arr = root["shortcuts"];
        for (const auto& item : arr) {
            Shortcut s;
            s.name = item["name"].asString();
            s.shortcut = item["shortcut"].asString();
            s.app = item["app"].asString();
            s.category = item["category"].asString();
            shortcuts.push_back(s);
        }
        index();
    }

    void index() {
        for (auto& s : shortcuts) {
            string app = s.app.empty() ? "Без приложения" : s.app;
            apps[app].push_back(s);

            string cat = s.category.empty() ? "Без категории" : s.category;
            categories[cat].push_back(s);

            if (!s.shortcut.empty()) {
                conflicts[s.shortcut].push_back(s);
            }
        }
    }

    vector<Shortcut> search(const string& query) {
        vector<Shortcut> results;
        string q = query;
        transform(q.begin(), q.end(), q.begin(), ::tolower);
        for (auto& s : shortcuts) {
            string name = s.name;
            transform(name.begin(), name.end(), name.begin(), ::tolower);
            string sc = s.shortcut;
            transform(sc.begin(), sc.end(), sc.begin(), ::tolower);
            if (name.find(q) != string::npos || sc.find(q) != string::npos) {
                results.push_back(s);
            }
        }
        return results;
    }

    map<string, vector<Shortcut>> getConflicts() {
        map<string, vector<Shortcut>> result;
        for (auto& kv : conflicts) {
            if (kv.second.size() > 1) {
                result[kv.first] = kv.second;
            }
        }
        return result;
    }

    void exportCSV(const string& filename) {
        ofstream f(filename);
        if (!f) { cerr << colorize("Ошибка создания файла", RED) << endl; return; }
        f << "name,shortcut,app,category\n";
        for (auto& s : shortcuts) {
            f << s.name << "," << s.shortcut << "," << s.app << "," << s.category << "\n";
        }
        cout << colorize("Экспортировано в " + filename, GREEN) << endl;
    }

    void display(const vector<Shortcut>* data, bool group, bool conflictsOnly, bool verbose) {
        const vector<Shortcut>* items = data ? data : &shortcuts;
        if (items->empty()) {
            cout << colorize("Нет данных для отображения.", YELLOW) << endl;
            return;
        }
        if (conflictsOnly) {
            auto conf = getConflicts();
            if (conf.empty()) {
                cout << colorize("Конфликтов не найдено.", GREEN) << endl;
                return;
            }
            cout << colorize("🔍 Конфликты (одинаковые комбинации):", BOLD) << endl;
            for (auto& kv : conf) {
                cout << colorize("  " + kv.first + ":", RED) << endl;
                for (auto& s : kv.second) {
                    cout << "    - " << s.name << " (" << s.app << ")" << endl;
                }
            }
            return;
        }
        if (group) {
            for (auto& kv : apps) {
                cout << colorize("\n📁 " + kv.first + " (" + to_string(kv.second.size()) + "):", BLUE) << endl;
                for (auto& s : kv.second) {
                    printItem(s, verbose);
                }
            }
        } else {
            for (auto& s : *items) {
                printItem(s, verbose);
            }
        }
    }

    void printItem(const Shortcut& s, bool verbose) {
        if (verbose) {
            cout << "  " << colorize(s.name, BOLD) << " → " << colorize(s.shortcut, CYAN)
                 << "  (приложение: " << s.app << ", категория: " << s.category << ")" << endl;
        } else {
            cout << "  " << s.name << " → " << s.shortcut << endl;
        }
    }
};

Json::Value loadJSON(const string& filename) {
    ifstream f(filename);
    Json::Value root;
    f >> root;
    return root;
}

int main(int argc, char* argv[]) {
    string file, search, exportFile;
    bool group = false, conflictsOnly = false, verbose = false;

    for (int i=1; i<argc; ++i) {
        string arg = argv[i];
        if (arg == "-s" && i+1 < argc) search = argv[++i];
        else if (arg == "-g") group = true;
        else if (arg == "-c") conflictsOnly = true;
        else if (arg == "-e" && i+1 < argc) exportFile = argv[++i];
        else if (arg == "-v") verbose = true;
        else if (arg == "-h") {
            cout << "Usage: shortcut_analyzer <file> [-s search] [-g] [-c] [-e file] [-v]" << endl;
            return 0;
        } else if (file.empty()) file = arg;
    }
    if (file.empty()) {
        cerr << colorize("Укажите файл конфигурации.", RED) << endl;
        return 1;
    }

    Json::Value root;
    try {
        root = loadJSON(file);
    } catch (const exception& e) {
        cerr << colorize("Ошибка загрузки файла: " + string(e.what()), RED) << endl;
        return 1;
    }

    Analyzer analyzer(root);

    if (!search.empty()) {
        auto results = analyzer.search(search);
        if (!results.empty()) {
            cout << colorize("🔍 Найдено " + to_string(results.size()) + " совпадений:", BOLD) << endl;
            analyzer.display(&results, group, false, verbose);
        } else {
            cout << colorize("Ничего не найдено.", YELLOW) << endl;
        }
        return 0;
    }

    if (conflictsOnly) {
        analyzer.display(nullptr, false, true, false);
        return 0;
    }

    if (!exportFile.empty()) {
        if (exportFile.find(".csv") != string::npos) {
            analyzer.exportCSV(exportFile);
        } else if (exportFile.find(".json") != string::npos) {
            // Для JSON используем простой вывод, в реальном проекте можно улучшить
            ofstream f(exportFile);
            f << root.toStyledString();
            cout << colorize("Экспортировано в " + exportFile, GREEN) << endl;
        } else {
            cerr << colorize("Неизвестный формат экспорта. Используйте .csv или .json", RED) << endl;
        }
        return 0;
    }

    analyzer.display(nullptr, group, false, verbose);
    return 0;
}
