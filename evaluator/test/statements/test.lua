
local E_a = false
local function E_thing()
	return E_a
end
if E_thing() then
elseif E_thing() then
else 
end
local function E_thing2(E_data)
	local E_a = E_data["numbers"]
	local E_b = function()
		if #E_a < 3 then
			return 1
		else 
			return 2
		end
	end
	for E_i, E_v in ipairs(E_a) do
		if E_a[E_i] >= 6 and E_v < 9 then
			local H1 = E_data["booleans"][1]
			if H1 == true then
				return 9, E_b
			elseif H1 == false then
				return 0, E_b
			else
				goto GL_
			end
			::GL0::
		end
		::GL_::
	end
	return 0, function()
		return -1
	end
end
local E_data = {
	numbers = {1, 2}, 
	booleans = {false}
}
E_thing2(E_data)
