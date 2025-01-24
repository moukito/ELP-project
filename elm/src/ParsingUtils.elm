module ParsingUtils exposing (read, programParser)

import Parser exposing (..)

-- Parse individual instructions
instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ map Forward (token "Forward" |> followedBy spaces |> andThen int)
        , map Left (token "Left" |> followedBy spaces |> andThen int)
        , map Right (token "Right" |> followedBy spaces |> andThen int)
        , repeatParser
        ]

repeatParser : Parser Instruction
repeatParser =
    map2 Repeat
        (token "Repeat" |> followedBy spaces |> andThen int)
        (inBrackets (separatedBy (symbol ",") instructionParser))

-- Program parser: Parse entire program
programParser : Parser Program
programParser =
    inBrackets (separatedBy (symbol ",") instructionParser)

-- Entry point for parsing
read : String -> Result (List Parser.DeadEnd) Program
read input =
    run programParser input