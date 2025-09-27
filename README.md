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
2. Create the simple CI/CD pipeline
```bash
amazing-automata <filename>.yml
```
3. Generate only a CI or CD pipeline
```bash
amazing-automata <filename>.yml --ci
```
```bash
amazing-automata <name.yml> --cd
```
4. Edit an existing pipeline
```bash
amazing-automata --append <path/to/workflow.yml>
```
5. Dry-run review
```bash
amazing-automata --dry-run 
```
6. Matrix usage
```bash
amazing-automata my-workflow.yml --matrix go=1.24,1.25 os=ubuntu-latest,macos-latest
```
## License
This project is licensed under the MIT License. 