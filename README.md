# Amazing-Automata
A universal pipeline generator for Github Actions

# Table of contents
1. [Overview](#overview)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Usage](#usage)
5. [License](#license)

## Overview
Amazing-Automata is a universal CI/CD pipeline generator for GitHub Actions that empowers DevOps engineers to create, customize, and maintain workflows in seconds.

## Requirements
 - Go 1.24.6

## Installation

## Usage
### Options
1. Check a documentation in shell
```bash
amazing-automata -h
```
2. Generate only a CI or CD pipeline
```bash
amazing-automata --ci -o <filename>.yml
```
```bash
amazing-automata --cd  -o <filename.yml> 
```
3. Dry-run review
```bash
amazing-automata --dry-run --ci 
```
## License
This project is licensed under the MIT License. 