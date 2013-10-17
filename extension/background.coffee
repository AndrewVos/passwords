fill = ->
  url = "http://localhost:8080/?q=" + encodeURIComponent(host())
  $.getJSON url, (data) ->
    fillUsername(data.Username)
    fillPassword(data.Password)

fillUsername = (username) ->
  username_fields = [
    "login_email",
    "email"
  ]
  tryFillIn(username_fields, username)

fillPassword = (password) ->
  password_fields = [
    "login_password",
    "password"
  ]
  tryFillIn(password_fields, password)

tryFillIn = (fields, value) ->
  for field in fields
    do (field) ->
      $("input").each (index, element) ->
        element = $(element)
        if element.attr("name") == field
          element.attr "value", value

host = ->
  if !window.location.origin
    window.location.origin = window.location.protocol+"//"+window.location.host
  window.location.origin

Mousetrap.bindGlobal 'ctrl+\\', (e) ->
  fill()
  false

Mousetrap.bindGlobal 'cmd+\\', (e) ->
  fill()
  false
