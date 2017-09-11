module BuildDuration exposing (view, viewFailDuration)

import Date exposing (Date)
import Date.Format
import Duration exposing (Duration)
import Html exposing (Html)
import Html.Attributes exposing (class, title)
import Time exposing (Time)
import Concourse


view : Concourse.BuildDuration -> Time.Time -> Html a
view duration now =
    Html.table [ class "dictionary build-duration" ] <|
        case ( duration.startedAt, duration.finishedAt ) of
            ( Nothing, Nothing ) ->
                [ pendingLabel "pending" ]

            ( Just startedAt, Nothing ) ->
                [ labeledRelativeDate "started" now startedAt ]

            ( Nothing, Just finishedAt ) ->
                [ labeledRelativeDate "finished" now finishedAt ]

            ( Just startedAt, Just finishedAt ) ->
                let
                    durationElmIssue =
                        -- https://github.com/elm-lang/elm-compiler/issues/1223
                        Duration.between (Date.toTime startedAt) (Date.toTime finishedAt)
                in
                    [ labeledRelativeDate "started" now startedAt
                    , labeledRelativeDate "finished" now finishedAt
                    , labeledDuration "duration" durationElmIssue
                    ]


viewFailDuration : Concourse.BuildDuration -> Time.Time -> Html a
viewFailDuration duration time =
    case duration.startedAt of
        Nothing ->
            Html.div [ class "build-duration" ] []

        Just startedAt ->
            let
                now =
                    time

                oneSecondLater =
                    Time.inMilliseconds time + 1000 * Time.millisecond

                durationUntilNow =
                    Duration.between (Date.toTime startedAt) now

                durationUntilOneSecondLater =
                    Duration.between (Date.toTime startedAt) oneSecondLater

                seconds =
                    truncate (durationUntilNow / 1000)

                remainingSeconds =
                    rem seconds 60

                minutes =
                    seconds // 60

                remainingMinutes =
                    rem minutes 60

                hours =
                    minutes // 60

                remainingHours =
                    rem hours 24

                days =
                    hours // 24

                secondsFoo =
                    truncate (durationUntilOneSecondLater / 1000)

                remainingSecondsFoo =
                    rem secondsFoo 60

                minutesFoo =
                    secondsFoo // 60

                remainingMinutesFoo =
                    rem minutesFoo 60

                hoursFoo =
                    minutesFoo // 60

                remainingHoursFoo =
                    rem hoursFoo 24

                daysFoo =
                    hoursFoo // 24

                da =
                    Html.td []
                        [ Html.span
                            [ if days /= daysFoo then
                                class "flip-next"
                              else
                                class ""
                            ]
                            [ Html.text <| toString days ]
                        , Html.span
                            [ if days /= daysFoo then
                                class "flip-back"
                              else
                                class ""
                            ]
                            [ Html.text <| toString days ]
                        , Html.span
                            [ if days /= daysFoo then
                                class "flip-top"
                              else
                                class ""
                            ]
                            [ Html.text <| toString daysFoo ]
                        ]

                ho =
                    Html.td []
                        [ Html.span
                            [ if remainingHours /= remainingHoursFoo then
                                class "flip"
                              else
                                class ""
                            ]
                            [ Html.text <| toString remainingHours ]
                        , Html.span
                            []
                            [ Html.text "h" ]
                        ]

                mi =
                    Html.td []
                        [ Html.span
                            [ if remainingMinutes /= remainingMinutesFoo then
                                class "flip"
                              else
                                class ""
                            ]
                            [ Html.text <| toString remainingMinutes ]
                        , Html.span
                            []
                            [ Html.text "m" ]
                        ]

                se =
                    Html.td []
                        [ Html.div [ class "flip-wrapper" ]
                            [ Html.div
                                [ if remainingSeconds /= remainingSecondsFoo then
                                    class "flip flip-next"
                                  else
                                    class ""
                                ]
                                [ Html.text <| toString remainingSecondsFoo ]
                            , Html.div
                                [ if remainingSeconds /= remainingSecondsFoo then
                                    class "flip flip-top"
                                  else
                                    class ""
                                ]
                                [ Html.text <| toString remainingSeconds ]
                            , Html.div
                                [ if remainingSeconds /= remainingSecondsFoo then
                                    class "flip flip-back"
                                  else
                                    class ""
                                ]
                                [ Html.text <| toString remainingSecondsFoo ]
                            , Html.div
                                [ if remainingSeconds /= remainingSecondsFoo then
                                    class "flip flip-bottom"
                                  else
                                    class ""
                                ]
                                [ Html.text <| toString remainingSeconds ]
                            ]
                        , Html.div
                            []
                            [ Html.text "s" ]
                        ]
            in
                case ( days, remainingHours, remainingMinutes, remainingSeconds ) of
                    ( 0, 0, 0, s ) ->
                        Html.div [ class "build-duration" ]
                            [ Html.tr []
                                [ Html.td [] [ Html.text "failing for:" ]
                                , se
                                ]
                            ]

                    ( 0, 0, m, s ) ->
                        viewBuildDuration mi se

                    ( 0, h, m, _ ) ->
                        viewBuildDuration ho mi

                    ( d, h, _, _ ) ->
                        viewBuildDuration da ho


viewBuildDuration : Html msg -> Html msg -> Html msg
viewBuildDuration first second =
    Html.div [ class "build-duration" ]
        [ Html.tr []
            [ Html.td [] [ Html.text "failing for:" ]
            , first
            , second
            ]
        ]


labeledRelativeDate : String -> Time -> Date -> Html a
labeledRelativeDate label now date =
    let
        ago =
            Duration.between (Date.toTime date) now
    in
        Html.tr []
            [ Html.td [ class "dict-key" ] [ Html.text label ]
            , Html.td
                [ title (Date.Format.format "%b %d %Y %I:%M:%S %p" date), class "dict-value" ]
                [ Html.span [] [ Html.text (Duration.format ago ++ " ago") ] ]
            ]


labeledDuration : String -> Duration -> Html a
labeledDuration label duration =
    Html.tr []
        [ Html.td [ class "dict-key" ] [ Html.text label ]
        , Html.td [ class "dict-value" ] [ Html.span [] [ Html.text (Duration.format duration) ] ]
        ]


pendingLabel : String -> Html a
pendingLabel label =
    Html.tr []
        [ Html.td [ class "dict-key" ] [ Html.text label ]
        , Html.td [ class "dict-value" ] []
        ]
