
local E_a = "string %s%v"

local E_b = string.byte(E_a, 2, 4)

local E_formatted = string.format(E_a, "a", 1)

