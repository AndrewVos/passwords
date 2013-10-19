filledPassword = false

fill = ->
  ensureLoggedIn ->
    self.filledPassword = true
    url = "http://localhost:8080/search/?q=" + encodeURIComponent(site())
    $.getJSON url, (data) ->
      fillUsername(data.Username)
      fillPassword(data.Password)

ensureLoggedIn (callback) ->
  $.post "http://localhost:8080/logged_in", (data) ->
    if data.logged_in
      callback()
    else
      login()

login = ->
  iframe = '
    <iframe id="login-frame" srcdoc="" allowtransparency="true" scrolling="no" style="position: absolute; z-index:99999999; border: 0; width: 100%; height: 50px; top: 0;">
    </iframe>
  '
  popup = '
    <div id="login-flash" style="
      position: relative;
      top: 0;
      left: 0;
      width: 99%;
      height: 100%;
      border-radius: 5px;
      background-color: lightgray;
      padding: 5px;
      margin-left: auto;
      margin-right: auto;
      ">
        Please enter your password to decrypt the passwords file
    <input id="password" type="password" value="">
    <input id="login" type="submit">
    </div>
  '
  $("body").append($(iframe))
  frame = $("#login-frame")
  frame.attr("srcdoc", popup)

  frame.load ->
    frame.contents().find("#login").click ->
      url = "http://localhost:8080/login"
      $.post url, items.lastFormSubmitted, ->
        chrome.storage.local.remove "lastFormSubmitted"
        $("#login-frame").hide()

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
  return if self.filledPassword
  if $(e.target).find("input[type='password']").length >= 1
    hash = $(e.target).serializeHash()
    hash["site"] = site()
    chrome.storage.local.set({'lastFormSubmitted': hash})

chrome.storage.local.get "lastFormSubmitted", (items) ->
  if items.lastFormSubmitted && items.lastFormSubmitted.site == site()
    iframe = '
      <iframe id="passwords-frame" srcdoc="" allowtransparency="true" scrolling="no" style="position: absolute; z-index:99999999; border: 0; width: 100%; height: 50px; top: 0;">
      </iframe>
    '
    popup = '
      <div id="passwords-flash" style="
        position: relative;
        top: 0;
        left: 0;
        width: 99%;
        height: 100%;
        border-radius: 5px;
        background-color: lightgray;
        padding: 5px;
        margin-left: auto;
        margin-right: auto;
        ">
      Want to store this password?
      <input id="passwords-yes" type="submit" value="yeah">
      <input id="passwords-no" type="submit" value="nope">
      </div>
    '
    $("body").append($(iframe))
    frame = $("#passwords-frame")
    frame.attr("srcdoc", popup)

    clickYes = ->
      url = "http://localhost:8080/store"
      $.post url, items.lastFormSubmitted, ->
        chrome.storage.local.remove "lastFormSubmitted"
        $("#passwords-frame").hide()
    clickNo = ->
      chrome.storage.local.remove "lastFormSubmitted"
      $("#passwords-frame").hide()

    frame.load ->
      frame.contents().find("#passwords-yes").click clickYes
      frame.contents().find("#passwords-no").click clickNo
