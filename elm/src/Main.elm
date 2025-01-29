module Main exposing (..)

{-| This is the main module of the application containing the entry point of the project.

It manages:

  - The overall app view.
  - The state of the model (such as user input, parsed code, and customizations).
  - User interaction (parsing code, changing colors, zooming, etc.).

-}

import Browser
import DrawingUtils exposing (Color, display)
import Html exposing (Html, button, div, h3, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import ParsingUtils exposing (parseErrorToString, read)
import Svg exposing (Svg, svg)
import Svg.Attributes
import TcTurtle exposing (Program)



-- MODEL


{-| Defines the `Model` for the main application state.

Fields:

  - `code` : String - The TcTurtle code entered by the user.
  - `error` : Maybe String - Any error messages displayed during parsing.
  - `svg` : Svg msg - The rendered SVG output based on parsed code.
  - `color` : Color - The current drawing color.
  - `program` : Maybe TcTurtle.Program - The parsed drawing instructions.
  - `customRed`, `customGreen`, `customBlue`, `customAlpha`: Strings - User-entered custom color values.
  - `zoom` : Float - The zoom level for SVG rendering.

-}
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


{-| The initial state (model) of the application.
-}
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


{-| Messages (`Msg`) define the possible user interactions within the app.

Variants:

  - `UpdateCode`: Fired when the user updates the input TcTurtle code.
  - `ParseCode`: Fired when the "Parse & Draw" button is clicked.
  - `ChangeColor`: Fired when the user chooses a predefined color button.
  - `UpdateCustomColor`: Fired when the custom color input fields are updated.
  - `SetCustomColor`: Fired when the "Set Custom Color" button is clicked.
  - `ZoomIn`: Fired to zoom into the SVG.
  - `ZoomOut`: Fired to zoom out of the SVG.

-}
type Msg
    = UpdateCode String
    | ParseCode
    | ChangeColor Color
    | UpdateCustomColor String String String String
    | SetCustomColor
    | ZoomIn
    | ZoomOut


{-| Handles app state updates based on `Msg`.

Parameters:

  - `msg`: The message triggered by user interaction.
  - `model`: The current application state.

Returns:

  - An updated version of the `Model`.

-}
update : Msg -> Model Never -> Model Never
update msg model =
    case msg of
        -- Handle updating the code entered by the user
        UpdateCode newCode ->
            { model | code = newCode }

        -- Handle the "Parse & Draw" event
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

        -- Handle changing the drawing color
        ChangeColor newColor ->
            case model.program of
                Just program ->
                    { model
                        | color = newColor
                        , svg = display program model.zoom newColor
                    }

                Nothing ->
                    { model | color = newColor }

        -- Update only the color if no program exists
        -- Handle updating values in the custom color input fields
        UpdateCustomColor red green blue alpha ->
            { model
                | customRed = red
                , customGreen = green
                , customBlue = blue
                , customAlpha = alpha
            }

        -- Handle applying the custom color to the drawing
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

        -- Handle zooming in the SVG
        ZoomIn ->
            case model.program of
                Just program ->
                    { model
                        | zoom = model.zoom * 1.1
                        , svg = display program (model.zoom * 1.1) model.color
                    }

                Nothing ->
                    { model | zoom = model.zoom * 1.1 }

        -- Handle zooming out the SVG
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


{-| Constructs the HTML structure of the app.

Parameters:

  - `model`: The current application state.

Returns:

  - The view (HTML) as a `Html Msg`.

-}
view : Model Never -> Html Msg
view model =
    div [ Html.Attributes.class "page" ]
        [ -- Input box for TcTurtle code
          input
            [ placeholder "Enter TcTurtle code"
            , Html.Attributes.class "input"
            , onInput UpdateCode
            , value model.code
            ]
            []
        , -- "Parse & Draw" Button
          button
            [ Html.Attributes.class "button"
            , onClick ParseCode
            ]
            [ text "Parse & Draw" ]
        , -- Color selection area
          div [ Html.Attributes.class "color-section" ]
            [ h3 [ Html.Attributes.class "color-title small-font" ] [ text "Choose a color for your pencil:" ]
            , -- Predefined color buttons
              div [ Html.Attributes.class "color-buttons" ]
                [ button
                    [ Html.Attributes.class "color-button black"
                    , onClick (ChangeColor { red = 0, green = 0, blue = 0, alpha = 1.0 })
                    ]
                    [ text "Black" ]
                , button
                    [ Html.Attributes.class "color-button red"
                    , onClick (ChangeColor { red = 255, green = 0, blue = 0, alpha = 1.0 })
                    ]
                    [ text "Red" ]
                , button
                    [ Html.Attributes.class "color-button green"
                    , onClick (ChangeColor { red = 0, green = 255, blue = 0, alpha = 1.0 })
                    ]
                    [ text "Green" ]
                , button
                    [ Html.Attributes.class "color-button blue"
                    , onClick (ChangeColor { red = 0, green = 0, blue = 255, alpha = 1.0 })
                    ]
                    [ text "Blue" ]
                ]
            , -- Custom color selection area
              div [ Html.Attributes.class "custom-color" ]
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
                , div [ Html.Attributes.class "center-button" ]
                    [ button
                        [ Html.Attributes.class "button"
                        , onClick SetCustomColor
                        ]
                        [ text "Set Custom Color" ]
                    ]
                ]
            ]
        , -- Zoom controls
          div [ Html.Attributes.class "zoom-controls" ]
            [ text "Zoom Controls: "
            , button [ Html.Attributes.class "zoom-button", onClick ZoomIn ] [ text "+" ]
            , button [ Html.Attributes.class "zoom-button", onClick ZoomOut ] [ text "−" ]
            ]
        , -- SVG output or error message
          case model.error of
            Nothing ->
                div [ Html.Attributes.class "svg" ] [ Html.map (always ParseCode) model.svg ]

            Just err ->
                div [ Html.Attributes.class "error" ] [ text err ]
        ]



-- PROGRAM


{-| The entry point for the program using `Browser.sandbox`.

Defines the main logic of the app by initializing, updating, and rendering the view.

-}
main =
    Browser.sandbox { init = init, update = update, view = view }
