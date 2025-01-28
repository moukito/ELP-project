module TcTurtle exposing (..)


type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)


type alias Program =
    List Instruction
