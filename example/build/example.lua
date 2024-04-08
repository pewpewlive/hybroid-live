
local test = require("/dynamic/import_test.lua")
local map = test.a

local one = 1

local function thing(param, param2)
	if param then
	elseif param2 == "oooo" then
		return "a"
	else 
		param2 = "looool"
		return "a"
	end	
	param2 = "aa"
	return param2
end
