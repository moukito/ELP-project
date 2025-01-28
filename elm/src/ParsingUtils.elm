module ParsingUtils exposing (parseErrorToString, programParser, read)

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


instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ forwardParser
        , leftParser
        , rightParser
        , repeatParser
        ]



-- Program parser: Parse entire program


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



-- Entry point for parsing


read : String -> Result (List Parser.DeadEnd) TcTurtle.Program
read input =
    run programParser input



--


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



-- Exemple de test
-- Fonction pour exécuter une commande


executeCommand : Instruction -> String
executeCommand command =
    case command of
        Forward distance ->
            "Executing Forward with distance: " ++ String.fromInt distance

        Left angle ->
            "Executing Left with angle: " ++ String.fromInt angle

        Right angle ->
            "Executing Right with angle: " ++ String.fromInt angle

        Repeat count commands ->
            String.join "\n"
                (List.concatMap (\_ -> List.map executeCommand commands) (List.repeat count ()))



-- Fonction pour exécuter un programme


executeProgram : TcTurtle.Program -> String
executeProgram program =
    String.join "\n" (List.map executeCommand program)



-- Fonction pour analyser et exécuter le programme


processInput : String -> String
processInput input =
    case run programParser input of
        Ok program ->
            executeProgram program

        Err _ ->
            "Invalid program"
