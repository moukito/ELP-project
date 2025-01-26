module Main exposing (..)

import Browser
import DrawingUtils exposing (display)
import Html exposing (Html, button, div, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import ParsingUtils exposing (read)
import Svg exposing (Svg, svg)



-- MODEL


type alias Model msg =
    { code : String
    , error : Maybe String
    , svg : Svg msg
    }


init : Model Never
init =
    { code = ""
    , error = Nothing
    , svg = svg [] []
    }



-- UPDATE


type Msg
    = UpdateCode String
    | ParseCode


update : Msg -> Model Never -> Model Never
update msg model =
    case msg of
        UpdateCode newCode ->
            { model | code = newCode }

        ParseCode ->
            case read model.code of
                Ok program ->
                    { model | error = Nothing, svg = display program }

                Err problems ->
                    { model
                        | error = Just "Syntax Error!"
                        , svg = svg [] []
                    }



-- VIEW


view : Model Never -> Html Msg
view model =
    div []
        [ input [ placeholder "Enter TcTurtle code", onInput UpdateCode, value model.code ] []
        , button [ onClick ParseCode ] [ text "Parse & Draw" ]
        , case model.error of
            Nothing ->
                text ""

            Just err ->
                div [] [ text err ]
        --, Html.map (always ()) (Html.fromSvg model.svg)
        ]



-- PROGRAM


main =
    Browser.sandbox { init = init, update = update, view = view }
