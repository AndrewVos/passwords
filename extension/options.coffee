$ ->
  $.getJSON "http://localhost:8080/passwords_file_exists/", (data) ->
    if data.passwords_file_exists
      $("#create-passwords-file").hide()
    else
      $("#create-passwords-file").show()

  toggleSubmitButton = ->
    button = $("#save")
    password = $("#password")
    confirm_password = $("#confirm_password")
    if password.val() == "" || (password.val() != confirm_password.val())
      button.attr("disabled", "true")
    else
      button.removeAttr("disabled")

  toggleSubmitButton()
  $("#password, #confirm_password").keyup toggleSubmitButton
