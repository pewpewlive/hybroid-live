local move_angle, move_distance, shoot_angle, shoot_distance = pewpew.get_player_inputs(0)

self.prev_input_distance = move_distance

--if self.first_input_released then pewpew.print("true") else pewpew.print("false") end

if self.dash_cooldown_timer > 0 and self.dash_time_left < 1 then
	self.dash_cooldown_timer = self.dash_cooldown_timer - 1
	return
end

pewpew.print(self.input_timer)

if self.input_timer > 0 then
	self.input_timer = self.input_timer - 1
end

if self.input_timer < 1 then
	self.input_timer = 0
end

if move_distance > 0.5fx and self.dash_time_left < 1 and self.input_timer < 1 and self.first_input_done == false then
	self.first_input_done = true
	self.input_timer = self.input_time
	self.first_input_angle = move_angle 
end

if move_distance < 0.5fx and self.input_timer < 1 and self.first_input_done == true then
	self.first_input_done = false
end

if move_distance < 0.5fx and self.input_timer > 0 and self.first_input_released == false and self.first_input_done == true then
	self.first_input_released = true
	--self.input_timer = self.input_time
end

if move_distance < 0.5fx and self.input_timer < 1 and self.first_input_released == true and self.first_input_done == true then
	self.first_input_done = false
	self.first_input_released = false
end

if self.first_input_done == true and self.first_input_released == true and self.input_timer > 0 and move_distance > 0.5fx --[[ and approx(self.dash_direction, move_angle, 0.3000fx)]]then
	self.first_input_released = false
	self.first_input_done = false
	self.input_timer = 0

	-- CAN HAVE A BUG WHEN THERE ARE SPEED BOOSTS ACTIVE
	self.saved_movement_speed = default_ship_speed
	pewpew.set_player_ship_speed(ship_id, 0fx, 0fx, -1)
	
	self.dash_time_left = self.dash_duration
	
	self.dash_dy, self.dash_dx = fmath.sincos(move_angle)
	self.dash_dy, self.dash_dx = self.dash_dy * 60fx, self.dash_dx * 60fx
	self.dash_direction = move_angle
	
	pewpew.make_player_ship_transparent(ship_id, self.dash_duration + 15)
end

if self.dash_time_left > 0 then
	self.dash_time_left = self.dash_time_left - 1
	self.dash_cooldown_timer = self.dash_cooldown 

	--self.dash_angle = self.dash_angle * 0.4000fx

	local sx, sy = pewpew.entity_get_position(ship_id)
	
	local entities_caught = pewpew.get_entities_colliding_with_disk(sx, sy, 20fx)
	for _, eid in ipairs(entities_caught) do
		pewpew.entity_react_to_weapon(eid, {
			type = pewpew.WeaponType.ATOMIZE_EXPLOSION, 
			x = sx, 
			y = sy, 
			player_index = 0
		})
	end

	for i = 0fx, 10fx, 1fx do
		local angle = (i - 1fx) * 2fx * fmath.tau() / 2fx / 10fx
		local pdy, pdx = fmath.sincos(angle)

		pewpew.add_particle(sx, sy, 0fx,
		pdx * fmath.random_fixedpoint(0.2000fx, 3fx),
		pdy * fmath.random_fixedpoint(0.2000fx, 3fx),
		0fx,
		color_helpers.make_color(255, 255, 0, 255),
		10
		)
	end

	self.dash_dx = self.dash_dx * 0.3400fx
	self.dash_dy = self.dash_dy * 0.3400fx
	
	pewpew.entity_set_position(ship_id, sx + self.dash_dx, sy + self.dash_dy)
end

if self.dash_time_left < 1 then
	pewpew.set_player_ship_speed(ship_id, self.saved_movement_speed, 0fx, -1)
end

self.prev_input_distance = move_distance