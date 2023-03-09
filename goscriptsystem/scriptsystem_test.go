package goscriptsystem

import (
	"testing"
)

func TestSimpleFuncCall(t *testing.T) {
	scriptSystem := New(NewScriptErrors())
	scriptSystem.DoString(`function testFunc() print("Hello world from testFunc()") end`)

	// err := scriptSystem.CallFuncSimple("testFuncs")

	// if err != nil {
	// 	t.Error("Failed to call lua function: testFuncs")
	// }

	err := scriptSystem.CallFuncSimple("testFunc")

	if err != nil {
		t.Error("Failed to call lua function: testFunc")
	}
}
