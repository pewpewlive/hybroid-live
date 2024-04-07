
local test = require("/dynamic/import_test.lua")
local map = test.a
local one = 1
local function thing(param, param2)
	if param then
		return "yeah"
	end
	
end

