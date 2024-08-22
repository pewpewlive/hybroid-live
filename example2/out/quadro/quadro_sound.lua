function split(str, pat)
   local t = {}  -- NOTE: use {n = 0} in Lua-5.0
   local fpat = "(.-)" .. pat
   local last_end = 1
   local s, e, cap = str:find(fpat, 1)
   while s do
      if s ~= 1 or cap ~= "" then
         table.insert(t,cap)
      end
      last_end = e+1
      s, e, cap = str:find(fpat, last_end)
   end
   if last_end <= #str then
      cap = str:sub(last_end)
      table.insert(t, cap)
   end
   return t
end

function ParseSound(link)
  local parts = split(link, '%%22')
  local sound = {}

  for i = 2, #parts,2 do
    local value = parts[i + 1]:sub(4, -4)
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

sounds = {{
	sampleRate = 44100, 
	attack = 0, 
	sustain = 0.07, 
	sustainPunch = 0, 
	decay = 0.06, 
	tremoloDepth = 0, 
	tremoloFrequency = 10, 
	frequency = 700, 
	frequencySweep = -1000, 
	frequencyDeltaSweep = -500, 
	repeatFrequency = 0, 
	frequencyJump1Onset = 33, 
	frequencyJump1Amount = 0, 
	frequencyJump2Onset = 66, 
	frequencyJump2Amount = 0, 
	harmonics = 0, 
	harmonicsFalloff = 0.5, 
	waveform = "sawtooth", 
	interpolateNoise = true, 
	vibratoDepth = 0, 
	vibratoFrequency = 10, 
	squareDuty = 50, 
	squareDutySweep = 0, 
	flangerOffset = 0, 
	flangerOffsetSweep = 0, 
	bitCrush = 16, 
	bitCrushSweep = 0, 
	lowPassCutoff = 22050, 
	lowPassCutoffSweep = -3500, 
	highPassCutoff = 0, 
	highPassCutoffSweep = 0, 
	compression = 1, 
	normalization = true, 
	amplification = 100

}, ParseSound("https://pewpew.live/jfxr/#%7B%22_version%22%3A1%2C%22_name%22%3A%22Explosion%201%22%2C%22_locked%22%3A%5B%5D%2C%22sampleRate%22%3A44100%2C%22attack%22%3A0%2C%22sustain%22%3A0.05%2C%22sustainPunch%22%3A0%2C%22decay%22%3A0.48%2C%22tremoloDepth%22%3A22%2C%22tremoloFrequency%22%3A96%2C%22frequency%22%3A8800%2C%22frequencySweep%22%3A-4800%2C%22frequencyDeltaSweep%22%3A-2800%2C%22repeatFrequency%22%3A0%2C%22frequencyJump1Onset%22%3A33%2C%22frequencyJump1Amount%22%3A0%2C%22frequencyJump2Onset%22%3A66%2C%22frequencyJump2Amount%22%3A0%2C%22harmonics%22%3A0%2C%22harmonicsFalloff%22%3A0.5%2C%22waveform%22%3A%22whitenoise%22%2C%22interpolateNoise%22%3Atrue%2C%22vibratoDepth%22%3A0%2C%22vibratoFrequency%22%3A10%2C%22squareDuty%22%3A50%2C%22squareDutySweep%22%3A0%2C%22flangerOffset%22%3A1%2C%22flangerOffsetSweep%22%3A3%2C%22bitCrush%22%3A16%2C%22bitCrushSweep%22%3A0%2C%22lowPassCutoff%22%3A22050%2C%22lowPassCutoffSweep%22%3A0%2C%22highPassCutoff%22%3A0%2C%22highPassCutoffSweep%22%3A0%2C%22compression%22%3A1%2C%22normalization%22%3Atrue%2C%22amplification%22%3A100%7D")}

