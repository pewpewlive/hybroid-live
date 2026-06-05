
local _ = 1fx
local _ = 0.409fx

_ = 0.7fx
local E_e, _, _ = function ()
	return 2, 3
end, 80247, 4294945535
local _, _ = E_e()
local _ = true
HEE_Entity = {}
function HEE_Entity_Spawn(E_x, E_y)
	local id = pewpew.new_customizable_entity(E_x, E_y)
	HEE_Entity[id] = {}
	local Self = HEE_Entity[id]
	Self[1] = function () end
	Self[2] = 2
	Self[2] = 1
	return id
end
function HEE_Entity_Destroy(id, _)
	local Self = HEE_Entity[id]
	pewpew.entity_destroy(id)
	HEE_Entity[id] = nil
end
function HEE_Entity_method1(id)
	local Self = HEE_Entity[id]
	HEE_Entity_method2(id)

end
function HEE_Entity_method2(id)
	local Self = HEE_Entity[id]
	HEE_Entity_method1(id)

end
local E_k = HEE_Entity_Spawn(0fx, 0fx)
HEE_Entity[E_k][1]()
local E_mp = {HEE_Entity[E_k][1], function () end}
E_mp[1]()
local E_l = {E_k, HEE_Entity_Spawn(200fx, 200fx)}
HEE_Entity[E_l[2]][1]()
HEE_Entity_Destroy(E_l[2], 2)
HEE_Entity_Destroy(E_k, 2)
local E_po = {
	["thing"] = E_l
}
HEE_Entity[E_po["thing"][2]][1]()
HEE_Entity_method1(E_po["thing"][2])

local _ = ToString(2)
