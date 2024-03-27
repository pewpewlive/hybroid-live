
local test = require("/dynamic/import_test.hyb")
local map = {
	a = {
		b = 5
	}
}
local one = 1
local p = 1 - map["a"]["b"]
