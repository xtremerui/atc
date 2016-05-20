module Main where

import Html exposing (Html)
import Effects
import StartApp
import Task exposing (Task)
import Time

import Autoscroll
import Build
import Navigation
import Scroll

port buildId : Int

port navState : Navigation.CurrentState

main : Signal Html
main =
  app.html

app : StartApp.App (Navigation.Model (Autoscroll.Model Build.Model))
app =
  let
    pageDrivenActions =
      Signal.mailbox Build.Noop
  in
    StartApp.start
      { init =
          Navigation.init navState <|
            Autoscroll.init
              Build.getScrollBehavior <|
                Build.init redirects.address pageDrivenActions.address buildId
      , update = Navigation.update (Autoscroll.update Build.update)
      , view = Navigation.view (Autoscroll.view Build.view)
      , inputs =
          [ Signal.map (Navigation.SubAction << Autoscroll.SubAction) pageDrivenActions.signal
          , Signal.merge
              (Signal.map (Navigation.SubAction << Autoscroll.FromBottom) Scroll.fromBottom)
              (Signal.map (Navigation.SubAction << always Autoscroll.ScrollDown) (Time.every (100 * Time.millisecond)))
          ]
      , inits = [Signal.map (Navigation.SubAction << Autoscroll.SubAction << Build.ClockTick) (Time.every Time.second)]
      }

redirects : Signal.Mailbox String
redirects = Signal.mailbox ""

port redirect : Signal String
port redirect =
  redirects.signal

port tasks : Signal (Task Effects.Never ())
port tasks =
  app.tasks
