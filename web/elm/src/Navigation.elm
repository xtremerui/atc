module Navigation exposing (..)

import Html exposing (Html)
import Html.App
import Html.Attributes exposing (action, class, classList, href, id, method, title, disabled, attribute)
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
    Html.nav [class "nav-sidebar"] [
      Html.form [class "magic-search"] [
        Html.input [Html.Attributes.type' "text", Html.Attributes.placeholder "filter…"] []
      ],

      case model.currentState.job of
        Nothing ->
          Html.text "i aint got NOOOO nav for you, boy"

        Just job ->
          Html.div [class "events"] [
            -- Html.div [class "event pipeline"] [
            --   Html.a [href Concourse.Pipeline.urlAll] [
            --     Html.text "pipelines"
            --   ]
            -- ],
            -- Html.div [class "event job"] [
            --   Html.a [href (Concourse.Pipeline.url job.pipelineName)] [
            --     Html.text job.pipelineName
            --   ]
            -- ],
            Html.div [class "event build started"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag started"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "5 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event resource"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-cube"] [],

                Html.a [href (Concourse.Pipeline.urlJobs job.pipelineName)] [
                  Html.text "concourse"
                ]
              ],
              Html.div [class "event-body resource-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "ref"],
                    Html.td [class "dict-value"] [Html.text "0fd4d2a444584fdb949fc99548a7bda604c9e368"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build failed"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag failed"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "failed"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event resource"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-cube"] [],

                Html.a [href (Concourse.Pipeline.urlJobs job.pipelineName)] [
                  Html.text "concourse"
                ]
              ],
              Html.div [class "event-body resource-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "version"],
                    Html.td [class "dict-value"] [Html.text "1.2.0"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event resource"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-cube"] [],

                Html.a [href (Concourse.Pipeline.urlJobs job.pipelineName)] [
                  Html.text "concourse"
                ]
              ],
              Html.div [class "event-body resource-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "ref"],
                    Html.td [class "dict-value"] [Html.text "0fd4d2a444584fdb949fc99548a7bda604c9e368"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event resource"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-cube"] [],

                Html.a [href (Concourse.Pipeline.urlJobs job.pipelineName)] [
                  Html.text "concourse"
                ]
              ],
              Html.div [class "event-body resource-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "ref"],
                    Html.td [class "dict-value"] [Html.text "0fd4d2a444584fdb949fc99548a7bda604c9e368"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event resource"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-cube"] [],

                Html.a [href (Concourse.Pipeline.urlJobs job.pipelineName)] [
                  Html.text "concourse"
                ]
              ],
              Html.div [class "event-body resource-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "ref"],
                    Html.td [class "dict-value"] [Html.text "0fd4d2a444584fdb949fc99548a7bda604c9e368"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build succeeded"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag succeeded"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "succeeded"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ],
            Html.div [class "event build aborted"] [
              Html.h3 [class "event-header"] [
                Html.i [class "fa fa-hashtag aborted"] [],

                Html.a [href (Concourse.Job.url job)] [
                  Html.text (job.name ++ " #23")
                ]
              ],
              Html.div [class "event-body build-info"] [
                Html.table [class "dictionary"] [
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "started"],
                    Html.td [class "dict-value"] [Html.text "5 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "aborted"],
                    Html.td [class "dict-value"] [Html.text "9 minutes ago"]
                  ],
                  Html.tr [] [
                    Html.td [class "dict-key"] [Html.text "duration"],
                    Html.td [class "dict-value"] [Html.text "4 minutes"]
                  ]
                ]
              ]
            ]
          ]
    ],

    Html.div [class "nav-content"] [
      Html.App.map SubAction (subView model.subModel)
    ]
  ]

fetchPipelines : Cmd (Action subAction)
fetchPipelines =
  Cmd.map PipelinesFetched << Task.perform Err Ok <|
    Concourse.Pipeline.fetchAll
