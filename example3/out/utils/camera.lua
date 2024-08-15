


function HSCamera_Update(Self)
	Self.Shake()
end
function HSCamera_Shake(Self)
	H4 = fmath.random_fixedpoint(-Self[4], Self[4])
	Self[1] = H4

	pewpew.configure_player(0, {
		camera_x_override = Self[1]
	
})
	Self[4] = Self[4] * (Self[5])

end
function HSE4Camera_New()
	local Self = {0fx,0fx,0fx,10fx,0.2048fx}
	pewpew.add_update_callback(function()
		Self.Update()
	end)
	return Self
end

