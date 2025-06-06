
local E_a = 1fx

local E_b = 0.409fx



E_d = 0.7fx

local E_e, E_f, E_g = function ()

	return 2, 3
end, 80247, 4294945535

local E_h, E_i = E_e()
local E_j = true

HEE_Entity = {}
function HEE_Entity_Spawn(E_x, E_y)
	local id = pewpew.new_customizable_entity(E_x, E_y)
	HEE_Entity[id] = {}
	local Self = HEE_Entity[id]
	Self[1] = function () end
	return id
end
function HEE_Entity_Destroy(id)
	local Self = HEE_Entity[id]
end
function HEE_Entity_method1(id)
	local Self = HEE_Entity[id]
end

local E_k = HEE_Entity_Spawn(0fx, 0fx)

HEE_Entity[E_k][1]()

local E_mp = {HEE_Entity[E_k][1], function () end}

E_mp[1]()

local E_l = {E_k, HEE_Entity_Spawn(200fx, 200fx)}

HEE_Entity[E_l[2]][1]()

local E_po = {
	["thing"] = E_l
}

HEE_Entity[E_po["thing"][2]][1]()

HEE_Entity_method1(E_po["thing"][2])

