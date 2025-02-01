# Discord Bot Setup

To setup a discord you first need to create an application (basically a bot), here are the steps to follow:

1. Visit [Discord Application](https://discord.com/developers/applications 'Discord Application')
2. Cick on "New Application" Button
   ![Screenshot 2024-11-08 at 1 52 53 AM](https://github.com/user-attachments/assets/380657ca-89b4-4053-96c9-6b73632d382c)
3. Fill in your Application Name
   ![Screenshot 2024-11-08 at 1 52 53 AM](https://github.com/user-attachments/assets/688bd69d-fcca-4a80-8780-9ab18bfc5037)
4. Now [here](https://discord.com/developers/applications) you will see your newly created application, click on it, this should open _"General Information"_
5. If you scroll down, you will se _PUBLIC KEY_, copy it and place it `.env` as _DISCORD_PUBLIC_KEY_
6. Now to create BOT_TOKEN, click on BOT > Reset Token
   ![Screenshot 2024-11-08 at 2 07 40 AM](https://github.com/user-attachments/assets/201f9e51-a44a-43af-9c96-4eaf453d02b0)
7. Once you have the token, place it against _BOT_TOKEN_ in `.env`
8. Now will be creating an invite URL and for that you need to click on OAuth2 > bot
   ![Screenshot 2024-11-13 at 11 40 21 PM](https://github.com/user-attachments/assets/aebad7fe-aa82-45de-bb17-25dc0fff0e5f)
9. Now as soon as you click on bot, a section to choose bot permission from, will shown up
   ![Screenshot 2024-11-14 at 10 53 20 AM](https://github.com/user-attachments/assets/b6fc4afb-4de4-449c-bf39-f8a0b4d3de06)
10. Check the following options

    1. [ ] Send Messages

11. Once you select all the bot permissions, scroll a bit down and you will see "Generated URL"
    ![Screenshot 2024-11-14 at 10 58 30 AM](https://github.com/user-attachments/assets/bbff4c6d-4ef5-46fd-89c7-9acf31c11cdd)
12. Copy and paste that URL in browser, a prompt will come up where it will ask you to select you own "Discord Server"
    ![Screenshot 2024-11-14 at 11 00 45 AM](https://github.com/user-attachments/assets/322caf6d-af84-4752-88db-0ce64e080d6d)
13. Once you add the Bot into your server, copy the "Server Id", by right clicking on the server avatar. Now place this id in `.env` against _GUILD_ID_

# Connecting Discord Service with Discord

Now as you have created the discord bot, now its time to connect it with discord service using the following steps:

1. You would need to register the commands first. That will be auto handled once you start the server
2. Now start the server using

```bash
   make run #or go run .
```

3. For IP tunneling, need to run NGROK, use the following command

```bash
   make ngrok #or ngrok http 8999
```

Since we are considering 8999 as default port for this service. If you wish to change it you can change it in `Makefile` & in `docker-compose.yml`

4. Copy the Ngrok URL and open the General Information on [Discord Developer Portal](https://discord.com/developers/applications) of your bot, paste the copied URL in Interactions Endpoint URL
   ![Screenshot 2024-11-14 at 10 58 30 AM](https://github.com/user-attachments/assets/53f372e4-44e7-4cdc-acfc-0e3b707f8607)
5. All Set 🚀🚀🚀. Now you can start with running hello command
