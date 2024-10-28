local menu = true
local game = false

function love.draw()
    if (menu) then
        love.graphics.print("In Menu", 0, 0)
    end
    if (game) then
        love.graphics.print("Hello", 400, 300)
        -- TODO Show its your turn
    end
end

function love.update(dt)
    -- menu mode
    -- prosses tcp
    if (menu) then
        -- get inputs
        -- TODO set piece color
        -- TODO set server ip
        -- TODO Join Button
    end
    if (game) then
        -- get inputs
        -- TODO Forfet Button
    end
end

-- Comms
function MainComms(type, message)
    -- type of comm
end

-- recive
function ReciveTurn()

end

function ReciveMap()

end

function ReciveMove()

end

-- send
function RequestMap()

end

function RequestTurn()

end

function SendMove()

end

function SendJoin()

end

function SendQuit()

end
