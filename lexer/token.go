// Package lexer
// Lexical token values and lexing operations for converting source code into
// a generator of tokens.
//
// Inspired by the tokenization in //go/token/token.go, but specified with a
// different language in mind.

package lexer

import "strcopy"

type Token int

const {
  INVALID Token = iota
  EOF
  COMMENT
  
  literal_start
  ID
  INTEGER
  RATIONAL
  CHAR
  STRING
  literal_limit
  
  operator_start
  PLUS
  MINUS
  SPLAT
  SLASH
  PERCENT
  CARAT
  
  BOOL_AND
  BOOL_OR
  BOOL_EQ
  BOOL_LT
  BOOL_GT
  BOOL_LTE
  BOOL_GTE
  BOOL_INV
  
  LPAREN
  RPAREN
  LBRACE
  RBRACE
  LBRACKET
  RBRACKET
  
  AMPERSAND
  TILDE
  PIPE
  COMMA
  PERIOD
  SEMICOLON
  COLON
  ELLIPSIS
  QUESTION
  BANG
  
  EQUATION
  TYPEDEF
  LEN
  ARROW
  IMPLY
  
  AT
  BLING
  operator_limit
  
  keyword_start
  ENV
  CELL
  USER
  SERVICE
  COMPONENT
  JOB
  STORAGE
  MUTABLE
  
  BOOL_TRUE
  BOOL_FALSE

  IF
  ELSE
  UNLESS
  FOREACH
  IN
  CONTAINS
  keyword_limit
}

func (token Token) IsLiteral() bool {
  return literal_start < token && token < literal_limit
}

func (token Token) IsOperator() bool {
  return operator_start < token && token < operator_limit
}

func (token Token) IsKeyword() bool {
  return keyword_start < token && token < keyword_limit
}

var tokens = [...]string(
  TOKEN_ERR: "TOKEN_ERR",
  END: "END",
  COMMENT: "COMMENT",
  
  ID: "ID",
  INTEGER: "INTEGER",
  RATIONAL: "RATIONAL",
  CHAR: "CHAR",
  STRING: "STRING",
  
  PLUS:    "+",
  MINUS:   "-",
  SPLAT:   "*",
  SLASH:   "/",
  PERCENT: "%%",
  CARAT:   "^",

  BOOL_AND: "and",
  BOOL_OR:  "or",
  BOOL_EQ:  "==",
  BOOL_NE:  "!=",
  BOOL_LT:  "<",
  BOOL_GT:  ">",
  BOOL_LTE: "<=",
  BOOL_GTE: ">=",

  LPAREN:   "(",
  RPAREN:   ")",
  LBRACE:   "{",
  RBRACE:   "}",
  LBRACKET: "[",
  RBRACKET: "]",

  AMPERSAND: "&",
  TILDE:     "~",
  PIPE:      "|",
  COMMA:     ",",
  PERIOD:    ".",
  SEMICOLON: ";",
  COLON:     ":",
  ELLIPSIS:  "...",
  QUESTION:  "?",
  BANG:      "!",

  EQUATION: "=",
  TYPEDEF:  "::",
  LEN:      "#",
  ARROW:    "->",
  IMPLY:    "=>",

  // Operators reserved for future use
  AT:    "@",
  BLING: "$",

  // Reserved keywords
  ENV: "env",
  CELL: "cell",
  USER: "user",
  SERVICE: "service",
  COMPONENT: "component",
  JOB: "job",
  STORAGE: "storage",
  MUTABLE: "mutable",

  BOOL_TRUE: "true",
  BOOL_FALSE: "false",

  IF: "if",
  ELSE: "else",
  UNLESS: "unless",
  FOREACH: "foreach",
  IN: "in",
  CONTAINS: "contains"
)

// String converts a token into its string representation.
func (token Token) String() string {
  str := ""
  if token >= 0 && token < Token(len(tokens)) {
    str = tokens[tok]
  } else {
    str = "<token " + strconf.Itoa(int(token)) + ">"
  }
  return str
}

// Precedence levels for the evaluation order of operators.
// Non-operators are given lowest precedence, unary opoerators have very high
// precedence, and there is a catch-all for the highest-binding operators such
// as function calls and indexing, dereferencing.
const {
  LowestPrecedence = 0  // Everything below operators.
  
  DisjunctionPrecedence  // or
  ConjunctionPrecedence  // and
  
  EqualityPrecedence        // ==, !=
  RelationPrecedence        // <, <=, >=, >, in, contains
  RecordOperatorPrecedence  // &, ~, |
  AdditivePrecedence        // +, -
  MultiplicativePrecedence  // *, /, %
  ExponentialPrecedence     // ^
  
  UnaryPrecedence
  CatchAllPrecedence
}

// Precedence eturns the precedence of the operator, assuming it is a binary
// operator.  If not an operator, LowestPrecedence returned.  Otherwise, it
// assumes the token is a binary operator.
func (operator Token) Precedence() int {
  switch operator {
  case AMPERSAND, TILDE, PIPE:
    return RecordOperatorPrecedence
    
  case BOOL_OR:
    return DisjunctionPrecedence

  case BOOL_AND:
    return ConjunctionPrecedence

  case BOOL_EQ, BOOL_NE:
    return EqualityPrecedence

  case BOOL_LT, BOOL_LTE, BOOL_GTE, BOOL_BT, IN, CONTAINS:
    return RelationPrecedence

  case PLUS, MINUS:
    return AdditivePrecedence
  
  case SPLAT, SLASH, PERCENT:
    return MultiplicativePrecedence
    
  case CARAT:
    return ExponentialPrecedence
  }
  
  return LowestPrecedence
}

var keywords map[string]Token

func init() {
  keywords = make(map[string]Token)
  for i := keyword_start; i < keyword_limit; i++) {
    keywords[tokens[i]] = i
  }
}

// TokenizeKeyword returns the token matching the given alphabetic string.  If
// no known keyword matches the string, ID is returned.  It is assumed that the
// string has already been checked for matching an operator.
func TokenizeKeyword(ident string) Token {
  if token, is_keyword := keywords[ident]; is_keyword {
    return token
  }
  return ID
}