window.Antfarm =

    game: null

    grid:
        tileWidth: 32
        tileHeight: 32

    width: ->
        return window.innerWidth

    height: ->
        return window.innerHeight

    renderWorld: ->
        h = @game.World.Cells.length
        w = @game.World.Cells[0].length
        for y in [0...h]
            for x in [0...w]
                cell = @game.World.Cells[y][x]
                @renderCell(x, y, cell)
        for actor in @game.Actors
            @renderActor(actor)

    renderCell: (x, y, cell) ->
        c = "Spr" + cell.Material.titleize()
        if cell.Solid
            c += "Wall"
        else
            c += "Floor"
        c += 0
        Crafty.e("Gridded, " + c).at(x, y)

    renderActor: (actor) ->
        console.log(actor)
        c = "Spr" + actor.Tile
        x = actor.X
        y = actor.Y
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

    gameLoop: ->
        socket = io.connect("http://localhost:9000")
        socket.on 'connect', ->
            console.log("connecting")
            socket.emit('ping')
        socket.on 'pong', ->
            console.log("pong")
        socket.on 'game', (game) ->
            console.log("tic")
            console.log(game)
            Antfarm.game = game
            Crafty.scene("sim")
        socket.on 'disconnect', ->
            console.log "You have been disconnected"

Antfarm.start()