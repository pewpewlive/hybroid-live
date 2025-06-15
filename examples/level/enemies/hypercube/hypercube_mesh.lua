require("/dynamic/hypercube/settings.lua")
local gfx = require("/dynamic/helpers/graphics_helper.lua")
local m = require("/dynamic/helpers/math_helper.lua")

local frames = {}

local size = 80
local points = {
    {0, -size}, {0, size}
}
local color = 0xffffffff

local pos_map = {
    {-size, 0, 0, 0.125, 1},
    {0, size, 0.125, 0.25, 2},
    {size, size/3, 0.25, 0.5, 1},
    {size/3, -size/3, 0.5, 0.75, 0},
    {-size/3, -size, 0.75, 1, 2},
}
local size_map = {
    {size, size*1.5, 0, 0.125, 2},
    {size*1.5, size, 0.125, 0.25, 1},
    {size, size*0.5, 0.25, 0.5, 2},
    {size*0.5, size*0.5, 0.5, 0.75, 0},
    {size*0.5, size, 0.75, 1, 2},
}
local r1,g1,b1 = 70, 30, 255
local r2,g2,b2 = 255, 0, 0
local color_map = {
    {r1,g1,b1, r1,g1,b1, 0, 0.25, 0},
    {r1,g1,b1, r2,g2,b2, 0.25, 0.5, 0},
    {r2,g2,b2, r2,g2,b2, 0.5, 0.75, 0},
    {r2,g2,b2, r1,g1,b1, 0.75, 1, 0},
}
local time = 0
for i = 1, HYPERCUBE_FRAMES do 
    time = time + 1/HYPERCUBE_FRAMES

    table.insert(frames, gfx.new_mesh())
    
    gfx.add_poly(frames[i], {m.lerp_map(time, pos_map), 0, 0}, 4, m.lerp_map_color(time, color_map), m.lerp_map(time, size_map))
    
    local offset1 = m.wrap(0,1,time+0.25)
    gfx.add_poly(frames[i], {m.lerp_map(offset1, pos_map), 0, 0}, 4, m.lerp_map_color(offset1, color_map), m.lerp_map(offset1, size_map))

    local offset3 = m.wrap(0,1,time+0.5)
    gfx.add_poly(frames[i], {m.lerp_map(offset3, pos_map), 0, 0}, 4, m.lerp_map_color(offset3, color_map), m.lerp_map(offset3, size_map))

    local offset2 = m.wrap(0,1,time+0.75)
    gfx.add_poly(frames[i], {m.lerp_map(offset2, pos_map), 0, 0}, 4, m.lerp_map_color(offset2, color_map), m.lerp_map(offset2, size_map))
end


for i = 1, #frames do 
    --pewpew.print(#frames[i].vertexes)
    for j = 1, #frames[i].vertexes do 
        local num = m.wrap(-1, #frames[i].vertexes-1, j+3)
        --if num ~= j+3 then pewpew.print("num: " .. num .. "og: " .. j+3) end
        table.insert(frames[i].segments, {j-1,num})
    end
end