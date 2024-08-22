



local function E0AddPolyWall(E0center, E0radius, E0sides)
	local E0x = E0center.x

	local E0y = E0center.y

	local E0vertices = {}

	local E0angle = fmath.tau() / fmath.to_fixed(E0sides)

	local E0sin, E0cos = fmath.sincos(E0angle)

	for _ = 1, E0sides, 1 do
		table.insert(E0vertices, {
			x = E0x + E0radius * E0cos, 
			y = E0y + E0radius * E0sin
		
})
		E0angle = E0angle + fmath.tau() / fmath.to_fixed(E0sides)

		H0, H1 = fmath.sincos(E0angle)
		E0sin, E0cos = H0, H1

		::GL_::
end
	for E0i, E0_ in ipairs(E0vertices) do
		if E0i == #E0vertices then
			pewpew.add_wall(E0vertices[E0i].x, E0vertices[E0i].y, E0vertices[1].x, E0vertices[1].y)
		else 
			pewpew.add_wall(E0vertices[E0i].x, E0vertices[E0i].y, E0vertices[E0i + 1].x, E0vertices[E0i + 1].y)
		end

		::GL2::
	end
end
