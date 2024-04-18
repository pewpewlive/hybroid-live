





local function process(a)
	return !a, "yeah"
end
local function thing(param, param2)
	if param then
		return param2, 1
	end	
	param2 = "aa"
	if param2 == "000" and param then
		return "aaaa", 2
	end	
	return "b", 0
end
local ma = {1fx, 9fx, 4fx}

local p = ma[1] - 9fx

