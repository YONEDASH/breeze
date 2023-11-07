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

    def node_name(self):
        if self.is_expr():
            return self.name + "Expr"
        elif self.is_stmt():
            return self.name + "Stmt"
        return self.name + "Decl"

    def node_id(self):
        return self.name + "Id"

    def token(self):
        return self.token_entry


class Expr(Node):
    pass


class Stmt(Node):
    pass


class Decl(Node):
    pass


def gen_struct(node):
    struct = "type " + node.node_name() + " struct {\n\tNode\n"

    if node.token() == 0:
        struct += "\ttoken scanner.Token\n"

    for entry in node.entries:
        struct += "\t" + entry.entry_name() + " " + entry.entry_type() + "\n"

    struct += "}\n"
    return struct


def gen_functions(node):
    node_type = "Decl"
    if node.is_expr():
        node_type = "Expr"
    elif node.is_stmt():
        node_type = "Stmt"

    get_type = "func (node *" + node.node_name() + ") GetType() NodeType {\n\treturn " + node_type + "\n}\n"

    get_id = "func (node *" + node.node_name() + ") GetId() NodeId {\n\treturn " + node.node_id() + "\n}\n"

    stringify = "func (node *" + node.node_name() + ") Stringify() string {\n\treturn \"(" + node.node_name()
    for entry in node.entries:
        stringify += " " + entry.entry_name() + "=\"+"
        if entry.entry_type() == "scanner.Token" or entry.entry_type() == "Expr" or entry.entry_type() == "Stmt":
            stringify += "node." + entry.entry_name() + ".Stringify()"
        else:
            stringify += "string(node." + entry.entry_name() + ")"
        stringify += "+\""
    stringify += ")\"\n}\n"

    get_token = "func (node *" + node.node_name() + ") GetToken() scanner.Token {\n\treturn "
    entry_token = node.token()
    if entry_token == 0:
        get_token += "node.token"
    else:
        name = entry_token.entry_name()
        get_token += "node." + name
    get_token += "\n}\n"

    return get_type + "\n" + get_id + "\n" + stringify + "\n" + get_token


def gen_source(ast_nodes):
    source_code = """package ast

import "breeze/scanner"

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
\tExpr NodeType = iota
\tStmt
\tDecl
)

type Node interface {
\tGetId() NodeId
\tGetType() NodeType
\tStringify() string
\tGetToken() scanner.Token
}

"""

    for i, node in enumerate(ast_nodes):
        source_code += gen_struct(node) + "\n"
        source_code += gen_functions(node)
        if i != ast_nodes.__len__() - 1:
            source_code += "\n"

    return source_code


# AST Nodes
nodes = {
    Expr("Binary", {Entry("Operator", "scanner.Token")}),
    Expr("Identifier", {Entry("Name", "string")}),
    Expr("Integer", {Entry("Value", "string")}),
    Expr("Floating", {Entry("Value", "string")}),
}

source = gen_source(nodes)

print(source)

file = open("ast/nodes.go", "w")
file.write(source)
file.close()