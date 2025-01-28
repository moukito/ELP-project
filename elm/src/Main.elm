module Main exposing (..)

import Browser
import DrawingUtils exposing (display)
import Html exposing (Html, button, div, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import ParsingUtils exposing (read, parseErrorToString)
import Svg exposing (Svg, svg)
import Svg.Attributes



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
    , svg = svg [ Svg.Attributes.viewBox "0 0 500 500", Svg.Attributes.width "500", Svg.Attributes.height "500" ] []
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
                        | error = Just (parseErrorToString problems) -- "Syntax Error!"
                        , svg = svg [] []
                    }



-- VIEW


view : Model Never -> Html Msg
view model =
    div [ Html.Attributes.class "page" ]
        [ input [ placeholder "Enter TcTurtle code", Html.Attributes.class "input", onInput UpdateCode, value model.code ] []
        , button [ Html.Attributes.class "button", onClick ParseCode ] [ text "Parse & Draw" ]
        , case model.error of
            Nothing ->
                div [ Html.Attributes.class "svg" ] [ Html.map (always ParseCode) model.svg ] -- Wrap here using Html.map


            Just err ->
                div [ Html.Attributes.class "error" ] [ text err ]

        ]



-- PROGRAM


main =
    Browser.sandbox { init = init, update = update, view = view }
