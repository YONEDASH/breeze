class Entry:
    def __init__(self, name, _type):
        self.name = name
        self._type = _type

    def entry_type(self):
        return self._type

    def entry_name(self):
        return self.name


class Node:
    def __init__(self, name, entries):
        self.name = name
        self.entries = entries

        self.token_entry = 0
        for entry in entries:
            if entry.entry_type() == "scanner.Token":
                self.token_entry = entry
                break

    def is_expr(self):
        return isinstance(self, Expr)

    def is_stmt(self):
        return isinstance(self, Stmt)

    def is_decl(self):
        return isinstance(self, Decl)

    def node_name(self):
        if self.is_expr():
            return self.name + "Expr"
        elif self.is_stmt():
            return self.name + "Stmt"
        elif self.is_decl():
            return self.name + "Decl"
        return self.name + "Node"

    def node_id(self):
        return self.name + "Id"

    def token(self):
        return self.token_entry


class Err(Node):
    pass


class Expr(Node):
    pass


class Stmt(Node):
    pass


class Decl(Node):
    pass


def gen_struct(node):
    struct = "type " + node.node_name() + " struct {\n\tNode\n"

    if node.token() == 0:
        struct += "\tToken scanner.Token\n"

    for entry in node.entries:
        struct += "\t" + entry.entry_name() + " " + entry.entry_type() + "\n"

    struct += "}\n"
    return struct


def gen_functions(node):
    node_type = "Err"
    if node.is_expr():
        node_type = "Expr"
    elif node.is_stmt():
        node_type = "Stmt"
    elif node.is_decl():
        node_type = "Decl"

    get_type = "func (node *" + node.node_name() + ") GetType() NodeType {\n\treturn " + node_type + "\n}\n"

    get_id = "func (node *" + node.node_name() + ") GetId() NodeId {\n\treturn " + node.node_id() + "\n}\n"

    str_prefix = ""
    for entry in node.entries:
        if entry.entry_type()[0] == "[":
            var_name = "str_" + entry.entry_name()
            str_prefix += "\t" + var_name + " := \"{\""
            str_prefix += """
\tfor i, n := range node.""" + entry.entry_name() + """ {
\t\t""" + "str_" + entry.entry_name() + """ += fmt.Sprintf("%s", n)
\t\tif i <= len(node.""" + entry.entry_name() + """)-1 {
\t\t\t""" + "str_" + entry.entry_name() + """ += ", "
\t\t}
\t}"""
            str_prefix += "\n\t" + var_name + " += \"}\"\n"

    stringify = "func (node *" + node.node_name() + ") String() string {\n" + str_prefix + "\treturn \"(" + node.node_name()
    for entry in node.entries:
        stringify += " " + entry.entry_name() + "=\"+"
        if entry.entry_type() == "scanner.Token" or entry.entry_type() == "Node":
            stringify += "fmt.Sprintf(\"%s\", node." + entry.entry_name() + ")"
        elif entry.entry_type()[0] == "[":
            stringify += "str_" + entry.entry_name()
        else:
            stringify += "string("
            if entry.entry_type()[0] == "*":
                stringify += "*"
            stringify += "node." + entry.entry_name() + ")"
        stringify += "+\""
    stringify += ")\"\n}\n"

    get_token = "func (node *" + node.node_name() + ") GetToken() scanner.Token {\n\treturn "
    entry_token = node.token()
    if entry_token == 0:
        get_token += "node.Token"
    else:
        name = entry_token.entry_name()
        get_token += "node." + name
    get_token += "\n}\n"

    visit = ("func (node *" + node.node_name() + ") Visit(visitor Visitor) any {\n\treturn visitor.Visit" + node.node_name() + "(node)\n}\n")

    return get_type + "\n" + get_id + "\n" + stringify + "\n" + get_token + "\n" + visit


def gen_source(ast_nodes):
    source_code = """package ast

import (
\t"breeze/scanner"
\t"fmt"
)

type NodeId uint8

const (
"""

    for i, node in enumerate(ast_nodes):
        source_code += "\t" + node.node_id()
        if i == 0:
            source_code += " NodeId = iota"
        source_code += "\n"

    source_code += """)

type NodeType uint8

const (
\tErr NodeType = iota
\tExpr
\tStmt
\tDecl
)

type Node interface {
\tGetId() NodeId
\tGetType() NodeType
\tString() string
\tGetToken() scanner.Token
\tVisit(visitor Visitor) any
}

"""

    source_code += "type Visitor interface {\n"

    for node in ast_nodes:
        source_code += "\tVisit" + node.node_name() + "(node *" + node.node_name() + ") any\n"

    source_code += "}\n\n"

    for i, node in enumerate(ast_nodes):
        source_code += gen_struct(node) + "\n"
        source_code += gen_functions(node)
        if i != ast_nodes.__len__() - 1:
            source_code += "\n"

    return source_code


# AST Nodes
nodes = {
    Err("Err", {Entry("Message", "string"), Entry("Hint", "string")}),
    Decl("Let", {Entry("Identifier", "string"), Entry("Type", "string")}),
    Decl("Function", {
        Entry("Identifier", "string"), Entry("Closure", "Node"),
        Entry("ReturnType", "string"), Entry("ParamType", "[]string"),
        Entry("ParamName", "[]string")
    }),
    Decl("Struct", {
        Entry("Identifier", "string"),
        Entry("ParentType", "string"),
        Entry("ParamType", "[]string"),
        Entry("ParamName", "[]string")
    }),
    Stmt("Debug", {Entry("Expression", "Node")}),
    Stmt("Return", {Entry("Expression", "Node")}),
    Stmt("Continue", {}),
    Stmt("Break", {}),
    Stmt("Block", {Entry("Nodes", "[]Node")}),
    Stmt("Conditional", {Entry("Condition", "Node"), Entry("Statement", "Node"), Entry("ElseStatement", "Node")}),
    Stmt("While", {Entry("Condition", "Node"), Entry("Statement", "Node")}),
    Stmt("Closure", {Entry("Block", "Node")}),
    Stmt("Expr", {Entry("Expression", "Node")}),
    Expr("Assign", {Entry("Operator", "scanner.Token"), Entry("Name", "scanner.Token"), Entry("Value", "Node")}),
    Expr("Set", {Entry("Expression", "Node"), Entry("Name", "scanner.Token"), Entry("Value", "Node")}),
    Expr("Binary", {Entry("Operator", "scanner.Token"), Entry("Left", "Node"), Entry("Right", "Node")}),
    Expr("Unary", {Entry("Operator", "scanner.Token"), Entry("Expression", "Node")}),
    Expr("Call", {Entry("Expression", "Node"), Entry("Arguments", "[]Node")}),
    Expr("Get", {Entry("Expression", "Node"), Entry("Name", "scanner.Token")}),
    Expr("IdentifierLit", {Entry("Name", "string")}),
    Expr("IntegerLit", {Entry("Value", "string")}),
    Expr("FloatingLit", {Entry("Value", "string")}),
    Expr("BooleanLit", {Entry("Value", "string")}),
}

source = gen_source(nodes)

print(source)

file = open("ast/nodes.go", "w")
file.write(source)
file.close()