# Amazing-Automata
Универсальный генератор пайплайнов для Github actions

# Содержание
1. [Обзор](#обзор)
2. [Зависимости](#зависимости)
3. [Установка](#установка)
4. [Использование](#использование)
5. [Лицензия](#лицензия)

## Обзор
Amazing-Automata это универсальная утилита которая позваляет DevOps инжинерам автоматические генерировать CI/CD Github Actions пайпланы а так же модифицировать существующие

## Зависимости
- Go 1.24.6

## Установка

## Использование
### Флаги
1. Просмотр документации в консоли
```bash
amazing-automata -h
```
2. Создание простого CI/CD пайплайна
```bash
amazing-automata <filename>.yml
```
3. Генерация только CI или CD пайпайна
```bash
amazing-automata <filename>.yml --ci
```
```bash
amazing-automata <name.yml> --cd
```
4. Изменение существующего пайплайна
```bash
amazing-automata --append <path/to/workflow.yml>
```
5. Просмотр в stdout изменений пайплана
```bash
amazing-automata --dry-run 
```
6. Использование Matrix
```bash
amazing-automata my-workflow.yml --matrix go=1.24,1.25 os=ubuntu-latest,macos-latest
```
## Лицензия
Этот проект лицензируются MIT