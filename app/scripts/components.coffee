Crafty.c 'Gridded',
    init: ->
        @requires('2D, DOM')
        @attr(
            w: Antfarm.grid.tileWidth
            h: Antfarm.grid.tileHeight
        )
    at: (x, y) ->
        unless x? && y?
            return {
                x: this.x/Antfarm.grid.tileWidth
                y: this.y/Antfarm.gird.tileHeight
            }
        else
            @attr(
                x: x * Antfarm.grid.tileWidth
                y: y * Antfarm.grid.tileHeight
            )
            return this