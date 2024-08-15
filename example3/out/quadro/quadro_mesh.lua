

local E3a = 20

meshes = {{
	vertexes = {{-E3a, -E3a}, {E3a, -E3a}, {E3a, E3a}, {-E3a, E3a}}, 
	segments = {{0, 1, 2, 3, 0}}, 
	colors = {0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff}

}}

local E3mesh = meshes[1]

table.insert(E3mesh.segments, {})
for E3i = 1, 10, 1 do
	table.insert(E3mesh.vertexes, {math.sin(E3i) * 10, math.cos(E3i) * 10})
	table.insert(E3mesh.segments[2], E3i + 2)
	table.insert(E3mesh.colors, 0xffffffff)
	::GL3::
end
