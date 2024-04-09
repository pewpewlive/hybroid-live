
local test = require("/dynamic/import_test.lua")
local map = {
	o = 1
}

local one = 1

local two = 2

one = one + (1)
two, map.o = 10, 20
local function thing(param, param2)
	if param then
	elseif param2 == "oooo" then
		return "a"
	else 
		param2 = "looool"
		return "a"
	end	
	param2 = "aa"
	if param2 == "000" and param then
		return "ppp"
	end	
	return param2
end
