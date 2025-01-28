module DrawingUtils exposing (display)

import Svg exposing (..)
import Svg.Attributes exposing (..)
import TcTurtle exposing (Instruction(..), Program)


type alias Position =
    { x : Float, y : Float, angle : Float }

type alias Color =
    { red : Int, green : Int, blue : Int, alpha : Float }


-- Convert a Color to a valid CSS RGBA string

colorToString : Color -> String
colorToString color =
    "rgba("
        ++ String.fromInt color.red ++ ", "
        ++ String.fromInt color.green ++ ", "
        ++ String.fromInt color.blue ++ ", "
        ++ String.fromFloat color.alpha ++ ")"

-- Execute a single turtle instruction and return the new position and SVG elements

execute : Position -> Instruction -> Color -> ( Position, List (Svg msg) )
execute pos instruction color =
    case instruction of
        Forward n ->
            let
                newX =
                    pos.x + toFloat n * cos (degrees pos.angle)
                newY =
                    pos.y + toFloat n * sin (degrees pos.angle)
                lineSvg =
                    line
                        [ x1 (String.fromFloat pos.x)
                        , y1 (String.fromFloat pos.y)
                        , x2 (String.fromFloat newX)
                        , y2 (String.fromFloat newY)
                        , stroke (colorToString color)
                        ]
                        []
            in
            ( { pos | x = newX, y = newY }, [ lineSvg ] )

        Right degrees ->
            ( { pos | angle = pos.angle - toFloat degrees }, [] )

        Left degrees ->
            ( { pos | angle = pos.angle + toFloat degrees }, [] )

        Repeat count instructions ->
            let
                ( newPos, svgs ) =
                    List.foldl
                        (\_ ( p, acc ) ->
                            List.foldl
                                (\instr ( pNext, accNext ) ->
                                    let
                                        ( nextPos, svg ) =
                                            execute pNext instr color
                                    in
                                    ( nextPos, accNext ++ svg )
                                )
                                ( p, acc )
                                instructions
                        )
                        ( pos, [] )
                        (List.repeat count ())
            in
            ( newPos, svgs )



-- Convert a complete program into SVG with support for zoom and color

display : Program -> Float -> Color -> Svg msg
display program zoom color =
    let
        initialPosition =
            { x = 500, y = 500, angle = 0 }

        ( _, svgs ) =
            List.foldl
                (\instr ( pos, acc ) ->
                    let
                        ( newPos, svg ) =
                            execute pos instr color
                    in
                    ( newPos, acc ++ svg )
                )
                ( initialPosition, [] )
                program

        center =
            500 / zoom

        viewBoxSize =
            1000 / zoom
    in
    svg 
        [ viewBox (String.join " " [ String.fromFloat (viewBoxSize - 2 * center), String.fromFloat (viewBoxSize - 2 * center), String.fromFloat viewBoxSize, String.fromFloat viewBoxSize ])
        , width "500"
        , height "500"
        ] 
        svgs
