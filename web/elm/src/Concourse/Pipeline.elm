module Concourse.Pipeline where

import Http
import Json.Decode exposing ((:=))
import Task exposing (Task)

type alias Pipeline =
  { name : String
  , paused : Bool
  }

fetch : String -> Task Http.Error Pipeline
fetch pipelineName =
  Http.get decode ("/api/v1/pipelines/" ++ pipelineName)

url : String -> String
url pipelineName =
  "/pipelines/" ++ pipelineName

urlAll : String
urlAll =
  "/pipelines"

urlJobs : String -> String
urlJobs pipelineName =
  "/pipelines/" ++ pipelineName ++ "/jobs"

fetchAll : Task Http.Error (List Pipeline)
fetchAll =
  Http.get (Json.Decode.list decode) "/api/v1/pipelines"

decode : Json.Decode.Decoder Pipeline
decode =
  Json.Decode.object2 Pipeline
    ("name" := Json.Decode.string)
    ("paused" := Json.Decode.bool)
