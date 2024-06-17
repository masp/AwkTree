# AwkTree (pronounced like octree)
AwkTree is [AWK](https://en.wikipedia.org/wiki/AWK) for Trees. You can parse, search, and transform any [Treesiter](https://tree-sitter.github.io/tree-sitter/) structured language source file with AWK like syntax:
```
// Print all identifiers in a program
(ident) {
    print(@)
}

// Rename a variable
(var_decl name: (ident) @n) {
    if n == "forbidden" {
        n = `correct`
    }
}
```

# Getting Started
Download from GitHub and install the executable. Make sure you have Go installed.
```sh
go install github.com/masp/awktree/cmd/tra
```

```sh
echo 'let a = 10; console.log(a)' > src.js
tra '(id){print(@)}' src.js
# You should see:
#  a
#  console
#  a
```

## About

With AWK editing CSV or tabular files is natural and easy to do in a few seconds, so why not combine the ease of AWK with the pattern matching of Treesitter? There are other tools for structured search, but none provide the power of a custom DSL (domain specific language) to make writing patterns that do exactly what you want.

## How it works

Treesitter produces concrete syntax trees. To navigate and match a syntax tree is different than records that AWK supports. AWK only supports "flat" trees, where you have a root node (program), followed by a list of records, where each record has 1-NF number of fields. The tree looks like:

$0
 | \
$1  $2

where $0 is the whole line, and $1, $2 are the fields. In a generic tree, numerals aren't enough. If we have a tree from treesitter like:

```
class MyClass {
    int field = 1;
}
```
it produces a tree like:
```
program
    class_declaration
        identifier: MyClass
        class_body
            field_declaration
                type: int
                identifier: field
```

which is more complex than the record/field delineation we had previously, which means we need a new syntax to describe them.

### Patterns
In AWK, the structure is predefined based on the -F flag and AWK's internal assumptions. The text and tree must be structured as a list of record, where each record (line) has 1 or more fields. The empty pattern matches any line:

```awk
{ print $0 }
```

In our case, our patterns describe trees. A simple pattern could be "find me every identifier in the Java program". To write patterns, we can use [Treesitter's queries](https://tree-sitter.github.io/tree-sitter/using-parsers#pattern-matching-with-queries) to describe what kind of nodes we want to find.

```
   Pattern      Action
| -------- | | -------- |
(identifier) { print(@) }
```

Any valid TreeSitter query can be used as a pattern.

TODO: Support patterns with example syntax, like:
```
// Print variable declarations where v is assigned null
`let @v@ = null` { print(@) }
```

## TODO: Editing
It's designed to be easy to use from the command line. Manipulating and editing nodes is as simple
as assigning different values.

```sh
# TODO: Edit main.go to replace interface{} with any
tra -i '`interface{}` { @ = "any" }' main.go
```