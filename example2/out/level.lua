function split(str, pat)
   local t = {}  -- NOTE: use {n = 0} in Lua-5.0
   local fpat = "(.-)" .. pat
   local last_end = 1
   local s, e, cap = str:find(fpat, 1)
   while s do
      if s ~= 1 or cap ~= "" then
         table.insert(t,cap)
      end
      last_end = e+1
      s, e, cap = str:find(fpat, last_end)
   end
   if last_end <= #str then
      cap = str:sub(last_end)
      table.insert(t, cap)
   end
   return t
end

function ParseSound(link)
  local parts = split(link, '%%22')
  local sound = {}

  for i = 2, #parts,2 do
    local value = parts[i + 1]:sub(4, -4)
    if parts[i] == 'waveform' then
      value = parts[i + 2]
    end
    if parts[i] == 'amplification' then
      value = value / 100.0
    end
    if value == "true" then
      value = true
    end
    if value == "false" then
      value = false
    end
    sound[parts[i]] = value
  end
  return sound
end
function ToString(value)
	local str
	if type(value) == "table" then
		str = "{"
		if #value == 0 then
			for k, v in pairs(value) do
				str = str .. k .. ": " .. ToString(v) .. ", "
			end
		else
			for _, v in ipairs(value) do
				str = str .. ToString(v) .. ", "
			end
		end
		if str ~= "{" then
			str = string.sub(str, 0, string.len(str)-2)
		end
		str = str .. "}"
	else
		str = value
	end
	return str
end
require("/dynamic/quadro/quadro.lua")



E1WIDTH, E1HEIGHT = 1000fx, 1000fx

pewpew.set_level_size(E1WIDTH, 1000fx)
E1Ship = pewpew.new_player_ship(E1WIDTH / 2fx, E1HEIGHT / 2fx, 0)

pewpew.configure_player_ship_weapon(E1Ship, {
	frequency = pewpew.CannonFrequency.FREQ_10, cannon = pewpew.CannonType.DOUBLE

})
local E1playerConfig = {
	shield = 5, camera_distance = -50fx

}

pewpew.configure_player(0, E1playerConfig)
local E1bg = pewpew.new_customizable_entity(0fx, 0fx)

pewpew.customizable_entity_set_mesh(E1bg, "/dynamic/map_mesh.lua", 0)
E3HEQuadro_Spawn(0fx, 0fx, E1Ship)
pewpew.add_wall(750fx, 250fx, 500fx, 500fx)
local E1time = 0
pewpew.add_update_callback(function()
E1time = E1time + 1
	if pewpew.get_player_configuration(0).has_lost == true then
		pewpew.stop_game()
	end

	if E1time % 30 == 0 then
	end

end)
