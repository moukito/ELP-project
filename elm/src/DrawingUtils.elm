module DrawingUtils exposing (Color, display)

{-| This module handles the rendering of TcTurtle programs into SVG elements.

Main functions:

  - `colorToString`: Converts a `Color` to an RGBA CSS string.
  - `execute`: Executes a single turtle instruction and returns updated state and SVG elements.
  - `display`: Converts a full program into an SVG graphic.

Types:

  - `Position`: Represents the turtle's position and orientation.
  - `Color`: Represents the drawing color in RGBA format.

-}

import Svg exposing (..)
import Svg.Attributes exposing (..)
import TcTurtle exposing (Instruction(..), Program)



-- TYPES


{-| Represents the turtle's position and orientation.

Fields:

  - `x`: X-coordinate of the turtle's position.
  - `y`: Y-coordinate of the turtle's position.
  - `angle`: Turtle's heading angle, in degrees.

-}
type alias Position =
    { x : Float, y : Float, angle : Float }


{-| Represents the Turtle's drawing color.

Fields:

  - `red`: Int (0-255) - Red color component.
  - `green`: Int (0-255) - Green color component.
  - `blue`: Int (0-255) - Blue color component.
  - `alpha`: Float (0.0-1.0) - Alpha (opacity) component.

-}
type alias Color =
    { red : Int, green : Int, blue : Int, alpha : Float }



-- UTILITY FUNCTIONS


{-| Converts a `Color` to a valid CSS RGBA string.

Parameters:

  - `color`: The `Color` to convert.

Returns:

  - A string in `rgba(r, g, b, a)` format.

-}
colorToString : Color -> String
colorToString color =
    "rgba("
        ++ String.fromInt color.red
        ++ ", "
        ++ String.fromInt color.green
        ++ ", "
        ++ String.fromInt color.blue
        ++ ", "
        ++ String.fromFloat color.alpha
        ++ ")"



-- EXECUTION


{-| Executes a single turtle instruction and returns the updated position and SVG elements.

Parameters:

  - `pos`: The current `Position` of the turtle.
  - `instruction`: The `Instruction` to execute.
  - `color`: The drawing `Color`.

Returns:

  - A tuple with the updated `Position` and a list of SVG elements.

-}
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



-- RENDERING


{-| Converts a TcTurtle program into an SVG graphic.

Parameters:

  - `program`: The `Program` to draw.
  - `zoom`: The zoom level for rendering.
  - `color`: The drawing color.

Returns:

  - An `Svg` element representing the full drawing.

-}
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
