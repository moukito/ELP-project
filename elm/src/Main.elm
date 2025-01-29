module Main exposing (..)

import Browser
import DrawingUtils exposing (display)
import Html exposing (Html, button, div, input, text, h3)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import ParsingUtils exposing (parseErrorToString, read)
import Svg exposing (Svg, svg)
import Svg.Attributes
import TcTurtle exposing (Program)


-- MODEL

type alias Color =
    { red : Int, green : Int, blue : Int, alpha : Float }


type alias Model msg =
    { code : String
    , error : Maybe String
    , svg : Svg msg
    , color : Color
    , program : Maybe Program
    , customRed : String
    , customGreen : String
    , customBlue : String
    , customAlpha : String
    , zoom : Float
    }


init : Model Never
init =
    { code = ""
    , error = Nothing
    , svg = svg [ Svg.Attributes.viewBox "0 0 500 500", Svg.Attributes.width "500", Svg.Attributes.height "500" ] []
    , color = { red = 0, green = 0, blue = 0, alpha = 1.0 } -- Default black
    , program = Nothing
    , customRed = "0"
    , customGreen = "0"
    , customBlue = "0"
    , customAlpha = "1.0"
    , zoom = 1.0
    }


-- UPDATE

type Msg
    = UpdateCode String
    | ParseCode
    | ChangeColor Color
    | UpdateCustomColor String String String String
    | SetCustomColor
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
                    { model
                        | error = Nothing
                        , svg = display program model.zoom model.color
                        , program = Just program
                    }

                Err problems ->
                    { model
                        | error = Just (parseErrorToString problems)
                        , svg = svg [] []
                        , program = Nothing
                    }

        ChangeColor newColor ->
            case model.program of
                Just program ->
                    { model
                        | color = newColor
                        , svg = display program model.zoom newColor
                    }

                Nothing ->
                    { model | color = newColor }

        UpdateCustomColor red green blue alpha ->
            { model
                | customRed = red
                , customGreen = green
                , customBlue = blue
                , customAlpha = alpha
            }

        SetCustomColor ->
            case model.program of
                Just program ->
                    let
                        newColor =
                            { red = String.toInt model.customRed |> Maybe.withDefault 0
                            , green = String.toInt model.customGreen |> Maybe.withDefault 0
                            , blue = String.toInt model.customBlue |> Maybe.withDefault 0
                            , alpha = String.toFloat model.customAlpha |> Maybe.withDefault 1.0
                            }
                    in
                    { model
                        | color = newColor
                        , svg = display program model.zoom newColor
                    }

                Nothing ->
                    let
                        newColor =
                            { red = String.toInt model.customRed |> Maybe.withDefault 0
                            , green = String.toInt model.customGreen |> Maybe.withDefault 0
                            , blue = String.toInt model.customBlue |> Maybe.withDefault 0
                            , alpha = String.toFloat model.customAlpha |> Maybe.withDefault 1.0
                            }
                    in
                    { model | color = newColor }

        ZoomIn ->
            case model.program of
                Just program ->
                    { model
                        | zoom = model.zoom * 1.1
                        , svg = display program (model.zoom * 1.1) model.color
                    }

                Nothing ->
                    { model | zoom = model.zoom * 1.1 }

        ZoomOut ->
            case model.program of
                Just program ->
                    { model
                        | zoom = model.zoom * 0.9
                        , svg = display program (model.zoom * 0.9) model.color
                    }

                Nothing ->
                    { model | zoom = model.zoom * 0.9 }



-- VIEW

view : Model Never -> Html Msg
view model =
    div [ Html.Attributes.class "page" ]
        [ input
            [ placeholder "Enter TcTurtle code"
            , Html.Attributes.class "input"
            , onInput UpdateCode
            , value model.code
            ]
            []
        , button
            [ Html.Attributes.class "button"
            , onClick ParseCode
            ]
            [ text "Parse & Draw" ]
        , div [ Html.Attributes.class "color-section" ]
            [ h3 [ Html.Attributes.class "color-title small-font" ] [ text "Choose a color for your pencil:" ]
            , div [ Html.Attributes.class "color-buttons" ]
                [ button
                    [ onClick (ChangeColor { red = 0, green = 0, blue = 0, alpha = 1.0 }) ]
                    [ text "Black" ]
                , button
                    [ onClick (ChangeColor { red = 255, green = 0, blue = 0, alpha = 1.0 }) ]
                    [ text "Red" ]
                , button
                    [ onClick (ChangeColor { red = 0, green = 255, blue = 0, alpha = 1.0 }) ]
                    [ text "Green" ]
                , button
                    [ onClick (ChangeColor { red = 0, green = 0, blue = 255, alpha = 1.0 }) ]
                    [ text "Blue" ]
                ]
            , div [ Html.Attributes.class "custom-color" ]
                [ text "Or decide the color you want:"
                , div [ Html.Attributes.class "color-inputs" ]
                [ input
                    [ placeholder "Red (0-255)"
                    , value model.customRed
                    , onInput (\v -> UpdateCustomColor v model.customGreen model.customBlue model.customAlpha)
                    , Html.Attributes.class "small-input"
                    ]
                    []
                , input
                    [ placeholder "Green (0-255)"
                    , value model.customGreen
                    , onInput (\v -> UpdateCustomColor model.customRed v model.customBlue model.customAlpha)
                    , Html.Attributes.class "small-input"
                    ]
                    []
                , input
                    [ placeholder "Blue (0-255)"
                    , value model.customBlue
                    , onInput (\v -> UpdateCustomColor model.customRed model.customGreen v model.customAlpha)
                    , Html.Attributes.class "small-input"
                    ]
                    []
                , input
                    [ placeholder "Alpha (0.0-1.0)"
                    , value model.customAlpha
                    , onInput (\v -> UpdateCustomColor model.customRed model.customGreen model.customBlue v)
                    , Html.Attributes.class "small-input"
                    ]
                    []
                ]   
                , button
                    [ Html.Attributes.class "button"
                    , onClick SetCustomColor
                    ]
                    [ text "Set Custom Color" ]
                ]
            ]
        , div [ Html.Attributes.class "zoom-controls" ]
            [ text "Zoom Controls: "
            , button [ Html.Attributes.class "zoom-button", onClick ZoomIn ] [ text "+" ]
            , button [ Html.Attributes.class "zoom-button", onClick ZoomOut ] [ text "−" ]
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

