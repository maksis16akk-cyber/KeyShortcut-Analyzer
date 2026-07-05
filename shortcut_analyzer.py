# shortcut_analyzer.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
import json
import csv
import argparse
from collections import defaultdict

# ANSI-цвета
COLORS = {
    'reset': '\033[0m',
    'green': '\033[92m',
    'red': '\033[91m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'cyan': '\033[96m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

class ShortcutAnalyzer:
    def __init__(self, data):
        self.shortcuts = data.get('shortcuts', [])
        self.apps = defaultdict(list)
        self.categories = defaultdict(list)
        self.conflicts = {}
        self._index()

    def _index(self):
        for s in self.shortcuts:
            self.apps[s.get('app', 'Без приложения')].append(s)
            self.categories[s.get('category', 'Без категории')].append(s)
            combo = s.get('shortcut', '')
            if combo:
                if combo not in self.conflicts:
                    self.conflicts[combo] = []
                self.conflicts[combo].append(s)

    def search(self, query):
        query = query.lower()
        results = []
        for s in self.shortcuts:
            if query in s.get('name', '').lower() or query in s.get('shortcut', '').lower():
                results.append(s)
        return results

    def get_conflicts(self):
        return {k: v for k, v in self.conflicts.items() if len(v) > 1}

    def group_by_app(self):
        return self.apps

    def export_csv(self, filename):
        if not self.shortcuts:
            return
        keys = ['name', 'shortcut', 'app', 'category']
        with open(filename, 'w', newline='', encoding='utf-8') as f:
            writer = csv.DictWriter(f, fieldnames=keys)
            writer.writeheader()
            writer.writerows(self.shortcuts)
        print(colorize(f"Экспортировано в {filename}", 'green'))

    def export_json(self, filename):
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump({'shortcuts': self.shortcuts}, f, indent=2, ensure_ascii=False)
        print(colorize(f"Экспортировано в {filename}", 'green'))

    def display(self, shortcuts=None, group=False, conflicts_only=False, verbose=False):
        data = shortcuts if shortcuts is not None else self.shortcuts
        if not data:
            print(colorize("Нет данных для отображения.", 'yellow'))
            return

        if conflicts_only:
            conflicts = self.get_conflicts()
            if not conflicts:
                print(colorize("Конфликтов не найдено.", 'green'))
                return
            print(colorize("🔍 Конфликты (одинаковые комбинации):", 'bold'))
            for combo, items in conflicts.items():
                print(colorize(f"  {combo}:", 'red'))
                for item in items:
                    print(f"    - {item.get('name', 'Без имени')} ({item.get('app', 'Без приложения')})")
            return

        if group:
            grouped = self.group_by_app()
            for app, items in grouped.items():
                print(colorize(f"\n📁 {app} ({len(items)}):", 'blue'))
                for item in items:
                    self._print_item(item, verbose)
        else:
            for item in data:
                self._print_item(item, verbose)

    def _print_item(self, item, verbose=False):
        name = item.get('name', 'Без имени')
        shortcut = item.get('shortcut', '')
        app = item.get('app', 'Без приложения')
        category = item.get('category', 'Без категории')
        if verbose:
            print(f"  {colorize(name, 'bold')} → {colorize(shortcut, 'cyan')}  (приложение: {app}, категория: {category})")
        else:
            print(f"  {name} → {shortcut}")

def load_json(filename):
    with open(filename, 'r', encoding='utf-8') as f:
        return json.load(f)

def main():
    parser = argparse.ArgumentParser(description="KeyShortcut Analyzer")
    parser.add_argument('file', help='Файл конфигурации (JSON)')
    parser.add_argument('-s', '--search', help='Поиск по названию или комбинации')
    parser.add_argument('-g', '--group', action='store_true', help='Группировать по приложениям')
    parser.add_argument('-c', '--conflicts', action='store_true', help='Показать конфликты')
    parser.add_argument('-e', '--export', help='Экспортировать в файл (CSV или JSON)')
    parser.add_argument('-v', '--verbose', action='store_true', help='Подробный вывод')
    args = parser.parse_args()

    try:
        data = load_json(args.file)
    except Exception as e:
        print(colorize(f"Ошибка загрузки файла: {e}", 'red'))
        sys.exit(1)

    analyzer = ShortcutAnalyzer(data)

    if args.search:
        results = analyzer.search(args.search)
        if results:
            print(colorize(f"🔍 Найдено {len(results)} совпадений:", 'bold'))
            analyzer.display(results, args.group, args.conflicts, args.verbose)
        else:
            print(colorize("Ничего не найдено.", 'yellow'))
        return

    if args.conflicts:
        analyzer.display(conflicts_only=True)
        return

    if args.export:
        if args.export.endswith('.csv'):
            analyzer.export_csv(args.export)
        elif args.export.endswith('.json'):
            analyzer.export_json(args.export)
        else:
            print(colorize("Неизвестный формат экспорта. Используйте .csv или .json", 'red'))
        return

    analyzer.display(group=args.group, verbose=args.verbose)

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print(colorize("\nВыход.", 'yellow'))
        sys.exit(0)
