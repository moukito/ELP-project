module ParsingUtils exposing (parseErrorToString, programParser, read)

{-| This module provides utilities for parsing TcTurtle programs, converting user input strings
to structured `Program` data, and handling errors during the parsing process.

Main functions:

  - `programParser`: Parses a list of turtle instructions (`Program`).
  - `read`: Entry point for parsing a string into a `Program`.
  - `parseErrorToString`: Converts parsing errors to human-readable messages.

Parsers for individual instructions are also included (e.g., `forwardParser`, `rightParser`, etc.).

-}

import Parser exposing (..)
import TcTurtle exposing (..)



-- PARSERS FOR INSTRUCTIONS


{-| Parses the `Forward` instruction, which moves the turtle forward by a specified number of units.
Example syntax: `Forward 100`
-}
forwardParser : Parser Instruction
forwardParser =
    succeed Forward
        |= (symbol "Forward" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


{-| Parses the `Left` instruction, which turns the turtle left by a specified angle in degrees.
Example syntax: `Left 90`
-}
leftParser : Parser Instruction
leftParser =
    succeed Left
        |= (symbol "Left" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


{-| Parses the `Right` instruction, which turns the turtle right by a specified angle in degrees.
Example syntax: `Right 90`
-}
rightParser : Parser Instruction
rightParser =
    succeed Right
        |= (symbol "Right" |> andThen (\_ -> spaces) |> andThen (\_ -> int))


{-| Parses the `Repeat` instruction, which repeats a given sequence of instructions a specified number of times.
Example syntax: `Repeat 4 [ Forward 100, Left 90 ]`
-}
repeatParser : Parser Instruction
repeatParser =
    succeed Repeat
        |= (symbol "Repeat"
                |> andThen (\_ -> spaces)
                |> andThen (\_ -> int)
           )
        |= (spaces
                |> andThen
                    (\_ ->
                        programParser
                    )
           )


{-| Parses any individual turtle instruction (`Forward`, `Left`, `Right`, or `Repeat`).
-}
instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ forwardParser
        , leftParser
        , rightParser
        , repeatParser
        ]



-- PARSER FOR ENTIRE PROGRAM


{-| Parses an entire list of instructions enclosed in brackets.
Example syntax: `[ Forward 100, Left 90, Forward 50 ]`
-}
programParser : Parser TcTurtle.Program
programParser =
    Parser.sequence
        { start = "["
        , separator = ","
        , end = "]"
        , spaces = spaces
        , item = instructionParser
        , trailing = Parser.Optional
        }



-- ENTRY POINT


{-| Parses the user input string into a `Program`.

Parameters:

  - `input`: User input string.

Returns:

  - A `Result` containing either the parsed program or a list of parsing errors.

-}
read : String -> Result (List Parser.DeadEnd) TcTurtle.Program
read input =
    run programParser input



-- ERROR HANDLING


{-| Converts a list of parsing errors into a human-readable message.

Parameters:

  - `errors`: A list of `Parser.DeadEnd` values representing parsing errors.

Returns:

  - A string describing the first error in the list.

-}
parseErrorToString : List Parser.DeadEnd -> String
parseErrorToString errors =
    case errors of
        [] ->
            "An unknown error occurred while parsing."

        firstError :: _ ->
            case firstError.problem of
                Parser.Expecting _ ->
                    "It seems like you forgot to close a bracket. Make sure all brackets match: [ ... ]"

                Parser.ExpectingSymbol _ ->
                    "Expected Forward, Left, Right or Repeat. Please check your syntax."

                Parser.ExpectingInt ->
                    "Expected an integer. Please check your syntax."

                Parser.BadRepeat ->
                    "It looks like there is an invalid format in a repeat block."

                _ ->
                    "Invalid syntax. Please check your program for errors."
