module Concourse.Pipeline exposing (..)

import Http
import Json.Decode exposing ((:=))
import Task exposing (Task)

type alias Pipeline =
  { name : String
  , teamName : String
  , paused : Bool
  }

fetch : String -> String -> Task Http.Error Pipeline
fetch teamName pipelineName =
  Http.get decode ("/api/v1/teams/" ++ teamName ++ "/pipelines/" ++ pipelineName)

url : { name : String, teamName : String } -> String
url {name, teamName} =
  "/teams/" ++ teamName ++ "/pipelines/" ++ name

urlAll : String -> String
urlAll teamName =
  "/teams/" ++ teamName ++ "/pipelines"

urlJobs : { name : String, teamName : String } -> String
urlJobs {name, teamName} =
  "/teams/" ++ teamName ++ "/pipelines/" ++ name ++ "/jobs"

fetchAll : String -> Task Http.Error (List Pipeline)
fetchAll teamName =
  Http.get (Json.Decode.list decode) ("/api/v1/teams/" ++ teamName ++ "/pipelines")

decode : Json.Decode.Decoder Pipeline
decode =
  Json.Decode.object3 Pipeline
    ("name" := Json.Decode.string)
    ("team_name" := Json.Decode.string)
    ("paused" := Json.Decode.bool)
