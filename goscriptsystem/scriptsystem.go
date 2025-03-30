package goscriptsystem

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// ScriptSystem use Lua for scripting.
type ScriptSystem struct {
	state  *lua.LState
	errors ScriptErrors
}

// New ScriptSystem Initializes the Lua Script System
func New(errors ScriptErrors) *ScriptSystem {
	var scriptSystem ScriptSystem

	scriptSystem.state = lua.NewState()
	scriptSystem.errors = errors

	return &scriptSystem
}

func (s *ScriptSystem) GetState() *lua.LState {
	return s.state
}

// CallFunc Call a Lua function
func (s *ScriptSystem) CallFunc(funcName string, numReturnValues int, returnError bool, args ...lua.LValue) (lua.LValue, error) {
	luaFunc := lua.P{
		Fn:      s.state.GetGlobal(funcName),
		NRet:    numReturnValues,
		Protect: returnError,
	}

	var returnVal lua.LValue

	err := s.state.CallByParam(luaFunc, args...)

	if numReturnValues == 1 {
		returnVal = s.state.Get(-1)
		s.state.Pop(1)
	}

	if err != nil {
		return returnVal, err
	}

	return returnVal, nil
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
func (s *ScriptSystem) CallFuncSimple(funcName string, args ...lua.LValue) error {
	_, err := s.CallFunc(funcName, 0, true, args...)

	if err != nil {
		return err
	}

	return nil
}

// CallFuncWithReturn Call a Lua function that has one return value
func (s *ScriptSystem) CallFuncWithReturn(funcName string, args ...lua.LValue) (lua.LValue, error) {
	value, err := s.CallFunc(funcName, 1, true, args...)

	if err != nil {
		return value, err
	}

	return value, nil
}

func (s *ScriptSystem) OnCreate(errOnNotFound bool) error {
	if s.HasFunc("OnCreate") {
		err := s.CallFuncSimple("OnCreate")

		if errOnNotFound && err != nil {
			return errors.New("failed to call OnCreate function. OnCreate function not found")
		}
	}

	return nil
}

func (s *ScriptSystem) OnDestroy(errOnNotFound bool) error {
	if s.HasFunc("OnDestroy") {
		err := s.CallFuncSimple("OnDestroy")

		if errOnNotFound && err != nil {
			return errors.New("failed to call OnDestroy function. OnDestroy function not found")
		}
	}

	return nil
}

// SetGlobal Just like the Lua version.
func (s *ScriptSystem) SetGlobal(name string, value interface{}) {
	s.state.SetGlobal(name, luar.New(s.state, value))
}

// GetString Get a string value from the Lua script variable
func (s *ScriptSystem) GetString(name string, defaultValue string) string {
	if s.HasFunc(name) {
		return string(s.state.GetGlobal(name).(lua.LString))
	}

	return defaultValue
}

// GetNumber Get a float64 value from the Lua script variable
func (s *ScriptSystem) GetNumber(name string, defaultValue float64) float64 {
	if s.HasFunc(name) {
		return float64(s.state.GetGlobal(name).(lua.LNumber))
	}

	return defaultValue
}

// GetBool Get a boolean value from the Lua script variable
func (s *ScriptSystem) GetBool(name string, defaultValue bool) bool {
	if s.HasFunc(name) {
		return bool(s.state.GetGlobal(name).(lua.LBool))
	}

	return defaultValue
}

// DestroyScriptSystem Calls lua.LState.Close
func (s *ScriptSystem) DestroyScriptSystem() {
	s.state.Close()
}

// NewTable Creates a new table
func (s *ScriptSystem) NewTable() *lua.LTable {
	return s.state.NewTable()
}

// RegisterFunction Register a function. This is here for convenience. SetGlobal is more flexible and should be preferred
func (s *ScriptSystem) RegisterFunction(name string, fn lua.LGFunction) {
	s.state.Register(name, fn)
}

// DoString Run the passed code string
func (s *ScriptSystem) DoString(code string) {
	err := s.state.DoString(code)

	if err != nil {
		s.errors.HandleError(err)
	}
}

// DoFile Load the file and run its code
func (s *ScriptSystem) DoFile(fileName string) {
	err := s.state.DoFile(fileName)

	if err != nil {
		s.errors.HandleError(err)
	}
}

// DoFiles Loads and processes files from the list generated by os.ReadDir
func (s *ScriptSystem) DoFiles(dirName string) {
	commandFiles, err := filepath.Glob(path.Join(dirName, "*.lua"))

	if err != nil {
		s.errors.HandleError(err)
	}

	for _, fileName := range commandFiles {
		s.DoFile(fileName)
	}
}

// LoadString load the passed code string
func (s *ScriptSystem) LoadString(code string) (*lua.LFunction, error) {
	luaFunc, err := s.state.LoadString(code)

	if err != nil {
		s.errors.HandleError(err)
	}

	return luaFunc, err
}

// LoadFile Load the file
func (s *ScriptSystem) LoadFile(fileName string) (*lua.LFunction, error) {
	luaFunc, err := s.state.LoadFile(fileName)

	if err != nil {
		s.errors.HandleError(err)
	}

	return luaFunc, err
}

// LoadFiles Loads files from the list generated by os.ReadDir
func (s *ScriptSystem) LoadFiles(dirName string) []*lua.LFunction {
	var funcList []*lua.LFunction

	files, readDirErr := os.ReadDir(dirName)

	if readDirErr != nil {
		s.errors.HandleError(readDirErr)
	}

	for _, file := range files {
		luaFunc, loadErr := s.LoadFile(file.Name())

		if loadErr != nil {
			s.errors.HandleError(loadErr)
		} else {
			funcList = append(funcList, luaFunc)
		}
	}

	return funcList
}
