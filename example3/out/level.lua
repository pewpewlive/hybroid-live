require("/dynamic/quadro/quadro.lua")



E0WIDTH, E0HEIGHT = 1000fx, 1000fx

pewpew.set_level_size(E0WIDTH, 1000fx)
E0Ship = pewpew.new_player_ship(E0WIDTH / 2fx, E0HEIGHT / 2fx, 0)

pewpew.configure_player_ship_weapon(E0Ship, {
	frequency = pewpew.CannonFrequency.FREQ_10, 
	cannon = pewpew.CannonType.DOUBLE

})
pewpew.configure_player(0, {
	shield = 5, 
	camera_distance = -50fx

})
local E0bg = pewpew.new_customizable_entity(0fx, 0fx)

pewpew.customizable_entity_set_mesh(E0bg, "/dynamic/map_mesh.lua", 0)
E2Quadro_Spawn(0fx, 0fx, E0Ship)
local E0time = 0
pewpew.add_update_callback(function()
E0time = E0time + 1
	if pewpew.get_player_configuration(0).has_lost == true then
		pewpew.stop_game()
	end

	if E0time % 30 == 0 then
	end

end)
