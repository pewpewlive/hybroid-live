
local E0a = 4

local E0test = {{1}}

meshes = {{
	vertexes = {{-E0a, -E0a}, {E0a, -E0a}, {E0a, E0a}, {-E0a, E0a}}, 
	segments = {{0, 1, 2, 3, 0}}, 
	colors = {0xffffffff}

}}

meshes[1].vertexes[1] = {1, 1}
