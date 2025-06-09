local lerp_colors = require("/dynamic/color_helpers.lua").lerp_colors
local grid = {
  debug = false
}

local grid_size_x, grid_size_y = 20, 20
local cell_size = 45fx
local cell_half_size = cell_size / 2fx
local grid_count = grid_size_x * grid_size_y

local offset_x, offset_y, offset_z = 0fx, 0fx, 0fx
local grid_color = 0x4445FF00
local active_grid_color = 0xFF221055

local level_width = (grid_size_x - 1) * cell_size
local level_height = (grid_size_y - 1) * cell_size

-- localizing the functions because its faster to call them this way
local pewpew = pewpew
local fmath = fmath

local sqrt = fmath.sqrt
local atan2 = fmath.atan2
local to_int = fmath.to_int
local print = pewpew.print
local insert = table.insert

local pewpew = pewpew
local entity_set_position = pewpew.entity_set_position
local customizable_entity_set_mesh_angle = pewpew.customizable_entity_set_mesh_angle
local customizable_entity_set_mesh_color = pewpew.customizable_entity_set_mesh_color
local customizable_entity_set_mesh_xyz_scale = pewpew.customizable_entity_set_mesh_xyz_scale

local grid_container = {}

function length(ax, ay, bx, by) -- ported
  local dx, dy = ax - bx, ay - by

  return sqrt(dx * dx + dy * dy)
end

local function clamp(min, max, num) -- ported
  if num > max then
    return max
  elseif num < min then
    return min
  end

  return num
end

function grid.get_point(x, y)
  return grid_container[(x - 1) * grid_size_x + y]
end

local get_point_local = grid.get_point

function grid.world_to_grid_position(x, y)
  return clamp(1, grid_size_x, to_int((x + cell_half_size) / cell_size) + 1), -- x
         clamp(1, grid_size_y, to_int((y + cell_half_size) / cell_size) + 1)  -- y
end

local world_to_grid_pos_local = grid.world_to_grid_position

local grid_lines = {}
local function grid_line_update(id)
  local base_point = grid_lines[id][1]
  local following_point = grid_lines[id][2]

  local px, py = base_point[3], base_point[4]
  entity_set_position(id, px, py)

  local dx, dy = following_point[3] - px, following_point[4] - py
  local mag = sqrt(dx * dx + dy * dy)

  customizable_entity_set_mesh_xyz_scale(id, mag, 1fx, 0fx, 0fx)
  customizable_entity_set_mesh_angle(id, atan2(dy, dx), 0fx, 0fx, 1fx)

  customizable_entity_set_mesh_color(id,
    lerp_colors(grid_color, active_grid_color,
      clamp(0fx, 80fx, (base_point[7] + following_point[7]) + (mag - cell_size)) / 80fx))
end

-- bp: base point; fp: following point
local function new_grid_line(base_point, following_point)
  local id = pewpew.new_customizable_entity(base_point[3], base_point[4])
  grid_lines[id] = { base_point, following_point }
  pewpew.customizable_entity_set_mesh(id, "/dynamic/grid/grid_meshes.lua", 0)

  pewpew.customizable_entity_set_position_interpolation(id, true)
  pewpew.customizable_entity_set_angle_interpolation(id, true)

  pewpew.customizable_entity_set_mesh_z(id, offset_z)
  pewpew.customizable_entity_start_spawning(id, 0)
  
  grid_line_update(id)
  pewpew.entity_set_update_callback(id, grid_line_update)
end 

local grid_custom_lines = {}
local function grid_custom_line_update(id)
  local base_point = grid_custom_lines[id][1]
  local following_point = grid_custom_lines[id][2]

  local px, py = base_point[3], base_point[4]
  entity_set_position(id, px, py)

  local dx, dy = following_point[3] - px, following_point[4] - py
  local mag = sqrt(dx * dx + dy * dy)

  customizable_entity_set_mesh_xyz_scale(id, mag / grid_custom_lines[id][3], 1fx, 0fx, 0fx)
  customizable_entity_set_mesh_angle(id, atan2(dy, dx), 0fx, 0fx, 1fx)

  customizable_entity_set_mesh_color(id,
    lerp_colors(grid_color, active_grid_color,
      clamp(0fx, 80fx, (base_point[7] + following_point[7]) + (mag - cell_size)) / 80fx))
end

-- bp: base point; fp: following point
local function new_grid_custom_line(base_point, following_point, mesh_info)
  local id = pewpew.new_customizable_entity(base_point[3], base_point[4])
  grid_custom_lines[id] = { base_point, following_point, mesh_info[3] }
  pewpew.customizable_entity_set_mesh(id, mesh_info[1], mesh_info[2])
  
  pewpew.customizable_entity_set_position_interpolation(id, true)
  pewpew.customizable_entity_set_angle_interpolation(id, true)

  pewpew.customizable_entity_set_mesh_z(id, offset_z)
  pewpew.customizable_entity_start_spawning(id, 0)
  
  grid_custom_line_update(id)
  pewpew.entity_set_update_callback(id, grid_custom_line_update)
end

local grid_points = {}
local function grid_point_update(id)
  local base_point = grid_points[id]

  entity_set_position(id, base_point[3], base_point[4])

  customizable_entity_set_mesh_color(id,
    lerp_colors(grid_color, active_grid_color,
      clamp(0fx, 80fx, base_point[7] * 2fx) / 80fx))
end

-- bp: base point; fp: following point
local function new_grid_point(base_point)
  local id = pewpew.new_customizable_entity(base_point[3], base_point[4])
  grid_points[id] = base_point
  pewpew.customizable_entity_set_mesh(id, "/dynamic/grid/grid_meshes.lua", 1)

  pewpew.customizable_entity_set_position_interpolation(id, true)

  pewpew.customizable_entity_start_spawning(id, 0)
  pewpew.customizable_entity_set_mesh_z(id, offset_z)

  grid_point_update(id)
  pewpew.entity_set_update_callback(id, grid_point_update)
end

-- bp: base point; fp: following point
local function new_grid_custom_point(base_point, mesh_info)
  local id = pewpew.new_customizable_entity(base_point[3], base_point[4])
  grid_points[id] = base_point
  pewpew.customizable_entity_set_mesh(id, mesh_info[1], mesh_info[2])

  pewpew.customizable_entity_set_position_interpolation(id, true)

  pewpew.customizable_entity_start_spawning(id, 0)
  pewpew.customizable_entity_set_mesh_z(id, offset_z)

  grid_point_update(id)
  pewpew.entity_set_update_callback(id, grid_point_update)
end

function grid.pulse(x, y, power, radius)
  local left_extent, bottom_extent = world_to_grid_pos_local(x - radius, y - radius)
  local right_extent, top_extent = world_to_grid_pos_local(x + radius, y + radius)

  for gx = left_extent, right_extent, 1 do
    for gy = bottom_extent, top_extent, 1 do
      local point = grid_container[(gx - 1) * grid_size_x + gy]

      local dx, dy = point[1] - x, point[2] - y
      local len = fmath.sqrt(dx * dx + dy * dy)

      if len < radius then
        local factor = (1fx - (len / radius)) * power * (len ^ -1)
        point[5], point[6] = point[5] + dx * factor, -- x
                             point[6] + dy * factor  -- y
      end
    end
  end
end

function grid.init(_size_x, _size_y, _cell_size, _offset_x, _offset_y, _offset_z) 
  grid_size_x, grid_size_y = _size_x, _size_y
  cell_size = _cell_size
  cell_half_size = cell_size / 2fx

  offset_x, offset_y, offset_z = _offset_x, _offset_y, _offset_z

  level_width = (grid_size_x - 1) * cell_size
  level_height = (grid_size_y - 1) * cell_size

  grid_count = grid_size_x * grid_size_y

  for x = 1fx, fmath.to_fixedpoint(grid_size_x), 1fx do
    for y = 1fx, fmath.to_fixedpoint(grid_size_y), 1fx do
      -- layout: 1, 2: anchor point; 3, 4 current position; 5, 6 velocity; 7 "power"
      local px = (x - 1fx) * cell_size + offset_x
      local py = (y - 1fx) * cell_size + offset_y
      
      insert(grid_container, { px, py, px, py, 0fx, 0fx, 0fx })

      if grid.debug then
        local debug_id = pewpew.new_customizable_entity(px, py)
        pewpew.customizable_entity_set_mesh_scale(debug_id, 0.1500fx)
        pewpew.customizable_entity_set_string(debug_id, "#ffffff55o")
      end
    end
  end

  pewpew.add_update_callback(function()
    for i = 1, grid_count, 1 do
      local point = grid_container[i]
      local pos_x, pos_y = point[3], point[4]
      local vel_x, vel_y = point[5], point[6]
  
      vel_x = (vel_x + (point[1] - pos_x)) * 0.2248fx
      vel_y = (vel_y + (point[2] - pos_y)) * 0.2248fx
  
      pos_x = pos_x + vel_x
      pos_y = pos_y + vel_y
  
      point[5], point[6] = vel_x, vel_y
      point[3], point[4] = pos_x, pos_y
      point[7] = length(point[1], point[2], pos_x, pos_y)
    end
  end)
end

function grid.set_colors(default, active)
  grid_color = default
  active_grid_color = active
end

function grid.create_line_grid(mesh_info)
  for x = 1, grid_size_x, 1 do
    for y = 1, grid_size_y, 1 do
      if x < grid_size_x then
        if mesh_info == nil then
          new_grid_line(get_point_local(x, y),  get_point_local(x + 1, y))
        else
          new_grid_custom_line(get_point_local(x, y),  get_point_local(x + 1, y), mesh_info)
        end
      end
      if y < grid_size_y then
        if mesh_info == nil then
          new_grid_line(get_point_local(x, y),  get_point_local(x, y + 1))
        else
          new_grid_custom_line(get_point_local(x, y),  get_point_local(x, y + 1), mesh_info)
        end
      end
    end
  end
end

function grid.create_point_grid(mesh_info)
  for x = 1, grid_size_x, 1 do
    for y = 1, grid_size_y, 1 do
      if mesh_info == nil then
        new_grid_point(get_point_local(x, y))
      else
        new_grid_custom_point(get_point_local(x, y), mesh_info)
      end
    end
  end
end

function grid.create_cross_lines_grid(mesh_info)
  for x = 1, grid_size_x, 1 do
    for y = 1, grid_size_y, 1 do
      if x < grid_size_x and y < grid_size_y then
        if mesh_info == nil then
          new_grid_line(get_point_local(x, y), get_point_local(x + 1, y + 1))
        else
          new_grid_custom_line(get_point_local(x, y), get_point_local(x + 1, y + 1), mesh_info)
        end
      end
      if x > 1 and y < grid_size_y then
        if mesh_info == nil then
          new_grid_line(get_point_local(x, y),  get_point_local(x - 1, y + 1))
        else
          new_grid_custom_line(get_point_local(x, y),  get_point_local(x - 1, y + 1), mesh_info)
        end
      end
    end
  end
end

function grid.get_size()
  return level_width, level_height
end

return grid