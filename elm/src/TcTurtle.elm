module TcTurtle exposing (..)

{-|
This module defines the core types and structures used to describe a `TcTurtle` drawing program.

Types:
  - `Instruction`: Represents instructions that a turtle can execute (e.g., moving forward, turning).
  - `Program`: A list of instructions to be interpreted and executed.
-}


{-|
The `Instruction` type describes individual turtle commands.

Variants:
  - `Forward Int`: Move forward by a given number of units.
  - `Left Int`: Turn left by a specified angle in degrees.
  - `Right Int`: Turn right by a specified angle in degrees.
  - `Repeat Int (List Instruction)`: Repeat a given number of times a sequence of instructions.
-}
type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)


{-|
The `Program` type is a list of `Instruction`s that define the TcTurtle program.
-}
type alias Program =
    List Instruction