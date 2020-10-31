package query_lexer

import "github.com/alecthomas/participle/lexer"

var queryLexer = lexer.Must(lexer.Regexp(`(\s+)` +
	`|(?P<Keyword>(?i)SELECT|FROM|WHERE|TRUE|FALSE|AS)` +
	`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
	`|(?P<Number>[-+]?\d*\.?\d+([eE][-+]?\d+)?)` +
	`|(?P<String>'[^']*'|"[^"]*")` +
	`|(?P<Operators><>|!=|<=|>=|[-+*/%,.()=<>])`,
))
