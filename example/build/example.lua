








local test = require("/dynamic/import_test.lua")
local map = {
	o = 1
}

local one = 1

local two = 2

one = one + (1)
two, map.o = 10, 20
local function thing(param, param2)
	return "a"
end
