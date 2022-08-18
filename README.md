# Monkey Programming Language Interpreter

This project is written by following the [Writing An Interpreter In Go](https://interpreterbook.com/) and [Writing A Compiler In Go](https://compilerbook.com/) books.

It demonstrates (Interpreter book):

- How to build an interpreter for a C-like programming language from scratch
- What a **lexer**, a **parser**, and an **Abstract Syntax Tree (AST)** are, and how to build your own
- What **closures** are and how and why they work
- What the **Pratt parsing technique** and a **recursive descent parser** is
- What others talk about when they talk about **built-in data structures**
- What **REPL** stands for and how to build one

It demonstrates (Compiler book):

- How to define our **bytecode instructions** and specify their operands and their encoding. Along the way, we also build a **mini-disassembler** for them
- How to write a **compiler** that takes in a **Monkey AST** and turns it into bytecode by emitting **instructions**
- At the same time, we build a **stack-based virtual machine** that executes the bytecode in its main loop
