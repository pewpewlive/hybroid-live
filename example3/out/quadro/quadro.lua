require("/dynamic/helpers/fmath_helpers.lua")


function Quadro_Move(id)
	H_, H0 = pewpew.entity_get_position(id)
	E2Quadro[Self][1], E2Quadro[Self][2] = H_, H0

	local E2sX, E2sY = 0fx, 0fx

	if pewpew.entity_get_is_alive(E2Quadro[Self][8]) then
		H1, H2 = pewpew.entity_get_position(E2Quadro[Self][8])
		E2sX, E2sY = H1, H2

	end

	local E2dirX, E2dirY = E_Normalize(E2sX - E2Quadro[Self][1], E2sY - E2Quadro[Self][2])

	E2Quadro[Self][3] = E2Quadro[Self][3] + ((E2dirX * 10fx * E2Quadro[Self][3]) * 0.819fx + 0.40fx)

	E2Quadro[Self][4] = E2Quadro[Self][4] + ((E2dirY * 10fx * E2Quadro[Self][4]) * 0.819fx + 0.40fx)

	pewpew.entity_set_position(id, E2Quadro[Self][1] + E2Quadro[Self][3] * 10fx, E2Quadro[Self][2] + E2Quadro[Self][4] * 10fx)
end
function Quadro_UpdateCooldown(id)
	E2Quadro[Self][5] = E2Quadro[Self][5] + (1)

	if E2Quadro[Self][5] == E2Quadro[Self][6] then
		E2Quadro[Self][9] = true

	end

end
local function E2QuadroHCb0(id, E2normalX, E2normalY)
end
local function E2QuadroHCb1(id, E2playerIndex, E2shipId)
	if pewpew.entity_get_is_alive(E2shipId) then
		pewpew.add_damage_to_player_ship(E2shipId, 1)
	end

	E2Quadro[Self][9] = false

end
local function E2QuadroHCb2(id, E2playerIndex, E2weaponType)
	if E2weaponType == pewpew.WeaponType.BULLET then
		E2Quadro[Self][7] = E2Quadro[Self][7] - (1)

	end

	if E2Quadro[Self][7] == 0 then
		E2Quadro[Self][10] = true


	end

end
local function E2QuadroHCb3(id)
	E2Quadro[Self].UpdateCooldown()
	E2Quadro[Self].Move()
end
E2Quadro = {}
function E2Quadro_Spawn(E2x, E2y, E2ship)
	local id = pewpew.new_customizable_entity(E2x, E2y)
	E2Quadro[id] = {0fx,0fx,0fx,0fx,0,5,20,nil,false,false}
	E2Quadro[id][8] = E2ship

	pewpew.customizable_entity_set_position_interpolation(id, true)
	pewpew.entity_set_position(id, E2x, E2y)
	pewpew.customizable_entity_set_mesh(id, "/dynamic/quadro/quadro_mesh.lua", 0)
	pewpew.entity_set_radius(id, 10.2048fx)
	pewpew.customizable_entity_configure_wall_collision(id, true, E2QuadroHCb0)
	pewpew.customizable_entity_set_player_collision_callback(id, E2QuadroHCb1)
	pewpew.customizable_entity_set_weapon_collision_callback(id, E2QuadroHCb2)
	pewpew.entity_set_update_callback(id, E2QuadroHCb3)
	return id
end
function E2Quadro_Destroy(id, E2DEATH_POWER)
	pewpew.customizable_entity_start_exploding(id, 10)
end

