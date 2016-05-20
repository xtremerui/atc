port module BuildPage exposing (..)

import Html.App
import Time

import Autoscroll
import Build
import Navigation
import Scroll


type alias Flags =
  { navState : Navigation.CurrentState
  , buildFlags : Build.Flags
  }

main : Program Flags
main =
  Html.App.programWithFlags
    { init = init
    , update = Navigation.update <| Autoscroll.update Build.update
    , view = Navigation.view <| Autoscroll.view Build.view
    , subscriptions =
        let
          scrolledUp =
            Scroll.fromBottom Autoscroll.FromBottom

          pushDown =
            Time.every (100 * Time.millisecond) (always Autoscroll.ScrollDown)

          tick =
            Time.every Time.second Build.ClockTick
        in \model ->
          Sub.map Navigation.SubAction <|
            Sub.batch
              [ scrolledUp
              , pushDown
              , Sub.map Autoscroll.SubAction <|
                  Sub.batch
                    [ tick
                    , case model.subModel.subModel.currentBuild `Maybe.andThen` .output of
                        Nothing ->
                          Sub.none
                        Just buildOutput ->
                          Sub.map Build.BuildOutputAction buildOutput.events
                    ]
              ]
    }

init : Flags -> (Navigation.Model (Autoscroll.Model Build.Model), Cmd (Navigation.Action (Autoscroll.Action Build.Action)))
init flags =
  Navigation.init flags.navState <|
    Autoscroll.init
    Build.getScrollBehavior <|
      Build.init flags.buildFlags
