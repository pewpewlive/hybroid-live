




local function thing(param, param2)
	if param then
		return param2, 1
	end	
	param2 = "aa"
	if (param2 == "000") and param then
		return "aaaa", 2
	end	
	return "a", 3
end
local b = thing(true, "xd")

b = 2
