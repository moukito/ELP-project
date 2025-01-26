module ParsingUtils exposing (programParser, read)

import Parser exposing (..)
import TcTurtle exposing (..)



-- Parse individual instructions


forwardParser : Parser Instruction
forwardParser =
    succeed Forward
        |= (symbol "Forward" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


leftParser : Parser Instruction
leftParser =
    succeed Left
        |= (symbol "Left" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


rightParser : Parser Instruction
rightParser =
    succeed Right
        |= (symbol "Right" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


repeatParser : Parser Instruction
repeatParser =
    succeed Repeat
        |= (symbol "Repeat" -- Match "Repeat"
                |> andThen (\_ -> spaces) -- Skip spaces after "Repeat"
                |> andThen (\_ -> int) -- Parse the repeat count
           )
        |= (spaces
                |> andThen (\_ ->
                    programParser -- Parse the inner program
                )
           )


instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ forwardParser
        , leftParser
        , rightParser
        , repeatParser
        ]



-- Program parser: Parse entire program


programParser : Parser (List Instruction)
programParser =
    Parser.sequence
        { start = "["
        , separator = ","
        , end = "]"
        , spaces = spaces
        , item = instructionParser
        , trailing = Parser.Optional
        }

-- Entry point for parsing


read : String -> Result (List Parser.DeadEnd) (List Instruction)
read input =
    run programParser input
