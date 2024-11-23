
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
function E_HCTest_method3(Self)
	return Self
end
function E_HCTest_method2(Self)
	return Self[5]
end
function E_HCTest_New()
	local Self = {1,2fx,E_EnumTest[4],nil,{{{
		xd = E_HCTest_New()
	}}}}
	Self[4] = Self

	return Self
end

local E_p = E_HCTest_New()

local E_asdasdads = E_p[4]

local E_assddxxx = E_p

local E_enumsss = E_asdasdads[3]

local E_enumss = E_p[4][3]

local E_asdsss = E_p[4][0]

local E_bizh = E_HCTest_method(E_HCTest_method3(E_HCTest_method3(E_HCTest_method2(E_HCTest_method3(E_HCTest_method3(E_HCTest_method2(E_p)
[0][0]["xd"])
)
[4]["xd"])
[0][0]["xd"])
)
, {
	field1 = 1, field2 = true, field3 = E_EnumTest[1]

})


local function E_function(E_param)
	if E_param == E_EnumTest[1] then
	end

end
local E_sadfsadgf = nil

E_HEaasd = {}
function E_HEaasd_Spawn(E_x, E_y)
	local id = pewpew.new_customizable_entity(E_x, E_y)
	E_HEaasd[id] = {}
	return id
end
function E_HEaasd_Destroy(id)
end

local E_aasd12 = pewpew.new_customizable_entity(0fx, 0fx)

local E_adf = function(param0, param1) end

local E_a = 1

local E_b = {1, 2, 3}

local E_asdsdssdds = {2fx, 3fx, 4fx}

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

E_e = true

E_f = false

E_g = not true

local E_g2 = not not not not false

E_h = -1

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

local E_q = E_p.a

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

local E_aasdsad = E_EnumTest[1]

if E_aasdsad == E_EnumTest[1] then
end

