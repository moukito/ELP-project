module TcTurtle exposing (..)

type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)

type Program = List Instruction