module Main exposing (..)

import Browser
import DrawingUtils exposing (display)
import Html exposing (Html, button, div, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import ParsingUtils exposing (parseErrorToString, read)
import Svg exposing (Svg, svg)
import Svg.Attributes



-- MODEL


type alias Model msg =
    { code : String
    , error : Maybe String
    , svg : Svg msg
    , zoom : Float
    , prog : Float -> Svg msg
    }


init : Model Never
init =
    { code = ""
    , error = Nothing
    , svg = svg [ Svg.Attributes.viewBox "0 0 500 500", Svg.Attributes.width "500", Svg.Attributes.height "500" ] []
    , zoom = 1.0
    , prog = display []
    }



-- UPDATE


type Msg
    = UpdateCode String
    | ParseCode
    | ZoomIn
    | ZoomOut


update : Msg -> Model Never -> Model Never
update msg model =
    case msg of
        UpdateCode newCode ->
            { model | code = newCode }

        ParseCode ->
            case read model.code of
                Ok program ->
                    { model | error = Nothing, prog = display program, svg = display program model.zoom }

                Err problems ->
                    { model
                        | error = Just (parseErrorToString problems) -- "Syntax Error!"
                        , svg = svg [] []
                    }

        ZoomIn ->
            { model | error = Nothing, zoom = model.zoom * 1.1, svg = model.prog (model.zoom * 1.1) }

        ZoomOut ->
            { model | error = Nothing, zoom = model.zoom * 0.9, svg = model.prog (model.zoom * 0.9) }



-- VIEW


view : Model Never -> Html Msg
view model =
    div [ Html.Attributes.class "page" ]
        [ input [ placeholder "Enter TcTurtle code", Html.Attributes.class "input", onInput UpdateCode, value model.code ] []
        , button [ Html.Attributes.class "button", onClick ParseCode ] [ text "Parse & Draw" ]
        , div [ Html.Attributes.class "zoom-controls" ]
            [ text "Zoom Controls: "
            , button [ Html.Attributes.class "zoom-button", onClick ZoomIn ] [ text "+" ] -- Use "+" for Zoom In
            , button [ Html.Attributes.class "zoom-button", onClick ZoomOut ] [ text "−" ] -- Use "−" for Zoom Out
            ]
        , case model.error of
            Nothing ->
                div [ Html.Attributes.class "svg" ] [ Html.map (always ParseCode) model.svg ]

            Just err ->
                div [ Html.Attributes.class "error" ] [ text err ]
        ]



-- PROGRAM


main =
    Browser.sandbox { init = init, update = update, view = view }
