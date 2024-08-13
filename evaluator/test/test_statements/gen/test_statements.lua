
local E_a = 1

if true then
	E_a = 2
end

if true then
	E_a = 2
end

while true do
	local E_a = 2

	::GL_::
end
if true then
	E_a = 2
end

local function E_test()
	E_a = 2
end
local function E_test2()
	if true then
		return 2
	end

end
local function E_test3()
	return 2, 3
end
local function E_test4(E_a)
	local H1
	if E_a == 2 then
		H1 = 2
		goto GL0

	else
		H1 = 1
		goto GL0

	end
	::GL0::
	return H1, 2
end
