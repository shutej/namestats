import Array
import History
import Html as H exposing (..)
import Html.Attributes as HA exposing (..)
import Html.Events exposing (on, onClick, targetValue)
import Http
import Json.Decode as Json exposing ((:=))
import Maybe exposing (..)
import String
import Svg as S exposing (..)
import Svg.Attributes as SA exposing (..)
import Task exposing (..)

-- These are stylistic options that control rendering.
-- TODO(shutej): It would be better to get these in a different way, especially
-- if we want the UI to be responsive.
options : GraphOptions
options =
    { width      = 1178
    , height     = 100
    , keyMin     = 1880
    , keyMax     = 2014
    , yMax       = 100000
    , yTransform = sqrt
    }

type Gender = Male | Female

stringToGender : String -> Gender
stringToGender string =
    case string of
      "m" -> Male
      "f" -> Female

decodeGender : Json.Decoder Gender
decodeGender = Json.map stringToGender Json.string

type alias Series a =
    { key  : Int
    , data : Array.Array a
    }

decodeMaybe : Json.Decoder a -> Json.Decoder (Maybe a)
decodeMaybe decoder =
    Json.oneOf [ Json.null Nothing
               , Json.map Just decoder ]

decodeSeries : Json.Decoder a -> Json.Decoder (Series a)
decodeSeries decoder =
    Json.object2 (Series)
        ("key"  := Json.int)
        ("data" := Json.array decoder)

type alias Name =
    { gender       : Gender
    , name         : String
    , rank         : Series (Maybe Int)
    , count        : Series (Maybe Int)
    , relatedNames : List String
    , totalCount   : Int
    }

decodeName : Json.Decoder Name
decodeName =
    Json.object6 (Name)
        ("gender"       := decodeGender)
        ("name"         := Json.string)
        ("rank"         := decodeSeries (decodeMaybe Json.int))
        ("count"        := decodeSeries (decodeMaybe Json.int))
        ("relatedNames" := Json.list Json.string)
        ("totalCount"   := Json.int)

type alias Response =
    { offset     : Int
    , limit      : Int
    , total      : Int
    , nameSearch : List Name
    }

decodeResponse : Json.Decoder Response
decodeResponse =
    Json.object4 (Response)
        ("offset"     := Json.int)
        ("limit"      := Json.int)
        ("total"      := Json.int)
        ("namesearch" := Json.list decodeName)

type alias Error = Maybe Http.Error

getNameSearch : String -> Task Error Response
getNameSearch query =
    Http.get decodeResponse ("/v1/namesearch/" ++ (Http.uriEncode query))
        |> mapError (Just)

type alias QueryChange = Maybe String

queryChanges : Signal.Mailbox QueryChange
queryChanges = Signal.mailbox Nothing

queryChangeToPath : QueryChange -> String
queryChangeToPath change =
    case change of
      Just query -> "/#" ++ (Http.uriEncode query)
      _ -> ""

port replaceHash : Signal (Task x ())
port replaceHash =
    Signal.map queryChangeToPath queryChanges.signal
          |> Signal.map History.replacePath

-- This maps the hash values to results displayed on the page.
-- TODO(j): Entering page and fragment sometimes redirects you to "/#" it seems?
results : Signal.Mailbox (Result Error Response)
results = Signal.mailbox (Err Nothing)

hashToQuery : String -> String
hashToQuery hash =
    case String.uncons hash of
      Just ('#', query) -> Http.uriDecode query
      _ -> ""

queries : Signal String
queries = Signal.map hashToQuery History.hash

port requests : Signal (Task x ())
port requests =
    queries
      |> Signal.map getNameSearch
      |> Signal.map (\task -> Task.toResult task `Task.andThen` Signal.send results.address)

viewNameSearchRangeNotEmpty : Response -> Html
viewNameSearchRangeNotEmpty response =
    let
        rangeStart = toString (response.offset + 1)
        rangeEnd   = toString (response.offset + List.length response.nameSearch)
        rangeTotal = toString response.total
    in
      div [ HA.class "namesearch" ]
              [ span [ HA.class "start" ] [ H.text rangeStart ]
              , span [ HA.class "end" ] [ H.text rangeEnd ]
              , span [ HA.class "total" ] [ H.text rangeTotal ] ]

viewNameSearchRangeEmpty : Response -> Html
viewNameSearchRangeEmpty response =
    div [ HA.class "namesearch" ] [ span [ HA.class "none" ] [] ]

viewNameSearchRange : Response -> Html
viewNameSearchRange response =
    if | response.total > 0 -> viewNameSearchRangeNotEmpty response
       | otherwise          -> viewNameSearchRangeEmpty response

viewResultName : String -> Html
viewResultName name =
    H.a [ HA.class "name"
        , onClick queryChanges.address (Just name) ] [ H.text name ]

viewRelatedName : String -> Html
viewRelatedName relatedName =
    li [ HA.class "related-name" ] [ viewResultName relatedName ]

viewRelatedNames : List String -> Html
viewRelatedNames relatedNames =
    List.map viewRelatedName relatedNames
        |>  ul [ HA.class "related-names" ]

viewGender : Gender -> Html -> Html
viewGender gender html =
    case gender of
      Male -> h1 [ HA.class "gender male" ] [ html ]
      Female -> h1 [ HA.class "gender female" ] [ html ]

viewTotalCount : Int -> Html
viewTotalCount totalCount =
    div [ HA.class "total-count" ] [ H.text <| toString totalCount ]

viewNameSearch : Name -> Html
viewNameSearch name =
    let
        list = seriesToRankCount name.rank name.count
    in
      div [ HA.class "result" ]
              [ viewGender name.gender (viewResultName name.name)
              , viewBarGraph list
              , viewTotalCount name.totalCount
              , viewRelatedNames name.relatedNames
              ]

viewResult : Result Error Response -> Html
viewResult result =
    case result of
      Err Nothing ->
          div [ HA.class "namesearch" ] [ H.text "Waiting for results..." ]
      Err (Just error) ->
          div [ HA.class "namesearch" ] [ H.text "Error retrieving results!" ]
      Ok response ->
          (viewNameSearchRange response :: List.map viewNameSearch response.nameSearch)
              |> div []

type alias GraphOptions =
    { width      : Float
    , height     : Float
    , keyMin     : Int
    , keyMax     : Int
    , yMax       : Float
    , yTransform : Float -> Float
    }

type alias RankCount a =
    { index : a
    , key   : a
    , rank  : a
    , count : a
    }

convertRankCount : (a -> b) -> RankCount a -> RankCount b
convertRankCount convert rankCount =
    { index = convert rankCount.index
    , key   = convert rankCount.key
    , rank  = convert rankCount.rank
    , count = convert rankCount.count
    }


dropMaybe : Maybe a -> List a -> List a
dropMaybe maybe list =
    case maybe of
      Just x -> x :: list
      Nothing -> list

seriesToRankCount : Series (Maybe Int) -> Series (Maybe Int) -> List (RankCount Int)
seriesToRankCount rankSeries countSeries =
    let
        list = [options.keyMin .. options.keyMax]
        helper k =
            let
                rank = Array.get (k - rankSeries.key) rankSeries.data
                count = Array.get (k - countSeries.key) countSeries.data
                i = k - options.keyMin
            in
              case (rank, count) of
                (Just (Just r), Just (Just c)) ->
                    Just { index = i, key = k, rank = r, count = c }
                _ -> Nothing
    in
      List.map helper list |> List.foldr dropMaybe []

viewBarGraphBar : Float -> Float -> (RankCount Int) -> Html
viewBarGraphBar unitWidth unitHeight rankCount =
    let
        rc     = convertRankCount toFloat rankCount
        width  = unitWidth - 1
        height = unitHeight * (rc.count |> options.yTransform)
        x      = rc.index * unitWidth
        y      = unitHeight * ((options.yMax |> options.yTransform) -
                               (rc.count |> options.yTransform))
        caption = "There were " ++ (toString rankCount.count)
                  ++ " births in " ++ (toString rankCount.key)
                  ++ ". (#" ++ (toString rankCount.rank) ++ ")"
    in
      S.g [] [ S.rect [ SA.class "bar"
                      , SA.width (width |> toString)
                      , SA.height (height |> toString)
                      , SA.x (x |> toString)
                      , SA.y (y |> toString)
                      ]
               []
             , S.rect [ SA.class "caption"
                      , SA.width (width |> toString)
                      , SA.height (options.height |> toString)
                      , SA.x (x |> toString)
                      , SA.y "0"
                      ]
               [ S.title [] [ H.text caption ] ] ]

viewBarGraph : List (RankCount Int) -> Html
viewBarGraph list =
    let
        unitWidth  = options.width / (options.keyMax - options.keyMin |> toFloat)
        unitHeight = options.height / (options.yMax |> options.yTransform)
    in
      S.svg [ version "1.1"
            , SA.width (options.width |> toString)
            , SA.height (options.height |> toString)
            , SA.class "bar-graph"]
           (List.map (viewBarGraphBar unitWidth unitHeight) list)

view : String -> Result Error Response -> Html
view value' result =
    let
        input' = input [ on "input"
                                (decodeMaybe targetValue)
                                (Signal.message queryChanges.address)
                       , placeholder "Type a name here..."
                       , value value' ] []
        output = viewResult result
    in
      div [] [ input', output ]

main =
    Signal.map2 view queries results.signal
