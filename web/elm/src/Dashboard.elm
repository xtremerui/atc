port module Dashboard exposing (Model, Msg, init, update, subscriptions, view)

import BuildDuration
import Concourse
import Concourse.BuildStatus
import Concourse.Job
import Concourse.Pipeline
import Date exposing (Date)
import Dict exposing (Dict)
import Graph exposing (Graph)
import Grid exposing (Grid)
import Html exposing (Html)
import Html.Attributes exposing (class, classList, id, href, src)
import RemoteData
import Task exposing (Task)
import Time exposing (Time)


type alias Model =
    { pipelines : RemoteData.WebData (List Concourse.Pipeline)
    , jobs : Dict Int (RemoteData.WebData (List Concourse.Job))
    , now : Maybe Time
    , turbulenceImgSrc : String
    }


type Msg
    = PipelinesResponse (RemoteData.WebData (List Concourse.Pipeline))
    | JobsResponse Int (RemoteData.WebData (List Concourse.Job))
    | ClockTick Time.Time
    | AutoRefresh Time


type Node
    = JobNode Concourse.Job
    | ConstrainedInputNode
        { resourceName : String
        , dependentJob : Concourse.Job
        , upstreamJob : Maybe Concourse.Job
        }


type alias ByName a =
    Dict String a


init : String -> ( Model, Cmd Msg )
init turbulencePath =
    ( { pipelines = RemoteData.NotAsked
      , jobs = Dict.empty
      , now = Nothing
      , turbulenceImgSrc = turbulencePath
      }
    , Cmd.batch [ fetchPipelines, getCurrentTime ]
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        PipelinesResponse response ->
            ( { model | pipelines = response }
            , case response of
                RemoteData.Success pipelines ->
                    Cmd.batch (List.map fetchJobs pipelines)

                _ ->
                    Cmd.none
            )

        JobsResponse pipelineId response ->
            ( { model | jobs = Dict.insert pipelineId response model.jobs }, Cmd.none )

        ClockTick now ->
            ( { model | now = Just now }, Cmd.none )

        AutoRefresh _ ->
            ( model, fetchPipelines )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Time.every Time.second ClockTick
        , Time.every (5 * Time.second) AutoRefresh
        ]


view : Model -> Html msg
view model =
    case model.pipelines of
        RemoteData.Success pipelines ->
            let
                pipelineStates =
                    List.sortBy pipelineStatusRank <|
                        List.filter ((/=) RemoteData.NotAsked << .jobs) <|
                            List.map
                                (\pipeline ->
                                    { pipeline = pipeline
                                    , jobs =
                                        Maybe.withDefault RemoteData.NotAsked <|
                                            Dict.get pipeline.id model.jobs
                                    }
                                )
                                pipelines

                pipelinesByTeam =
                    List.foldl
                        (\pipelineState byTeam ->
                            Dict.update pipelineState.pipeline.teamName
                                (\mps ->
                                    Just (pipelineState :: Maybe.withDefault [] mps)
                                )
                                byTeam
                        )
                        Dict.empty
                        (List.reverse pipelineStates)
            in
                Html.div [ class "dashboard" ]
                    (Dict.values (Dict.map (viewGroup model.now) pipelinesByTeam))

        RemoteData.Failure _ ->
            Html.div
                [ class "error-message" ]
                [ Html.div [ class "message" ]
                    [ Html.img [ src model.turbulenceImgSrc, class "seatbelt" ] []
                    , Html.p [] [ Html.text "experiencing turbulence" ]
                    , Html.p [ class "explanation" ] []
                    ]
                ]

        _ ->
            Html.text ""


pipelineStatusRank : PipelineState -> Int
pipelineStatusRank state =
    let
        order =
            [ Concourse.BuildStatusFailed, Concourse.BuildStatusErrored, Concourse.BuildStatusAborted, Concourse.BuildStatusSucceeded, Concourse.BuildStatusPending ]

        ranks =
            List.indexedMap
                (\index buildStatus ->
                    ( Concourse.BuildStatus.show buildStatus, index * 2 )
                )
                order

        unranked =
            (List.length ranks) * 2

        status =
            pipelineStatus state

        paused =
            state.pipeline.paused

        running =
            isPipelineRunning state

        rank =
            Maybe.withDefault unranked <| Dict.get status <| Dict.fromList ranks
    in
        if paused then
            unranked
        else if running then
            rank - 1
        else
            rank


type alias PipelineState =
    { pipeline : Concourse.Pipeline
    , jobs : RemoteData.WebData (List Concourse.Job)
    }


viewGroup : Maybe Time -> String -> List PipelineState -> Html msg
viewGroup now teamName pipelines =
    Html.div [ id teamName, class "dashboard-team-group" ]
        [ Html.div [ class "dashboard-team-name" ]
            [ Html.text teamName ]
        , Html.div [ class "dashboard-team-pipelines" ]
            (List.map (viewPipeline now) pipelines)
        ]


viewPipeline : Maybe Time -> PipelineState -> Html msg
viewPipeline now state =
    Html.div
        [ classList
            [ ( "dashboard-pipeline", True )
            , ( "dashboard-paused", state.pipeline.paused )
            , ( "dashboard-running", isPipelineRunning state )
            , ( "dashboard-status-" ++ pipelineStatus state, not state.pipeline.paused )
            ]
        ]
        [ Html.div [ class "dashboard-pipeline-banner" ] []
        , Html.a [ class "dashboard-pipeline-content", href state.pipeline.url ]
            [ Html.div [ class "dashboard-pipeline-icon" ]
                []
            , Html.div [ class "dashboard-pipeline-name" ]
                [ Html.text state.pipeline.name ]
            ]
        , Html.div [] [ viewPipelinePreview state ]
        , Html.div [] (timeSincePipelineFailed now state)
        ]


timeSincePipelineFailed : Maybe Time -> PipelineState -> List (Html a)
timeSincePipelineFailed time { jobs } =
    case jobs of
        RemoteData.Success js ->
            let
                failedJobs =
                    List.filter ((==) Concourse.BuildStatusFailed << jobStatus) <| js

                failedDurations =
                    List.map
                        (\job ->
                            Maybe.withDefault { startedAt = Nothing, finishedAt = Nothing } <|
                                Maybe.map .duration job.transitionBuild
                        )
                        failedJobs

                failedDuration =
                    List.head <|
                        List.sortBy
                            (\duration ->
                                case duration.startedAt of
                                    Just date ->
                                        Time.inSeconds <| Date.toTime date

                                    Nothing ->
                                        0
                            )
                            failedDurations
            in
                case ( time, failedDuration ) of
                    ( Just now, Just duration ) ->
                        [ BuildDuration.viewFailDuration duration now ]

                    _ ->
                        []

        _ ->
            []


isPipelineRunning : PipelineState -> Bool
isPipelineRunning { jobs } =
    case jobs of
        RemoteData.Success js ->
            List.any (\job -> job.nextBuild /= Nothing) js

        _ ->
            False


pipelineStatus : PipelineState -> String
pipelineStatus { jobs } =
    case jobs of
        RemoteData.Success js ->
            Concourse.BuildStatus.show (jobsStatus js)

        _ ->
            "unknown"


jobStatus : Concourse.Job -> Concourse.BuildStatus
jobStatus job =
    Maybe.withDefault Concourse.BuildStatusPending <| Maybe.map .status job.finishedBuild


jobsStatus : List Concourse.Job -> Concourse.BuildStatus
jobsStatus jobs =
    let
        statuses =
            List.map (\job -> Maybe.withDefault Concourse.BuildStatusPending <| Maybe.map .status job.finishedBuild) jobs
    in
        if List.member Concourse.BuildStatusFailed statuses then
            Concourse.BuildStatusFailed
        else if List.member Concourse.BuildStatusErrored statuses then
            Concourse.BuildStatusErrored
        else if List.member Concourse.BuildStatusAborted statuses then
            Concourse.BuildStatusAborted
        else if List.member Concourse.BuildStatusSucceeded statuses then
            Concourse.BuildStatusSucceeded
        else
            Concourse.BuildStatusPending


fetchPipelines : Cmd Msg
fetchPipelines =
    Cmd.map PipelinesResponse <|
        RemoteData.asCmd Concourse.Pipeline.fetchPipelines


fetchJobs : Concourse.Pipeline -> Cmd Msg
fetchJobs pipeline =
    Cmd.map (JobsResponse pipeline.id) <|
        RemoteData.asCmd <|
            Concourse.Job.fetchJobs
                { teamName = pipeline.teamName
                , pipelineName = pipeline.name
                }


getCurrentTime : Cmd Msg
getCurrentTime =
    Task.perform ClockTick Time.now


initGraph : List Concourse.Job -> Graph Node ()
initGraph jobs =
    let
        jobNodes =
            List.map JobNode jobs

        jobsByName =
            List.foldl (\job dict -> Dict.insert job.name job dict) Dict.empty jobs

        resourceNodes =
            List.concatMap (jobResourceNodes jobsByName) jobs

        graphNodes =
            List.indexedMap Graph.Node (List.concat [ jobNodes, resourceNodes ])
    in
        Graph.fromNodesAndEdges
            graphNodes
            (List.concatMap (nodeEdges graphNodes) graphNodes)


jobResourceNodes : ByName Concourse.Job -> Concourse.Job -> List Node
jobResourceNodes jobs job =
    List.concatMap (inputNodes jobs job) job.inputs


inputNodes : ByName Concourse.Job -> Concourse.Job -> Concourse.JobInput -> List Node
inputNodes jobs job { resource, passed } =
    if List.isEmpty passed then
        []
    else
        List.map (constrainedInputNode jobs resource job) passed


constrainedInputNode : ByName Concourse.Job -> String -> Concourse.Job -> String -> Node
constrainedInputNode jobs resourceName dependentJob upstreamJobName =
    ConstrainedInputNode
        { resourceName = resourceName
        , dependentJob = dependentJob
        , upstreamJob = Dict.get upstreamJobName jobs
        }


nodeEdges : List (Graph.Node Node) -> Graph.Node Node -> List (Graph.Edge ())
nodeEdges allNodes { id, label } =
    case label of
        JobNode _ ->
            []

        ConstrainedInputNode { dependentJob, upstreamJob } ->
            Graph.Edge id (jobId allNodes dependentJob) ()
                :: case upstreamJob of
                    Just upstream ->
                        [ Graph.Edge (jobId allNodes upstream) id () ]

                    Nothing ->
                        []


viewSerial : Grid Node () -> List (Html msg)
viewSerial grid =
    case grid of
        Grid.Serial prev next ->
            viewSerial prev ++ viewSerial next

        _ ->
            viewGrid grid


viewGrid : Grid Node () -> List (Html msg)
viewGrid grid =
    case grid of
        Grid.Cell { node } ->
            viewNode node

        Grid.Serial prev next ->
            [ Html.div [ class "serial-grid" ]
                (viewSerial prev ++ viewSerial next)
            ]

        Grid.Parallel grids ->
            let
                foo =
                    List.concatMap viewGrid grids
            in
                if List.length foo <= 1 then
                    foo
                else
                    [ Html.div [ class "parallel-grid" ] foo ]

        Grid.End ->
            []


viewNode : Graph.Node Node -> List (Html msg)
viewNode { id, label } =
    case label of
        JobNode job ->
            let
                linkAttrs =
                    case ( job.finishedBuild, job.nextBuild ) of
                        ( Just fb, Just nb ) ->
                            Concourse.BuildStatus.show fb.status

                        ( Just fb, Nothing ) ->
                            Concourse.BuildStatus.show fb.status

                        ( Nothing, Just nb ) ->
                            "no-builds"

                        ( Nothing, Nothing ) ->
                            "no-builds"
            in
                [ Html.div [ class ("node job " ++ linkAttrs) ] [ Html.text job.name ] ]

        _ ->
            []


viewPipelinePreview : PipelineState -> Html msg
viewPipelinePreview { jobs } =
    case jobs of
        RemoteData.Success js ->
            let
                graph =
                    initGraph js
            in
                Html.div [ class "pipeline-grid" ] (viewGrid (Grid.fromGraph graph))

        RemoteData.Failure _ ->
            Html.div [] []

        _ ->
            Html.div [] []


jobId : List (Graph.Node Node) -> Concourse.Job -> Int
jobId nodes job =
    case List.filter ((==) (JobNode job) << .label) nodes of
        { id } :: _ ->
            id

        [] ->
            Debug.crash "impossible: job index not found"
