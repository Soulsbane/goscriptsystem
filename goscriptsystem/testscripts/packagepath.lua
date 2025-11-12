local args = { ... }
local testMod = require "testmodule"

assert(testMod.add(args[1], 1) == 3)
