package main

import "log"

// ScriptErrors handler
type ScriptErrors struct {
	enable bool
}

// NewScriptErrors Creates a new ScriptErrors object
func NewScriptErrors() *ScriptErrors {
	var errors ScriptErrors
	errors.enable = true

	return &errors
}

// Enable Enables debug output.
func (s *ScriptErrors) Enable() {
	s.enable = true
}

// Disable Disables debug output.
func (s *ScriptErrors) Disable() {
	s.enable = false
}

// IsEnabled True if Debug Mode is enabled false otherwise
func (s *ScriptErrors) IsEnabled() bool {
	return s.enable
}

// Fatal Stops execution with an erro
func (s *ScriptErrors) Fatal(err ...interface{}) {
	if s.IsEnabled() {
		log.Fatal(err...)
	}
}
