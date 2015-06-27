package static

//go:generate elm-make Main.elm
//go:generate uglifyjs -m -c --lint=false --stats -o elm.min.js -- elm.js
