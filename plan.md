* people ssh into the game
    * if they are new, ask for a username
    * save username with ssh key
* players wait in the lobby
* settings menu to change name, color
    * possibly preferred direction?
* press join lobby to queue up
    * assigned a random direction
* game waits until there are four players
* game starts, snake spawns
* every time somebody presses their key
    * send to server
    * server simulates the game
* every timestep send state back to clients to display current position

* admin panel
    * request admin password
    * if correct, open admin panel
    * change names, kick players, ban players
