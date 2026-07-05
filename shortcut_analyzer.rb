#!/usr/bin/env ruby
# shortcut_analyzer.rb
# encoding: UTF-8

require 'json'
require 'csv'
require 'set'

COLORS = {
  reset: "\e[0m",
  green: "\e[92m",
  red: "\e[91m",
  yellow: "\e[93m",
  blue: "\e[94m",
  cyan: "\e[96m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

class Analyzer
  attr_reader :shortcuts, :apps, :categories, :conflicts

  def initialize(data)
    @shortcuts = data['shortcuts'] || []
    @apps = Hash.new { |h, k| h[k] = [] }
    @categories = Hash.new { |h, k| h[k] = [] }
    @conflicts = Hash.new { |h, k| h[k] = [] }
    index
  end

  def index
    @shortcuts.each do |s|
      app = s['app'] || 'Без приложения'
      @apps[app] << s
      cat = s['category'] || 'Без категории'
      @categories[cat] << s
      @conflicts[s['shortcut']] << s if s['shortcut'] && !s['shortcut'].empty?
    end
  end

  def search(query)
    q = query.downcase
    @shortcuts.select do |s|
      s['name'].downcase.include?(q) || s['shortcut'].downcase.include?(q)
    end
  end

  def get_conflicts
    @conflicts.select { |_, v| v.size > 1 }
  end

  def group_by_app
    @apps
  end

  def export_csv(filename)
    CSV.open(filename, 'w') do |csv|
      csv << ['name', 'shortcut', 'app', 'category']
      @shortcuts.each { |s| csv << [s['name'], s['shortcut'], s['app'], s['category']] }
    end
    puts colorize("Экспортировано в #{filename}", :green)
  end

  def export_json(filename)
    File.write(filename, JSON.pretty_generate({ 'shortcuts' => @shortcuts }))
    puts colorize("Экспортировано в #{filename}", :green)
  end

  def display(items = nil, group: false, conflicts_only: false, verbose: false)
    data = items || @shortcuts
    return puts colorize('Нет данных для отображения.', :yellow) if data.empty?
    if conflicts_only
      conf = get_conflicts
      if conf.empty?
        puts colorize('Конфликтов не найдено.', :green)
        return
      end
      puts colorize('🔍 Конфликты (одинаковые комбинации):', :bold)
      conf.each do |combo, list|
        puts colorize("  #{combo}:", :red)
        list.each { |s| puts "    - #{s['name']} (#{s['app']})" }
      end
      return
    end
    if group
      group_by_app.each do |app, items|
        puts colorize("\n📁 #{app} (#{items.size}):", :blue)
        items.each { |s| print_item(s, verbose) }
      end
    else
      data.each { |s| print_item(s, verbose) }
    end
  end

  def print_item(s, verbose)
    if verbose
      puts "  #{colorize(s['name'], :bold)} → #{colorize(s['shortcut'], :cyan)}  (приложение: #{s['app']}, категория: #{s['category']})"
    else
      puts "  #{s['name']} → #{s['shortcut']}"
    end
  end
end

def main
  file = nil
  search = nil
  export_file = nil
  group = false
  conflicts_only = false
  verbose = false

  i = 0
  while i < ARGV.size
    arg = ARGV[i]
    case arg
    when '-s' then search = ARGV[i+1]; i += 1
    when '-g' then group = true
    when '-c' then conflicts_only = true
    when '-e' then export_file = ARGV[i+1]; i += 1
    when '-v' then verbose = true
    when '-h', '--help'
      puts "Usage: ruby shortcut_analyzer.rb <file> [-s search] [-g] [-c] [-e file] [-v]"
      return
    else file = arg if file.nil?
    end
    i += 1
  end

  unless file
    puts colorize('Укажите файл конфигурации.', :red)
    exit 1
  end

  begin
    data = JSON.parse(File.read(file))
  rescue => e
    puts colorize("Ошибка загрузки файла: #{e.message}", :red)
    exit 1
  end

  analyzer = Analyzer.new(data)

  if search
    results = analyzer.search(search)
    if results.any?
      puts colorize("🔍 Найдено #{results.size} совпадений:", :bold)
      analyzer.display(results, group: group, verbose: verbose)
    else
      puts colorize('Ничего не найдено.', :yellow)
    end
    return
  end

  if conflicts_only
    analyzer.display(conflicts_only: true)
    return
  end

  if export_file
    if export_file.end_with?('.csv')
      analyzer.export_csv(export_file)
    elsif export_file.end_with?('.json')
      analyzer.export_json(export_file)
    else
      puts colorize('Неизвестный формат экспорта. Используйте .csv или .json', :red)
    end
    return
  end

  analyzer.display(group: group, verbose: verbose)
end

main if __FILE__ == $0
