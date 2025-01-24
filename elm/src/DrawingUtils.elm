module DrawingUtils exposing (display)

import Svg exposing (..)
import Svg.Attributes exposing (..)
import TcTurtle exposing (Instruction(..), Program)

type alias Position =
    { x : Float, y : Float, angle : Float }

-- Execute a single instruction
execute : Position -> Instruction -> (Position, List (Svg msg))
execute pos instruction =
    case instruction of
        Forward n ->
            let
                newX = pos.x + toFloat n * cos (degrees pos.angle)
                newY = pos.y + toFloat n * sin (degrees pos.angle)
                lineSvg =
                    line
                        [ x1 (String.fromFloat pos.x), y1 (String.fromFloat pos.y)
                        , x2 (String.fromFloat newX), y2 (String.fromFloat newY)
                        , stroke "black"
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
                (newPos, svgs) =
                    List.foldl
                        (\_ (p, acc) ->
                            List.foldl
                                (\instr (pNext, accNext) ->
                                    let
                                        (nextPos, svg) =
                                            execute pNext instr
                                    in
                                    (nextPos, accNext ++ svg)
                                )
                                (p, acc)
                                instructions
                        )
                        (pos, [])
                        (List.repeat count ())
            in
            (newPos, svgs)

-- Convert complete program to SVG
display : Program -> Svg msg
display program =
    let
        initialPosition =
            { x = 0, y = 0, angle = 0 }

        (_, svgs) =
            List.foldl
                (\instr (pos, acc) ->
                    let
                        (newPos, svg) =
                            execute pos instr
                    in
                    (newPos, acc ++ svg)
                )
                (initialPosition, [])
                program
    in
    svg [ viewBox "0 0 500 500", width "500", height "500" ] svgs