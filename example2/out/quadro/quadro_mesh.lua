

local E4a = 20

meshes = {{
	vertexes = {{-E4a, -E4a}, {E4a, -E4a}, {E4a, E4a}, {-E4a, E4a}}, 
	segments = {{0, 1, 2, 3, 0}}, 
	colors = {0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff}

}}

local E4mesh = meshes[1]

table.insert(E4mesh.segments, {})
for E4i = 1, 10, 1 do
	table.insert(E4mesh.vertexes, {math.sin(E4i) * 10, math.cos(E4i) * 10})
	table.insert(E4mesh.segments[2], E4i + 2)
	table.insert(E4mesh.colors, 0xffffffff)
	::GL8::
end
