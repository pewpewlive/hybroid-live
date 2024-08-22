


function E6HCCamera_Update(Self)
	E6HCCamera_Shake(Self)

end
function E6HCCamera_Shake(Self)
	H9 = fmath.random_fixedpoint(-Self[4], Self[4])
	Self[1] = H9

	pewpew.configure_player(0, {
		camera_x_override = Self[1]
	
})
	Self[4] = Self[4] * (Self[5])

end
function E6HCCamera_New()
	local Self = {0fx,0fx,0fx,10fx,0.2048fx}
	pewpew.add_update_callback(function()
		E6HCCamera_Update(Self)

	end)
	return Self
end

