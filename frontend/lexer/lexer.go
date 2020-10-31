package lexer

import "github.com/alecthomas/participle/lexer"

var QueryLexer = lexer.Must(lexer.Regexp(`(\s+)` +
	`|(?P<Keyword>(?i)SELECT|FROM|WHERE|TRUE|FALSE|AS|AND|OR)` +
	`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
	`|(?P<Number>[-+]?\d*\.?\d+([eE][-+]?\d+)?)` +
	`|(?P<String>'[^']*'|"[^"]*")` +
	`|(?P<Operators><>|!=|<=|>=|[-+*/%,.()=<>])`,
))
