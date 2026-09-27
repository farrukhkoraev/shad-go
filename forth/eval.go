//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PLUS  string = "_plus_"
	MINUS string = "_minus_"
	MUL   string = "_mul_"
	DIV   string = "_div_"
	DUP   string = "_dup_"
	OVER  string = "_over_"
	DROP  string = "_drop_"
	SWAP  string = "_swap_"
)

type EvalError struct {
	msg string
}

func (e EvalError) Error() string {
	return e.msg
}

type IntStack []int

func (s IntStack) isEmpty() bool {
	return len(s) == 0
}

func (s IntStack) push(v int) IntStack {
	return append(s, v)
}

func (s IntStack) pop() (int, IntStack, error) {
	if s.isEmpty() {
		return -1, s, EvalError{"empty stack"}
	}
	return s[len(s)-1], s[0 : len(s)-1], nil
}

func (s IntStack) peek() (int, error) {
	if s.isEmpty() {
		return -1, EvalError{"empty stack"}
	}
	return s[len(s)-1], nil
}

type Evaluator struct {
	stack IntStack
	env   map[string][]string
}

func (e *Evaluator) expandWords(words []string) ([]string, error) {
	result := []string{}
	for _, word := range words {
		if _, err := strconv.Atoi(word); err != nil {
			// not a number, check if in env
			if w, ok := e.env[word]; ok {
				result = append(result, w...)
			} else {
				return []string{}, EvalError{"uknown instruction"}
			}
		} else {
			// number, add to result
			result = append(result, word)
		}
	}
	return result, nil
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		stack: make(IntStack, 0, 10),
		env: map[string][]string{
			"+":    []string{PLUS},
			"-":    []string{MINUS},
			"*":    []string{MUL},
			"/":    []string{DIV},
			"dup":  []string{DUP},
			"over": []string{OVER},
			"drop": []string{DROP},
			"swap": []string{SWAP},
		},
	}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	row = strings.ToLower(row)
	words := strings.Split(row, " ")

	if len(words) == 0 {
		return e.stack, nil
	}
	if words[0] == ":" {
		// : name word1 word2 ... ;
		if len(words) < 4 || words[len(words)-1] != ";" {
			return e.stack, EvalError{"invalid definition"}
		}
		// cannot redefine numbers
		if _, err := strconv.Atoi(words[1]); err == nil {
			return e.stack, EvalError{"cannot redefine numbers"}
		}
		expanded, err := e.expandWords(words[2 : len(words)-1])
		if err != nil {
			return e.stack, err
		}

		e.env[words[1]] = expanded
		return e.stack, nil
	}
	expanded, err := e.expandWords(words)
	if err != nil {
		return e.stack, err
	}

	return e.stack, e.eval(expanded)
}

func (e *Evaluator) eval(words []string) error {
	fmt.Printf("words: %v\n", words)

	for _, word := range words {
		// try number
		n, convErr := strconv.Atoi(word)
		if convErr == nil {
			e.stack = e.stack.push(n)
			continue
		}

		var fst int
		var err error
		var snd int
		fst, e.stack, err = e.stack.pop()
		if err != nil {
			return err
		}

		switch word {
		case DUP:
			e.stack = e.stack.push(fst)
			e.stack = e.stack.push(fst)
		case DROP:
			// fst already popped
		case OVER:
			snd, err = e.stack.peek()
			if err != nil {
				return err
			}
			e.stack = e.stack.push(fst)
			e.stack = e.stack.push(snd)
		case SWAP:
			snd, e.stack, err = e.stack.pop()
			if err != nil {
				return err
			}
			e.stack = e.stack.push(fst)
			e.stack = e.stack.push(snd)
		default:
			snd, e.stack, err = e.stack.pop()
			if err != nil {
				return err
			}
			switch word {
			case PLUS:
				e.stack = e.stack.push(snd + fst)
			case MINUS:
				e.stack = e.stack.push(snd - fst)
			case MUL:
				e.stack = e.stack.push(snd * fst)
			case DIV:
				if fst == 0 {
					return EvalError{"division by zero"}
				}
				e.stack = e.stack.push(snd / fst)
			default:
				return EvalError{"unknown instruction"}
			}
		}
	}
	return nil
}
