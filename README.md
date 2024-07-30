# Bracket Bot
#### Still work in progress | Made using discordgo / GoLang
#### A discord bot that utilize Challonge API to manage tournaments from just discord
#### Must have a Challonge account

````
TODO:
Set permission based command
Allow admin to set certain role to certain commands
Update tournament
Delete tournament
Create tournament
need more ideas
````

### SETUP:

Set up a config.toml file with this format \
You can get Challonge API key from here : https://challonge.com/settings/developer
````
[Bot]
token  = {discord bot token goes here}
guild_id = {discord server id goes here}

[API]
token = {challonge api key goes here}
````

After setting up config file, you can just run ``go run main.go``

### Commands:
````
/showalltournaments 
- Shows all tournaments
````
![img.png](img.png)
````
/showtournament {tourney-id} 
- Shows information of the tournament with given ID
````
![img_1.png](img_1.png)

````
/showparticipants {tourney-id} 
- Shows all participant in that tournament with given ID
````
![img_2.png](img_2.png)

````
/addparticipant {tourney-id} {name} optionals : {seed} {misc} {username} {email} 
- Add a participant with those parameters
````


````
/removeparticipant {tourney-id} {name} 
- Removes a participants from tournament
````
