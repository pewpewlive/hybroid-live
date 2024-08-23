local color_helper = require("/dynamic/helpers/color_helpers.lua")
local helper = {}

--- Creates a new empty mesh.
function helper.new_mesh()
  local mesh = {}
  mesh.vertexes = {}
  mesh.segments = {}
  mesh.colors = {}
  return mesh
end

function helper.lerp(a,b,t)--point a to point b, t has to be somewhere from 0 to 1
  local v = (1 - t) * a + b * t
  return v
end

function helper.invLerp(a,b,v)--point a to point b, but value is used
  local t = (v - a) / (b - a)
  return t
end

function helper.normalize(x,y)
  local mag = math.sqrt((x*x)+(y*y))
  if mag > 0 then
      return x/mag,y/mag
  end
end

function helper.copy(obj, seen)--some code from stack overflow that I might never understand
  if type(obj) ~= 'table' then return obj end
  if seen and seen[obj] then return seen[obj] end
  local s = seen or {}
  local res = setmetatable({}, getmetatable(obj))
  s[obj] = res
  for k, v in pairs(obj) do res[helper.copy(k, s)] = helper.copy(v, s) end
  return res
end

--- Adds a line to a mesh.
-- @params mesh table: a properly formated mesh. Possibly created with `new_mesh`.
-- @params vertex table: a list of either 2D or 3D vertexes. The coordinates of the vertexes are floats.
-- @params colors table: a list of colors. There should be as many colors as there are vertexes.
-- @params close_loop boolean: whether the line should be a closed loop, with the first vertex being linked with the last vertex.
function helper.add_line_to_mesh(mesh, vertexes, colors, close_loop)
  local vertex_count = #mesh.vertexes
  local color_count = #mesh.colors
  local segment_count = #mesh.segments
  local number_of_new_segments = #vertexes - 1
  local segments = {}

  for i = 1, #vertexes do
    table.insert(mesh.vertexes, vertexes[i])
    table.insert(mesh.colors, colors[i])
  end

  table.insert(segments, vertex_count)
  for i = 1, number_of_new_segments do
    table.insert(segments, vertex_count + i)
  end
  if close_loop then
    table.insert(segments, vertex_count)
  end
  table.insert(mesh.segments, segments)
end

function helper.add_poly(mesh, center, sides, color, radius, height)
  if sides > 2 then
    local x = center[1]
    local y = center[2]
    local z = center[3]
    local vertices = {} 
    local vertices2 = {} 
    local colors = {}

    --Plots the vertices
    for i = 1, sides do
      local angle = (math.pi * 2 * i)/sides
      table.insert(vertices, {x + radius * math.cos(angle),y + radius * math.sin(angle), z + height})
      table.insert(vertices2, {x + radius * math.cos(angle),y + radius * math.sin(angle), z - height})
      table.insert(colors, color)
    end

    --Draws the segments between each vertex
    helper.add_line_to_mesh(mesh, vertices, colors ,true)
    helper.add_line_to_mesh(mesh, vertices2, colors ,true)
    for i = 1, sides do
      helper.add_line_to_mesh(mesh, {vertices[i], vertices2[i]}, {color,color})
    end
  end
end

function helper.make_polygon(mesh,vertices)
  for i = 1, #vertices do
    if i == #vertices then
      helper.add_line_to_mesh(mesh,{vertices[1],vertices[i]})
    else
      helper.add_line_to_mesh(mesh,{vertices[i],vertices[i+1]})
    end
  end
end

function helper.add_flat_poly_angle(mesh, center, sides, color, radius, startangle)
  if sides > 2 then
    local x = center[1]
    local y = center[2]
    local z = center[3]
    local vertices = {} 
    local colors = {}
    for i = 1, sides do
      local angle = (math.pi * 2 * i)/sides + startangle
      table.insert(vertices, {x + radius * math.cos(angle),y + radius * math.sin(angle), z})
      table.insert(colors, color)
    end

    helper.add_line_to_mesh(mesh, vertices, colors ,true)
  end
end

function helper.add_flat_poly_angle_verts(mesh, center, sides, color, radius, startangle)
  if sides > 2 then
    local x = center[1]
    local y = center[2]
    local z = center[3]
    local vertices = {} 
    local colors = {}
    for i = 1, sides do
      local angle = (math.pi * 2 * i)/sides + startangle
      table.insert(mesh.vertexes, {x + radius * math.cos(angle),y + radius * math.sin(angle), z})
      table.insert(mesh.colors, color)
    end
  end
end

function helper.add_grid_map(mesh,center,color,width,height,frequency)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  for x = center[1],width-frequency[1]*2,frequency[1] do
    local a = {x+frequency[1],y,z}
    local b = {x+frequency[1],y+height,z}

    helper.add_line_to_mesh(mesh,{a,b},{color,color})
  end
  for y = center[2],height-frequency[2]*2,frequency[2] do
    local a = {x,y+frequency[2],z}
    local b = {x+width,y+frequency[2],z}

    helper.add_line_to_mesh(mesh,{a,b},{color,color})
  end
end

function helper.add_elevating_squares(mesh,center,color,width,height,end_z,frequency,z_extra)
  local x = center[1]
  local y = center[2]
  local z = center[3]

  local r = color[1]
  local g = color[2]
  local b = color[3]
  local a2 = color[4]

  for z = center[3], end_z, frequency do
    local a1 = {x-width/2-z/4,y-height/2-z/4,z-z_extra}
    local b1 = {x+width/2+z/4,y-height/2-z/4,z-z_extra}
    local c1 = {x+width/2+z/4,y+height/2+z/4,z-z_extra}
    local d1 = {x-width/2-z/4,y+height/2+z/4,z-z_extra}

    --[[if frequency % frequency*2 == 0 then
      z_extra = - z_extra
    end]]
    local a2 = a2 - z/3
    local true_color = color_helper.make_color(r,g,b,a2)
    local true_color2 = color_helper.make_color(r+80,g,b,a2)

    helper.add_line_to_mesh(mesh,{a1,b1,c1,d1},{true_color,true_color2,true_color,true_color2},true)
  end
end

function helper.add_cube_to_mesh(mesh, center, side_length, color)
  local half = side_length / 2
  local x = center[1]
  local y = center[2]
  local z = center[3]

  local a = {x - half, y - half, z - half}
  local b = {x - half, y + half, z - half}
  local c = {x + half, y + half, z - half}
  local d = {x + half, y - half, z - half}
  local e = {x - half, y - half, z + half}
  local f = {x - half, y + half, z + half}
  local g = {x + half, y + half, z + half}
  local h = {x + half, y - half, z + half}

  helper.add_line_to_mesh(mesh, {a, b, c, d}, {color, color, color, color}, true)
  helper.add_line_to_mesh(mesh, {e, f, g, h}, {color, color, color, color}, true)
 
  helper.add_line_to_mesh(mesh, {a, e}, {color, color})
  helper.add_line_to_mesh(mesh, {b, f}, {color, color})
  helper.add_line_to_mesh(mesh, {c, g}, {color, color})
  helper.add_line_to_mesh(mesh, {d, h}, {color, color})
end

function helper.add_dot(mesh, center, color_list)
  local x = center[1]
  local y = center[2]
  local z = center[3]

  local color = color_list[math.random(1,#color_list)]

  local a = {x-1,y,z}
  local b = {x+1,y,z}
  local a2 = {x,y-1,z}
  local b2 = {x,y+1,z}

  helper.add_line_to_mesh(mesh, {a,b}, {color, color})
  helper.add_line_to_mesh(mesh, {a2,b2}, {color, color})
end

function helper.add_stars(mesh,center,width,height,depth,amount)
  local x = center[1]
  local y = center[2]
  local z = center[3]

  for i = 1, amount do
    local x2 = fmath.random_int(x,width)
    local y2 = fmath.random_int(y,height)
    local z2 = fmath.random_int(z,depth)

    helper.add_dot(mesh,{x2,y2,z2},{0xffffff90,0xffffff70,0xffffff50,0xffffff30})
  end
end

function helper.add_sliced_rectangle(mesh,center,width,height,slice_freq,color)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  local vertices = {}
  local colors = {}
  local ww,hh = 0,height//slice_freq
  for i = 0, slice_freq do
    table.insert(vertices,{x-width/2+ww,y-height/2,z})
    ww = ww + width//slice_freq
  end
  ww = width//slice_freq
  for i = 1, slice_freq do
    table.insert(vertices,{x+width/2,y-height/2+hh,z})
    hh = hh + height//slice_freq
  end
  hh = height//slice_freq

  for i = 1, slice_freq do
    table.insert(vertices,{x+width/2-ww,y+height/2,z})
    ww = ww + width//slice_freq
  end
  ww = width//slice_freq
  for i = 1, slice_freq do
    table.insert(vertices,{x-width/2,y+height/2-hh,z})
    hh = hh + height//slice_freq
  end
  hh = height//slice_freq
  for i = 1, #vertices do
    table.insert(colors,color)
  end
  helper.add_line_to_mesh(mesh,vertices,colors,false)
end

function helper.add_circle(mesh, center, color, radius)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  local vertices = {}
  local colors = {}
  local vertex_amount = 50
  local increment = math.pi/vertex_amount
  for i = 0, vertex_amount do
      table.insert(vertices, {x + radius*(math.cos(increment*i*2)), y + radius*(math.sin(increment*i*2)), z})
      table.insert(colors,color)
  end

  helper.add_line_to_mesh(mesh, vertices, colors)
end

function helper.add_circle2(mesh, center, color, radius)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  local vertices = {}
  local colors = {}
  local vertex_amount = 50
  local increment = math.pi/vertex_amount
  for i = 0, vertex_amount do
      table.insert(vertices, {x+ radius*(math.sin(increment*i*2)), y , z+ radius*(math.cos(increment*i*2))})
      table.insert(colors,color)
  end

  helper.add_line_to_mesh(mesh, vertices, colors)
end

function helper.make_hemisphere(mesh,center,radius,color)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  local vertices = {}
  local colors = {}
  local vertex_amount = 50
  local increment = math.pi/vertex_amount
  for i = 0, vertex_amount-vertex_amount/2 do
      table.insert(vertices, {x + radius*(math.cos(increment*i*2)), y + radius*(math.sin(increment*i*2)), z})
      table.insert(colors,color)
  end

  helper.add_line_to_mesh(mesh, vertices, colors)
end

local function f(x)
  return math.sqrt(16-x*x)
end
local function f2(x)
  return math.sqrt(9-x*x)
end

function helper.add_sphere(mesh, center, color, radius)
  local x = center[1]
  local y = center[2]
  local z = center[3]
  local plus_reduce = 8
  local vertexes = {}
  local colors = {}

  for i = -3, 3 do
    for j = 1, (1+f(i)^2) do
      local angle = (math.pi * 2 * j)/(1+f(i)^2)
      local scale = radius
      local xchange = i * scale
      local ychange = math.sin(angle) * scale * f(i)
      local zchange = math.cos(angle) * scale * f(i)
      
      table.insert(vertexes,{x + xchange, y + ychange, z + zchange})
      table.insert(colors,color)
    end
    helper.add_line_to_mesh(mesh, vertexes, colors, true)
    colors = {}
    vertexes = {}
  end
end

function f3(x)
  return 2 * x * x * x * x
end

function dist(ex,ey,tx,ty)
  local dx,dy = tx-ex,ty-ey
  return math.sqrt((dx*dx)+(dy*dy))
end

function helper.add_cave(mesh,center,color,steps,rows,z_start,z_end,length,width_end,off)
  local x = center[1]
  local y = center[2]

  local space = 20
  local angle_step = (math.pi*2)/steps
  local offset = off or 0

  for i = 1, rows do
    local angle = (math.pi * 2 * i)/rows
    for j = 1, steps do 
      local sin, cos = math.sin(angle),math.cos(angle)
      local posx,posy = x+cos*(j*length+space+offset),y+sin*(j*length+space+offset)
      local t = helper.invLerp(0, steps*length,dist(x,y,posx,posy))
      local cwdidth = helper.lerp(0,width_end,t)
      t = f3(t)
      local z = helper.lerp(z_start,z_end,t)
      local square = {}
      --print("x: "..posx-posy*(cwdidth/2))
      --print("cos: "..cos..", sin: "..sin)
      table.insert(square,{posx-sin*(cwdidth/2),posy+cos*(cwdidth/2),z})
      table.insert(square,{posx+sin*(cwdidth/2),posy-cos*(cwdidth/2),z})
      local posx2,posy2 = x+cos*((j+1)*length+offset),y+sin*((j+1)*length+offset)
      --print("i: "..i..", j: "..j)

      t = helper.invLerp(0, steps*length,dist(x,y,posx2,posy2))
      cwdidth = helper.lerp(0,width_end,t)
      t = f3(t)
      --print("extend: "..t)
      z = helper.lerp(z_start,z_end,t)
      table.insert(square,{posx2+sin*(cwdidth/2),posy2-cos*(cwdidth/2),z})
      table.insert(square,{posx2-sin*(cwdidth/2),posy2+cos*(cwdidth/2),z})
      helper.add_line_to_mesh(mesh,square,{color,color,color,color},true)
    end
  end
end

return helper