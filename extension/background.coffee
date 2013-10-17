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
    chrome.storage.local.set({'lastFormSubmitted': hash})

chrome.storage.local.get "lastFormSubmitted", (items) ->
  if items.lastFormSubmitted && items.lastFormSubmitted.site == site()
    popup = '
      <div id="passwords-flash" style="position: absolute; z-index:99999999; width: 100%; height: 20px; top: 0; background-color: lightgray; border: 1px solid black;">
      Want to store this password? <input id="passwords-yes" type="submit" value="yeah"><input id="passwords-no" type="submit" value="nope">
      </div>
    '
    $("body").append($(popup))

    $(document).on "click", "#passwords-yes", ->
      url = "http://localhost:8080/store"
      $.post url, items.lastFormSubmitted, ->
        chrome.storage.local.remove "lastFormSubmitted"
        $("#passwords-flash").hide()
    $(document).on "click", "#passwords-no", ->
      chrome.storage.local.remove "lastFormSubmitted"
      $("#passwords-flash").hide()
