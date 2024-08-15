
local E_a = 1

local E_b = {1, 2, 3}

local E_c = {
	a = {
		b = function ()
			return 2
		end
	}
}

local E_d = {
	a = 1, 
	b = 2, 
	c = 3

}

local E_e = true

local E_f = false

local E_g = not true

local E_g2 = not not not not false

local E_h = -1

local E_i = function ()
end

local E_j = function ()
	return 2
end

local E_k = function (param)
	return false
end

local E_l = function ()
	return "a", {false, true}
end

local E_m = function ()
	return function ()
		return function ()
		end
	end
end

local E_n = {function (a, b)
	return false
end, function (a, b)
	return true
end}

local E_o = "string"

local E_EnumTest = {
	0, 
	1, 
	2, 
	3
}
function E_HCTest_method(Self, E_param1)
	return false
end
function E_HCTest_method1(Self, E_param1)
	return false
end
function E_HCTest_New()
	local Self = {1,2fx,E_EnumTest[4]}
	return Self
end

local E_p = E_HCTest_New()

local E_q = E_p[1]

local E_r = E_p[2]

local E_s = E_HCTest_method(E_p, {
	field1 = 1, 
	field2 = true, 
	field3 = E_EnumTest[1]

})


local E_t = E_b[1]

local E_u = E_c["a"]["b"]()

local E_v = {
	a = {{
		field = {E_HCTest_New(), E_HCTest_New()}, 
		field2 = 0.142fx
	
}}
}

