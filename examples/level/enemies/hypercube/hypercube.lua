local hc = {}

function mh.lerp(a, b, t, inter) 
    local inter_type = inter or 0
    if inter == 1 then
        t = t^1.5 -- it aint ^2 cuz it just dont look good
    elseif inter == 2 then
        t = t^0.9 -- same thing here
    end
    return a+(b-a)*t
end

-- map: {[{start_value, end_value, start_t, end_t, inter_type}]...}
function mh.lerp_map(t, map) 
    for i = 1, #map do 
        if map[i][4] >= t and map[i][3] <= t then
            return mh.remap(map[i][3], map[i][4], t, map[i][1], map[i][2], map[i][5])
        end
    end
    --pewpew.print("shi i aint found nun'")
    return 0
end

function mh.lerp_map_color(t, colormap) 
    for i = 1, #colormap do 
        if colormap[i][8] >= t and colormap[i][7] <= t then
            local color = cp.make_color(mh.remap(colormap[i][7], colormap[i][8], t, colormap[i][1], colormap[i][4], colormap[i][9]),
            mh.remap(colormap[i][7], colormap[i][8], t, colormap[i][2], colormap[i][5], colormap[i][9]),
            mh.remap(colormap[i][7], colormap[i][8], t, colormap[i][3], colormap[i][6], colormap[i][9]),
            255)
            return color
        end
    end
    --pewpew.print("shi i aint found nun': " .. t )
    return 0xffffffff
end

function mh.invlerp(a,b,v) 
    if b-a ~= 0 then
        return (v-a)/(b-a)
    else
        return 0
    end
end

-- params: list of {max:int, inter:INTERPOLATION_TYPE}, max should be between 0 and 1
function mh.parametric_remap(a,b, v, c,d, params) 
    local t = mh.invlerp(a, b, v)

    for i = 1, #params do  

        local param = params[i]

        if i <= param.max then 
            if param.inter == 1 then 
                t = math.sqrt(t)
            elseif param.inter == 2 then
                t = t*t
            end
        end

        return mh.lerp(c,d,t)
    end
end

function hc.new(x, y, rotate) 
    local id = pewpew.new_customizable_entity(x, y)
    pewpew.customizable_entity_set_mesh(id, "/dynamic/hypercube/graphics.lua", 0)
    
    local time = 0
    local frame_counter = 0
    pewpew.entity_set_update_callback(id, function(id)
        time = time + 1
        if time < 60 then return end
        pewpew.customizable_entity_set_mesh(id, "/dynamic/hypercube/graphics.lua", frame_counter % (HYPERCUBE_FRAMES-1))
        frame_counter = frame_counter + 1
        if rotate then 
            pewpew.customizable_entity_add_rotation_to_mesh(id, fmath.tau()/90, 1fx, 1fx, 0fx)
        end
    end)
end

return hc