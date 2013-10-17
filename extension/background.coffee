fill = ->
  url = "http://localhost:8080/search/?q=" + encodeURIComponent(site())
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

site = ->
  if !window.location.origin
    window.location.origin = window.location.protocol+"//"+window.location.host
  window.location.origin

Mousetrap.bindGlobal 'ctrl+\\', (e) ->
  fill()
  false

Mousetrap.bindGlobal 'cmd+\\', (e) ->
  fill()
  false

$("form").submit (e) ->
  if $(e.target).find("input[type='password']").length >= 1
    hash = $(e.target).serializeHash()
    hash["site"] = site()
    url = "http://localhost:8080/store"
    $.post url, hash
