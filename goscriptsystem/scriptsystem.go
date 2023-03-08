package goscriptsystem

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// ScriptSystem use Lua for scripting.
type ScriptSystem struct {
	state  *lua.LState
	errors *ScriptErrors
}

// NewScriptSystem Initializes the Lua Script System
func New(errors *ScriptErrors) *ScriptSystem {
	var scriptSystem ScriptSystem

	scriptSystem.state = lua.NewState()
	scriptSystem.errors = errors

	return &scriptSystem
}

// CallFunc Call a Lua function
func (s *ScriptSystem) CallFunc(funcName string, numReturnValues int, returnError bool, args ...lua.LValue) lua.LValue {
	luaFunc := lua.P{
		Fn:      s.state.GetGlobal(funcName),
		NRet:    numReturnValues,
		Protect: returnError,
	}

	var returnVal lua.LValue

	/*err := s.state.CallByParam(luaFunc, args...)

	if err != nil {
		fmt.Println("Function name", funcName, "not found")
	}*/

	s.state.CallByParam(luaFunc, args...)

	if numReturnValues == 1 {
		returnVal = s.state.Get(-1)
		s.state.Pop(1)
	}

	return returnVal
}

func (s *ScriptSystem) HasFunc(funcName string) bool {
	exists := s.state.GetGlobal(funcName)

	if exists == lua.LNil {
		return false
	} else {
		return true
	}
}

// CallFuncSimple This is just sugar for calling a Lua function without having to deal with additional parameters.
func (s *ScriptSystem) CallFuncSimple(funcName string, args ...lua.LValue) {
	s.CallFunc(funcName, 0, true)
}

// CallFuncWithReturn Call a Lua function that has one return value
func (s *ScriptSystem) CallFuncWithReturn(funcName string, args ...lua.LValue) lua.LValue {
	return s.CallFunc(funcName, 1, true)
}

func (s *ScriptSystem) onCreate() {
	s.CallFuncSimple("OnCreate")
	//returnVal := s.CallFuncWithReturn("TestArgFunc", lua.LNumber(10), lua.LString("hello world"))
}

func (s *ScriptSystem) onDestroy() {
	s.CallFuncSimple("OnDestroy")
}

// SetGlobal Just like the Lua version.
func (s *ScriptSystem) SetGlobal(name string, value interface{}) {
	s.state.SetGlobal(name, luar.New(s.state, value))
}

// DestroyScriptSystem Calls lua.LState.Close
func (s *ScriptSystem) DestroyScriptSystem() {
	s.onDestroy()
	s.state.Close()
}

// NewTable Creates a new table
func (s *ScriptSystem) NewTable() *lua.LTable {
	return s.state.NewTable()
}

// DoString Run the passed code string
func (s *ScriptSystem) DoString(code string) {
	s.state.DoString(code)
}

// DoFile Load the file and run its code
func (s *ScriptSystem) DoFile(fileName string, callOnCreate bool) {
	err := s.state.DoFile(fileName)

	if err != nil {
		s.errors.Fatal(err)
	}

	if callOnCreate {
		s.onCreate()
	}
}

// DoFiles Loads and processes files from the list generated by ioutil.ReadDir
func (s *ScriptSystem) DoFiles(dirName string, callOnCreate bool) {
	commandFiles, err := filepath.Glob(path.Join(dirName, "*.lua"))

	if err != nil {
		log.Fatal(err)
	}

	for _, fileName := range commandFiles {
		s.DoFile(fileName, callOnCreate)
	}
}

// LoadString load the passed code string
func (s *ScriptSystem) LoadString(code string) (*lua.LFunction, error) {
	return s.state.LoadString(code)
}

// LoadFile Load the file
func (s *ScriptSystem) LoadFile(fileName string) (*lua.LFunction, error) {
	luaFunc, err := s.state.LoadFile(fileName)

	if err != nil {
		s.errors.Fatal(err)
	}

	return luaFunc, err
}

// LoadFiles Loads files from the list generated by ioutil.ReadDir
func (s *ScriptSystem) LoadFiles(dirName string) {
	// INFO: No point in this function at this time. Since we can't return for each file.
	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		s.LoadFile(file.Name())
	}
}
