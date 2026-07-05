⌨️ KeyShortcut Analyzer – Анализатор клавиатурных сокращений на 7 языках
Инструмент для анализа, поиска и управления клавиатурными сокращениями в ваших приложениях.
Загружайте файлы конфигурации, находите конфликты, группируйте по приложениям, экспортируйте статистику.
Реализован на 7 языках программирования с единым интерфейсом командной строки и расширенными возможностями.

🚀 Возможности
Загрузка сокращений – чтение из JSON-файлов с описанием горячих клавиш.

Поиск – быстрый поиск по названию или комбинации клавиш.

Группировка – сортировка сокращений по приложениям или категориям.

Обнаружение конфликтов – поиск дублирующихся комбинаций клавиш.

Статистика – количество сокращений, самые популярные комбинации.

Экспорт – сохранение результатов в CSV или JSON.

Цветной вывод – наглядное отображение в терминале.

Интерактивный режим – просмотр и фильтрация в реальном времени.

Кроссплатформенность – работает на Windows, Linux и macOS.

📖 Использование
Синтаксис (единый для всех версий):

bash
<команда> <файл_конфигурации> [опции]
Опции
Опция	Описание
-s, --search <текст>	Поиск по названию или комбинации
-g, --group	Группировать по приложениям
-c, --conflicts	Показать только конфликтующие сокращения
-e, --export <файл>	Экспортировать в CSV или JSON
-v, --verbose	Подробный вывод
-h, --help	Справка
Примеры
bash
# Анализ файла shortcuts.json
python shortcut_analyzer.py shortcuts.json

# Поиск сокращений с "Ctrl"
python shortcut_analyzer.py shortcuts.json -s Ctrl

# Группировка по приложениям и экспорт в CSV
python shortcut_analyzer.py shortcuts.json -g -e report.csv

# Поиск конфликтов
python shortcut_analyzer.py shortcuts.json -c
🛠 Установка и запуск
Python
bash
python shortcut_analyzer.py <файл> [опции]
Требуется Python 3.6+.

C++
bash
g++ -std=c++17 shortcut_analyzer.cpp -o shortcut_analyzer
./shortcut_analyzer <файл> [опции]
Go
bash
go build shortcut_analyzer.go
./shortcut_analyzer <файл> [опции]
JavaScript (Node.js)
bash
node shortcut_analyzer.js <файл> [опции]
C#
bash
csc shortcut_analyzer.cs
mono shortcut_analyzer.exe <файл> [опции]   # или dotnet run
Ruby
bash
ruby shortcut_analyzer.rb <файл> [опции]
Java
bash
javac shortcut_analyzer.java
java shortcut_analyzer <файл> [опции]
🧠 Формат входного файла
Программа ожидает JSON-файл со следующей структурой:

json
{
  "shortcuts": [
    {
      "name": "Сохранить",
      "shortcut": "Ctrl+S",
      "app": "Текстовый редактор",
      "category": "Файл"
    },
    {
      "name": "Копировать",
      "shortcut": "Ctrl+C",
      "app": "Система",
      "category": "Правка"
    }
  ]
}
✨ Дополнительные фичи
Автоопределение формата – поддержка JSON и CSV на входе.

Интерактивный режим – возможность фильтрации и поиска в реальном времени.

Подсветка конфликтов – цветное выделение дублирующихся комбинаций.

Сохранение настроек – запоминание последних параметров поиска.

📂 Состав репозитория
Язык	Файл	Статус
Python	shortcut_analyzer.py	✅
Go	shortcut_analyzer.go	✅
C++	shortcut_analyzer.cpp	✅
JavaScript	shortcut_analyzer.js	✅
C#	shortcut_analyzer.cs	✅
Ruby	shortcut_analyzer.rb	✅
Java	shortcut_analyzer.java	✅
🤝 Вклад в проект
Приветствуются улучшения:

Поддержка дополнительных форматов (YAML, XML).

Графический интерфейс.

Интеграция с системными реестрами Windows/macOS/Linux.

Создавайте Issues и Pull Requests.

📜 Лицензия
MIT License – свободное использование, модификация и распространение.

📂 Исходный код
Первая строка каждого файла – его имя. Скопируйте блок целиком и сохраните в соответствующий файл.

