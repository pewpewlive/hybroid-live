pewpew.print("START")

grid = require("/dynamic/grid/grid.lua")
weapon_system = require("/dynamic/weapon_system.lua")
local pewpew = pewpew

local sqrt = fmath.sqrt

function length(ax, ay, bx, by)
  local dx, dy = ax - bx, ay - by

  return sqrt(dx * dx + dy * dy)
end

pewpew.enable_flag(0)

grid.init(18, 18, 50fx, 0fx, 0fx, 0fx)

local level_width, level_height = grid.get_size()
pewpew.set_level_size(level_width, level_height)

grid.set_colors(0x4445FF11, 0xFF221055)
grid.create_line_grid()

local ship_id = pewpew.new_player_ship(level_width / 2fx, level_height / 2fx, 0)
pewpew.configure_player(0, {camera_distance = 150fx})




local time = 0
local fx_time = 0fx



local wormhole_id = pewpew.new_customizable_entity(200fx, 200fx)
pewpew.entity_set_update_callback(wormhole_id, function ()
  local x, y = pewpew.entity_get_position(wormhole_id)
  local power = 15fx
  local radius = 300fx
  local kill_radius = 20fx
  local pulse_power, _ = (power / fmath.sincos(fx_time / 2fx) + 1fx)
    grid.pulse(x, y, -power * 2fx, radius)
  
  local entities = pewpew.get_entities_colliding_with_disk(x, y, 300fx)
  --pewpew.print(#entities)
  
  for key, id in ipairs(entities) do
    if not pewpew.entity_get_is_alive(id) or pewpew.entity_get_is_started_to_be_destroyed(id) then
      goto continue
    end

    local ex, ey = pewpew.entity_get_position(id)
    local dx, dy = x - ex, y - ey
    local len = sqrt(dx * dx + dy * dy)
    
    if len < radius then
      local factor = (1fx - (len / radius)) * power / 2fx * (len ^ -1)
      if pewpew.get_entity_type(id) == pewpew.EntityType.CUSTOMIZABLE_ENTITY then
        if pewpew.customizable_entity_get_tag(id) == 2 then
          weapon_system.pull_projectile(id, dx * factor, dy * factor)
        end
      else
        pewpew.entity_move(id, dx * factor, dy * factor)
      end
    end

    if len < kill_radius then
      if pewpew.get_entity_type(id) == pewpew.EntityType.CUSTOMIZABLE_ENTITY then
        if pewpew.customizable_entity_get_tag(id) == 2 then
          weapon_system.destroy_projectile(id)
        end
      elseif pewpew.get_entity_type(id) == pewpew.EntityType.SHIP then
        pewpew.add_damage_to_player_ship(ship_id, 1)
      end
    end
      ::continue::
  end
end)



--pewpew.add_wall(0fx, 0fx, 300fx, 300fx)

local ppx, ppy = 0fx, 0fx

local can_shoot = false
pewpew.add_update_callback(function()
  time = time + 1
  fx_time = fx_time + 1

  if pewpew.get_player_configuration(0)["has_lost"] == true then
    pewpew.stop_game()
    return
  end

  collectgarbage("collect")
  pewpew.print_debug_info()
  
  local _, _, sa, sd = pewpew.get_player_inputs(0)
  local px, py = pewpew.entity_get_position(ship_id)
  grid.pulse(px, py, 10fx, 100fx)
  
  if time % weapon_system.weapon_stats.shoot_freq == 0 then
    can_shoot = true
  end

  if sd ~= 0fx and can_shoot then
      local bullet_id = weapon_system.shoot(px, py, ppx - px, ppy - py, sa) 

    if can_shoot then
      time = time - time % weapon_system.weapon_stats.shoot_freq
    end
    can_shoot = false
  end

  -- if time % 30 == 0 then
  --   grid_pulse(500fx, 500fx, 600fx, 500fx)
  -- end

  ppx, ppy = px, py
end)
