

function E_Length(E_x, E_y)
	return fmath.sqrt(E_x * E_x + E_y * E_y)
end
function E_Normalize(E_x, E_y)
	local E_len = E_Length(E_x, E_y)

	return E_x / E_len, E_y / E_len
end
