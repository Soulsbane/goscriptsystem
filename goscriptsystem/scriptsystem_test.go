package goscriptsystem

import (
	"fmt"
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
	scriptSystem.DoString(`function exampleReturnFunc() return true end`)

	value, err := scriptSystem.CallFuncWithReturn("exampleReturnFunc")

	if err != nil {
		t.Error("Failed to call lua function: exampleReturnFunc")
	} else {
		if value.(lua.LBool) != true {
			t.Error("Return value is not true: ", value)
		}
	}
}

func TestRegisterFunc(t *testing.T) {
	scriptSystem := New(NewScriptErrors())

	scriptSystem.RegisterFunction("testFunc", func(L *lua.LState) int {
		fmt.Println("Hello world from TestRegisterFunc")
		return 0
	})

	scriptSystem.RegisterFunction("testFuncWithReturn", func(L *lua.LState) int {
		fmt.Print("Hello world from TestRegisterFuncWithReturn(bool): ")
		L.Push(lua.LBool(true))

		return 1
	})

	scriptSystem.DoString(`testFunc() print(testFuncWithReturn())`)
}
