fillIn = () ->
  console.log window.location

Mousetrap.bind 'ctrl+\\', ->
  fillIn()
  return false

Mousetrap.bind 'cmd+\\', ->
  fillIn()
  return false
