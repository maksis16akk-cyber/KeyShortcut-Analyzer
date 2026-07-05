// shortcut_analyzer.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const readline = require('readline');

const COLORS = {
    reset: '\x1b[0m',
    green: '\x1b[92m',
    red: '\x1b[91m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    cyan: '\x1b[96m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

class Analyzer {
    constructor(data) {
        this.shortcuts = data.shortcuts || [];
        this.apps = {};
        this.categories = {};
        this.conflicts = {};
        this.index();
    }

    index() {
        for (const s of this.shortcuts) {
            const app = s.app || 'Без приложения';
            if (!this.apps[app]) this.apps[app] = [];
            this.apps[app].push(s);

            const cat = s.category || 'Без категории';
            if (!this.categories[cat]) this.categories[cat] = [];
            this.categories[cat].push(s);

            if (s.shortcut) {
                if (!this.conflicts[s.shortcut]) this.conflicts[s.shortcut] = [];
                this.conflicts[s.shortcut].push(s);
            }
        }
    }

    search(query) {
        const q = query.toLowerCase();
        return this.shortcuts.filter(s =>
            s.name.toLowerCase().includes(q) ||
            s.shortcut.toLowerCase().includes(q)
        );
    }

    getConflicts() {
        const result = {};
        for (const [key, val] of Object.entries(this.conflicts)) {
            if (val.length > 1) result[key] = val;
        }
        return result;
    }

    groupByApp() {
        return this.apps;
    }

    exportCSV(filename) {
        const lines = ['name,shortcut,app,category'];
        for (const s of this.shortcuts) {
            lines.push(`${s.name},${s.shortcut},${s.app||''},${s.category||''}`);
        }
        fs.writeFileSync(filename, lines.join('\n'));
        console.log(colorize(`Экспортировано в ${filename}`, 'green'));
    }

    exportJSON(filename) {
        fs.writeFileSync(filename, JSON.stringify({ shortcuts: this.shortcuts }, null, 2));
        console.log(colorize(`Экспортировано в ${filename}`, 'green'));
    }

    display(items, group, conflictsOnly, verbose) {
        const data = items || this.shortcuts;
        if (data.length === 0) {
            console.log(colorize('Нет данных для отображения.', 'yellow'));
            return;
        }
        if (conflictsOnly) {
            const conf = this.getConflicts();
            const keys = Object.keys(conf);
            if (keys.length === 0) {
                console.log(colorize('Конфликтов не найдено.', 'green'));
                return;
            }
            console.log(colorize('🔍 Конфликты (одинаковые комбинации):', 'bold'));
            for (const combo of keys) {
                console.log(colorize(`  ${combo}:`, 'red'));
                for (const item of conf[combo]) {
                    console.log(`    - ${item.name} (${item.app || 'Без приложения'})`);
                }
            }
            return;
        }
        if (group) {
            const grouped = this.groupByApp();
            for (const [app, items] of Object.entries(grouped)) {
                console.log(colorize(`\n📁 ${app} (${items.length}):`, 'blue'));
                for (const item of items) {
                    this.printItem(item, verbose);
                }
            }
        } else {
            for (const item of data) {
                this.printItem(item, verbose);
            }
        }
    }

    printItem(item, verbose) {
        if (verbose) {
            console.log(`  ${colorize(item.name, 'bold')} → ${colorize(item.shortcut, 'cyan')}  (приложение: ${item.app || 'Без приложения'}, категория: ${item.category || 'Без категории'})`);
        } else {
            console.log(`  ${item.name} → ${item.shortcut}`);
        }
    }
}

function main() {
    const args = process.argv.slice(2);
    let file = null, search = null, exportFile = null;
    let group = false, conflictsOnly = false, verbose = false;

    for (let i = 0; i < args.length; i++) {
        const arg = args[i];
        if (arg === '-s' && i+1 < args.length) search = args[++i];
        else if (arg === '-g') group = true;
        else if (arg === '-c') conflictsOnly = true;
        else if (arg === '-e' && i+1 < args.length) exportFile = args[++i];
        else if (arg === '-v') verbose = true;
        else if (arg === '-h' || arg === '--help') {
            console.log('Usage: node shortcut_analyzer.js <file> [-s search] [-g] [-c] [-e file] [-v]');
            process.exit(0);
        } else if (!file) file = arg;
    }
    if (!file) {
        console.log(colorize('Укажите файл конфигурации.', 'red'));
        process.exit(1);
    }

    let data;
    try {
        data = JSON.parse(fs.readFileSync(file, 'utf8'));
    } catch (err) {
        console.log(colorize(`Ошибка загрузки файла: ${err.message}`, 'red'));
        process.exit(1);
    }

    const analyzer = new Analyzer(data);

    if (search) {
        const results = analyzer.search(search);
        if (results.length) {
            console.log(colorize(`🔍 Найдено ${results.length} совпадений:`, 'bold'));
            analyzer.display(results, group, false, verbose);
        } else {
            console.log(colorize('Ничего не найдено.', 'yellow'));
        }
        return;
    }

    if (conflictsOnly) {
        analyzer.display(null, false, true, false);
        return;
    }

    if (exportFile) {
        if (exportFile.endsWith('.csv')) {
            analyzer.exportCSV(exportFile);
        } else if (exportFile.endsWith('.json')) {
            analyzer.exportJSON(exportFile);
        } else {
            console.log(colorize('Неизвестный формат экспорта. Используйте .csv или .json', 'red'));
        }
        return;
    }

    analyzer.display(null, group, false, verbose);
}

main();
