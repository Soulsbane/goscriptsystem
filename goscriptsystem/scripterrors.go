package goscriptsystem

import (
	"fmt"
	"log"
)

type ScriptErrors interface {
	Enable()
	Disable()
	IsEnabled() bool
	Fatal(err ...any)
	Println(err ...any)
}

// StdOutScriptErrors handler
type StdOutScriptErrors struct {
	enable bool
}

// NewScriptErrors Creates a new StdOutScriptErrors object
func NewStdOutScriptErrors() *StdOutScriptErrors {
	var errors StdOutScriptErrors
	errors.enable = true

	return &errors
}

// Enable Enables debug output.
func (s *StdOutScriptErrors) Enable() {
	s.enable = true
}

// Disable Disables debug output.
func (s *StdOutScriptErrors) Disable() {
	s.enable = false
}

// IsEnabled True if Debug Mode is enabled false otherwise
func (s *StdOutScriptErrors) IsEnabled() bool {
	return s.enable
}

// Fatal Stops execution with an error
func (s *StdOutScriptErrors) Fatal(err ...any) {
	if s.IsEnabled() {
		log.Fatal(err...)
	}
}

func (s *StdOutScriptErrors) Println(err ...any) {
	if s.IsEnabled() {
		fmt.Println(err...)
	}
}
