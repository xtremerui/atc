module Navigation exposing (..)

import Html exposing (Html)
import Html.App
import Html.Attributes exposing (action, class, classList, href, id, method, style, title, disabled, attribute)
import Http
import Task

import Concourse.Job exposing (Job)
import Concourse.Pipeline exposing (Pipeline)

type alias CurrentState =
  { job : Maybe { name : String, teamName : String, pipelineName : String }
  }

type alias Model subModel =
  { subModel : subModel
  , currentState : CurrentState
  , pipelines : List Pipeline
  , jobs : List Job
  }

type Action subAction
  = SubAction subAction
  | PipelinesFetched (Result Http.Error (List Pipeline))
  -- | JobsFetched (Result Http.Error (List Job))

init : CurrentState -> (subModel, Cmd subAction) -> (Model subModel, Cmd (Action subAction))
init currentState (subModel, subCmd) =
  ( Model subModel currentState [] []
  , Cmd.batch
      [ Cmd.map SubAction subCmd
      , fetchPipelines
      ]
  )

update : (subAction -> subModel -> (subModel, Cmd subAction)) -> Action subAction -> Model subModel -> (Model subModel, Cmd (Action subAction))
update subUpdate action model =
  case action of
    SubAction subAction ->
      let
        (subModel, subCmd) = subUpdate subAction model.subModel
      in
        ({ model | subModel = subModel }, Cmd.map SubAction subCmd)

    PipelinesFetched (Err err) ->
      Debug.log ("failed to fetch pipelines: " ++ toString err) <|
        (model, Cmd.none)

    PipelinesFetched (Ok pipelines) ->
      ({ model | pipelines = pipelines }, Cmd.none)

view : (subModel -> Html subAction) -> Model subModel -> Html (Action subAction)
view subView model =
  Html.div [class "nav-page"] [
    Html.div [class "nav-content"] [
      Html.App.map SubAction (subView model.subModel)
    ],

    Html.div [class "nav-radar"] [
      failedBuild 1,
      failedBuild 8,
      failedBuild 54,
      failedBuild 99,
      failedBuild 182,
      startedBuild 34,
      startedBuild 68,
      startedBuild 96,
      startedBuild 110
    ]
  ]

failedBuild : Int -> Html (Action subAction)
failedBuild =
  buildEvent "failed"

startedBuild : Int -> Html (Action subAction)
startedBuild =
  buildEvent "started"

buildEvent : String -> Int -> Html (Action subAction)
buildEvent status pct =
  Html.div [
    classList [
      ("event", True),
      ("build", True),
      (status, True)
    ]
  ] [
    Html.div [class "event-header"] [
      Html.span [class "build"] [Html.text "#306"],
      Html.span [class "hl"] [Html.text "main"],
      Html.text "/",
      Html.span [class "hl"] [Html.text "testflight"]
    ],

    Html.div [class "progress"] [
      Html.span [class "fixed-duration"] [Html.text (toString (round <| 15 * (toFloat pct / 100.0)) ++ "m")],

      Html.div [class "fg", style [("width", toString (min pct 100) ++ "%")]] [
        Html.span [class "glyph"] [
          -- Html.span [class "embedded-duration"] [Html.text (toString (round <| 15 * (toFloat pct / 100.0)) ++ "m")],
          Html.i [class ("fa " ++ iconForStatus status)] []
        ]
      ],

      if pct > 100 then
        Html.div [class "over", style [("width", toString (pct - 100) ++ "%")]] []
      else
        Html.div [class "not-over-lol"] []
    ]
  ]

iconForStatus : String -> String
iconForStatus status =
  case status of
    "failed" ->
      "fa-close"

    "started" ->
      "fa-plane"

    _ ->
      "fa-beer"

fetchPipelines : Cmd (Action subAction)
fetchPipelines =
  Cmd.map PipelinesFetched << Task.perform Err Ok <|
    Concourse.Pipeline.fetchAll "main"
