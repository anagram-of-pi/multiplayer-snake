* people ssh into the game
    * if they are new, ask for a username
    * save username with ssh key
* players wait in the lobby
    * open settings menu to change name, color
    * possibly preferred direction?
    * possibly other stuff?
* press join to queue up
    * assigned a random open direction
* game waits until there are four players
    * After four players are in, waits for ready
* game starts, snake spawns, etc
* every time somebody presses their key
    * send to server
    * server simulates the game
* every timestep send state back to clients to display current position


...
* admin panel
    * request admin password
    * if correct, open admin panel
    * change names, kick players, ban players



---

ssh shows the tui via middleware
* How does the tui connect to and talk with the server?
    * ClientMessage and ServerMessage?
    * Channel for them?
* How does the server start?
    * In main, start a goroutine for the multiplayer game manager?


## Main
* Starts Game Manager goroutine
* Creates the ssh server + middleware
* Serves the server

## TUI
* Gets displayed when sshing
* Requests from server the username
    * If none, displays a page with a textbox
    * Sends back to server
* Main page
    * Logo
    * Join button
    * Settings
* Periodically polls for number of players waiting
    * Shows under join button
    * Use channel with ClientBoundPacket?
* When join button pressed
    * Sends message to server
    * Asks server for information
    * Polls for players
    * Or wait on a gameStart channel?
* Game page
    * Grid using bubble tea (how? lipstick?)
    * Create snake board client side
    * Display who is in the game
    * Display score
    * Display controls/current direction
    * Leave button


## Game Manager
* Manages all connections and converts between TUI and snake game
* Every time a player 















---
