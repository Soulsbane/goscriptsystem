package goscriptsystem

import (
	"fmt"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

type testGlobalStruct struct {
	Name string
}

func (t testGlobalStruct) TestStructFunc() {
	fmt.Println("TestStruct:TestStructFunc()")
}

func TestSetGlobal(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	testStruct := testGlobalStruct{Name: "This is the name of the test struct"}

	scriptSystem.SetGlobal("testGlobal", testSetGlobal)
	scriptSystem.SetGlobal("TestStruct", testStruct)
	scriptSystem.DoString(`function testGlobal() print(TestStruct.Name) TestStruct:TestStructFunc() end`)
}

func testSetGlobal() {
	fmt.Println("Hello world from testGlobal()")
}

func TestGetString(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	scriptSystem.DoString(`Author = "Soulsbane"`)

	testStringValue := scriptSystem.GetString("Author", "Neo")
	testDefaultValue := scriptSystem.GetString("NonExistent", "Neo")

	if testStringValue != "Soulsbane" {
		t.Errorf("Failed to get global string variable got: %q want: Soulsbane", testStringValue)
	}

	if testDefaultValue != "Neo" {
		t.Errorf("Failed to get default value for non-existent global variable got: %q< want: Neo", testDefaultValue)
	}
}

func TestGetNumber(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	scriptSystem.DoString(`Version = 100`)

	testNumberValue := scriptSystem.GetNumber("Version", 0)
	testDefaultValue := scriptSystem.GetNumber("NonExistent", 255)

	if testNumberValue != 100 {
		t.Errorf("Failed to get global number variable got: %v want: 100", testNumberValue)
	}

	if testDefaultValue != 255 {
		t.Errorf("Failed to get default value for non-existent global variable got: %v want: 0", testDefaultValue)
	}
}

func TestGetBool(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	scriptSystem.DoString(`TestBool = true`)

	testBoolValue := scriptSystem.GetBool("TestBool", false)
	testDefaultValue := scriptSystem.GetBool("NonExistent", false)

	if testBoolValue != true {
		t.Errorf("Failed to get global bool variable got: %v want: true", testBoolValue)
	}

	if testDefaultValue != false {
		t.Errorf("Failed to get default value for non-existent global variable got: %v want: false", testDefaultValue)
	}
}

func TestSimpleFuncCall(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	scriptSystem.DoString(`function testFunc() print("Hello world from testFunc()") end`)

	err := scriptSystem.CallFuncSimple("testFunc")

	if err != nil {
		t.Error("Failed to call lua function: testFunc")
	}
}

func TestFuncWithReturn(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	scriptSystem.DoString(`function exampleReturnFunc() return true end`)

	value, err := scriptSystem.CallFuncWithReturn("exampleReturnFunc")

	if err != nil {
		t.Error("Failed to call lua function: exampleReturnFunc")
	} else if value.(lua.LBool) != true {
		t.Error("Return value is not true: ", value)
	}
}

func TestRegisterFunc(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())

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

func TestLoadString(t *testing.T) {
	scriptSystem := New(NewStdOutScriptErrors())
	luaFunc, err := scriptSystem.LoadString(`print("Hello world from TestLoadString()")`)

	if err != nil {
		t.Error("Failed to load string: ", err)
	}

	scriptSystem.GetState().Push(luaFunc)
	err = scriptSystem.GetState().PCall(0, 0, nil)

	if err != nil {
		t.Error("Failed to call lua function: ", err)
	}
}
