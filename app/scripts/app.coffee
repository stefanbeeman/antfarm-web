window.Antfarm =

    world: null

    grid:
        tileWidth: 32
        tileHeight: 32

    width: ->
        return window.innerWidth

    height: ->
        return window.innerHeight

    renderWorld: ->
        h = @world.Cells.length
        w = @world.Cells[0].length
        for y in [0...h]
            for x in [0...w]
                cell = @world.Cells[y][x]
                cell = cell.Data
                @renderCell(x, y, cell)
        for actor in @world.Actors
            @renderActor(actor)

    renderCell: (x, y, cell) ->
        material = @world.Materials[cell.material]
        c = "Spr" + material.Name
        if cell.solid
            c += "Wall"
        else
            c += "Floor"
        c += 0
        Crafty.e("Gridded, " + c).at(x, y)

    renderActor: (actor) ->
        c = "SprWorm"
        x = actor.Location.X
        y = actor.Location.Y
        Crafty.e("Unit, " + c).at(x, y)

    start: ->
        Crafty.init(@width(), @height())
        Crafty.background('black')
        Crafty.scene('Loading')

    loadSprites: (filename, fn) ->
        async.parallel(
            data: (done) ->
                $.get '/data/sprites_' + filename + '.yml', (yml) ->
                    data = jsyaml.load(yml)
                    done(null, data)
            gfx: (done) ->
                Crafty.load ['/gfx/tiles/' + filename + '.png'], ->
                    done(null, true)
        , (err, results) ->
            Crafty.sprite(32, '/gfx/tiles/' + filename + '.png', results.data)
            fn()
        )

Antfarm.start()