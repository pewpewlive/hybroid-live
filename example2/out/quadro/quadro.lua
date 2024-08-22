require("/dynamic/helpers/fmath_helpers.lua")



function E3HEQuadro_Move(id)
	local E3length = 1fx

	if not E3HEQuadro[id][10] and pewpew.entity_get_is_alive(E3HEQuadro[id][8]) then
		local E3sX, E3sY = pewpew.entity_get_position(E3HEQuadro[id][8])

		local E3dX, E3dY = E3sX - E3HEQuadro[id][1], E3sY - E3HEQuadro[id][2]

		H7 = E_Length(E3dX, E3dY)
		E3length = H7

		pewpew.print(E3length)
		E3HEQuadro[id][3] = E3HEQuadro[id][3] + (E3dX / E3length)

		E3HEQuadro[id][4] = E3HEQuadro[id][4] + (E3dY / E3length)

		E3HEQuadro[id][3] = E3HEQuadro[id][3] * (0.3686fx)

		E3HEQuadro[id][4] = E3HEQuadro[id][4] * (0.3686fx)

	end

	pewpew.entity_set_position(id, E3HEQuadro[id][1] + E3HEQuadro[id][3], E3HEQuadro[id][2] + E3HEQuadro[id][4])
end
function E3HEQuadro_UpdateCooldown(id)
	if not E3HEQuadro[id][9] then
		E3HEQuadro[id][6] = E3HEQuadro[id][6] - (1)

	end

	if E3HEQuadro[id][6] <= 0 and not E3HEQuadro[id][9] then
		E3HEQuadro[id][6] = 10

		E3HEQuadro[id][9] = true

	end

end
local function E3HEQuadroHCb0(id, E3normalX, E3normalY)
	if E3HEQuadro[id][10] then
		return 
	end

	H3, H4 = E_Reflect(E3HEQuadro[id][3], E3HEQuadro[id][4], E3normalX, E3normalY)
	E3HEQuadro[id][3], E3HEQuadro[id][4] = H3, H4

end
local function E3HEQuadroHCb1(id, E3playerIndex, E3shipId)
	if pewpew.entity_get_is_alive(E3shipId) and E3HEQuadro[id][9] then
		E3HEQuadro[id][3] = -E3HEQuadro[id][3] * 1.2048fx
		E3HEQuadro[id][4] = -E3HEQuadro[id][4] * 1.2048fx

		pewpew.add_damage_to_player_ship(E3shipId, 1)
	end

	E3HEQuadro[id][9] = false

end
local function E3HEQuadroHCb2(id, E3playerIndex, E3weaponType)
	if E3HEQuadro[id][10] then
		return false
	end

	if E3weaponType == pewpew.WeaponType.BULLET then
		E3HEQuadro[id][7] = E3HEQuadro[id][7] - (1)

		if E3HEQuadro[id][7] > 0 then
			pewpew.play_sound("/dynamic/quadro/quadro_sound.lua", 0, E3HEQuadro[id][1], E3HEQuadro[id][2])
		end

	end

	if E3HEQuadro[id][7] <= 0 then
		E3HEQuadro_Destroy(id)
	end

	return true
end
local function E3HEQuadroHCb3(id)
	H5, H6 = pewpew.entity_get_position(id)
	E3HEQuadro[id][1], E3HEQuadro[id][2] = H5, H6

	E3HEQuadro_UpdateCooldown(id)

	E3HEQuadro_Move(id)

end
E3HEQuadro = {}
function E3HEQuadro_Spawn(E3x, E3y, E3ship)
	local id = pewpew.new_customizable_entity(E3x, E3y)
	E3HEQuadro[id] = {0fx,0fx,0.40fx,0.40fx,5fx,0,5,nil,false,false}
	E3HEQuadro[id][8] = E3ship

	pewpew.customizable_entity_set_position_interpolation(id, true)
	pewpew.entity_set_position(id, E3x, E3y)
	pewpew.customizable_entity_set_mesh(id, "/dynamic/quadro/quadro_mesh.lua", 0)
	pewpew.entity_set_radius(id, 15fx)
	pewpew.customizable_entity_configure_wall_collision(id, true, E3HEQuadroHCb0)
	pewpew.customizable_entity_set_player_collision_callback(id, E3HEQuadroHCb1)
	pewpew.customizable_entity_set_weapon_collision_callback(id, E3HEQuadroHCb2)
	pewpew.entity_set_update_callback(id, E3HEQuadroHCb3)
	return id
end
function E3HEQuadro_Destroy(id)
	E3HEQuadro[id][10] = true

	E3HEQuadro[id][9] = false

	E3HEQuadro[id][6] = 999999999

	pewpew.play_sound("/dynamic/quadro/quadro_sound.lua", 1, E3HEQuadro[id][1], E3HEQuadro[id][2])
	pewpew.customizable_entity_start_exploding(id, 10)
end

