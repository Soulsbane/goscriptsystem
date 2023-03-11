package goscriptsystem

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSimpleFuncCall(t *testing.T) {
	scriptSystem := New(NewScriptErrors())
	scriptSystem.DoString(`function testFunc() print("Hello world from testFunc()") end`)

	err := scriptSystem.CallFuncSimple("testFunc")

	if err != nil {
		t.Error("Failed to call lua function: testFunc")
	}
}

func TestFuncWithReturn(t *testing.T) {
	scriptSystem := New(NewScriptErrors())
	scriptSystem.DoString(`function exampleReturnFunc() return false end`)

	value, err := scriptSystem.CallFuncWithReturn("exampleReturnFunc")

	if err != nil {
		t.Error("Failed to call lua function: exampleReturnFunc")
	} else {
		if value.(lua.LBool) != true {
			t.Error("Return value is not true: ", value)
		}
	}
}
