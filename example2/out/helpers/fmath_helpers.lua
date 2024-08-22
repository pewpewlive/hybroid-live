

function E_Length(E_x, E_y)
	return fmath.sqrt(E_x * E_x + E_y * E_y)
end
function E_Normalize(E_x, E_y)
	if E_x == 0fx and E_y == 0fx then
		return 0fx, 0fx
	end

	local E_len = E_Length(E_x, E_y)

	if E_len == 0fx then
		return 0fx, 0fx
	end

	return E_x / E_len, E_y / E_len
end
function E_Clamp(E_val, E_min, E_max)
	if E_val > E_max then
		return E_max
	end

	if E_val < E_min then
		return E_min
	end

	return E_val
end
function E_Lerp(E_a, E_b, E_t)
	return E_a + (E_b - E_a) * E_t
end
function E_InvLerp(E_a, E_b, E_v)
	return (E_b - E_a) / (E_v - E_a)
end
function E_Dot(E_x, E_y, E_nX, E_nY)
	return E_x * E_nX + E_y * E_nY
end
function E_Reflect(E_x, E_y, E_nX, E_nY)
	local E_dot = E_Dot(E_x, E_y, E_nX, E_nY) * 2fx

	E_x = E_x - (E_nX * E_dot)

	E_y = E_y - (E_nY * E_dot)

	return E_x, E_y
end
