Crafty.scene 'Loading', ->

    $text_css = { 'font-size': '24px', 'font-family': 'Arial', 'color': 'white', 'text-align': 'center'}
    Crafty.e('2D, DOM, Text').text('Loading...').attr({ x: 0, y: Antfarm.height()/2 - 24, w: Antfarm.width() }).css($text_css)
    async.parallel(
        constuction: (done) ->
            Antfarm.loadSprites "construction", ->
                done(null, true)
        material: (done) ->
            Antfarm.loadSprites "material", ->
                done(null, true)
        unit: (done) ->
            Antfarm.loadSprites "unit", ->
                done(null, true)
        world: (done) ->
            $.get '/world', (world) ->
                Antfarm.world = world
                done(null, true)
    , (error, result) ->
        Crafty.scene("Sim")
    )

Crafty.scene "Sim", ->
    Crafty.background("black")
    Antfarm.renderWorld()