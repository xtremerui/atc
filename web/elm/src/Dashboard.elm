port module Dashboard exposing (Model, view)

import Concourse
import Html exposing (Html)
import List.Extra exposing (groupWhile)


type alias Model =
    { pipelines : (List Concourse.Pipeline)
    }

view : Model -> Html msg
view model =
    let
        groupedPipelines = groupWhile
            (\pipeline1 pipeline2 -> pipeline1.teamName == pipeline2.teamName)
            model.pipelines
        teamNames = List.map (\pipelines ->
            case pipelines of
                p :: _ -> p.teamName
                _ -> ""
            ) groupedPipelines
    in
        Html.div [] (List.map Html.text teamNames)
