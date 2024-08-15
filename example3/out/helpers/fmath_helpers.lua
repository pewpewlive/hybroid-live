

function E_Length(E_x, E_y)
	return fmath.sqrt(E_x * E_x + E_y * E_y)
end
function E_Normalize(E_x, E_y)
	if E_x == 0fx and E_y == 0fx then
		return 0fx, 0fx
	end

	local E_len = E_Length(E_x, E_y)

	return E_x / E_len, E_y / E_len
end
