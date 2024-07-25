function split(str, pat)
   local t = {}  -- NOTE: use {n = 0} in Lua-5.0
   local fpat = "(.-)" .. pat
   local last_end = 1
   local s, e, cap = string.find(str, fpat, 1)
   while s do
      if s ~= 1 or cap ~= "" then
         table.insert(t,cap)
      end
      last_end = e+1
      s, e, cap = string.find(str, fpat, last_end)
   end
   if last_end <= #str then
      cap = string.sub(str, last_end)
      table.insert(t, cap)
   end
   return t
end

function parseSound(link)
  local parts = split(link, '%%22')
  local sound = {}

  for i = 2, #parts,2 do
    local value = string.sub(parts[i + 1], 4, -4)
    if parts[i] == 'waveform' then
      value = parts[i + 2]
    end
    if parts[i] == 'amplification' then
      value = value / 100.0
    end
    if value == "true" then
      value = true
    end
    if value == "false" then
      value = false
    end
    sound[parts[i]] = value
  end
  return sound
end

function EditedSound(sound, new_kvs)
    local sound_copy = {}
    for k,v in pairs(sound) do
      sound_copy[k] = v
    end
    for k,v in pairs(new_kvs) do
      sound_copy[k] = v
    end
    return sound_copy
  end

--string.byte(s: string|number, i?:integer, j?:integer) -> ...integer
--string.char(byte:integer, ...integer) -> string
--string.dump(f: fun(...any):...unknown, strip?:boolean) -> string
--[[ string.find(s: string|number, pattern:string|number, init?:number, plain?:boolean) 
  -> start: integer|nil
  2. end: integer|nil
  3. ...any ]]
--string.format(s: string|number, ...any) -> string
