[![LinkedIn][linkedin-shield]][linkedin-url]


<br />
<p align="center">
  <h1 align="center">Topic Sync</h1>
  <h2 align="center">Low-latency IRC<->Discord Channel Topic Syncing</h2>
  <h3 align="center">Written in Go</h3>
</p>


## Table of Contents
* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
  * [Docker](#Docker)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)


## About The Project

[![Short Topic Sync Demo][product-screenshot]]

This is a simple, small, and easy to configure golang bot for connecting and bidirectionally syncing a __single__ IRC channel's topic to the topic of a __single__ Discord channel. 

Any topic updates in those channels will be synced with the other nearly instantly and a topic change alert will be posted in Discord to mirror the alert normally shown in IRC.

It was made to compliment an existing kubernetes cluster of chat and site services for an independent game streaming community spread out over two chat networks. Since [matterbridge](https://github.com/42wim/matterbridge) doesn't support topic syncing for the time being, this was developed to fill in the feature gap.

If there is demand, this project could be extended to support multiple channels, networks, or other configuration needs so feel free to <a href="https://github.com/CoordSpace/topicsync/issues">request a feature!</a>


### Built With

* [discordgo](https://github.com/bwmarrin/discordgo)
* [girc](https://github.com/lrstanley/girc)
* [viper](https://github.com/spf13/viper)


## Getting Started

To get a local copy compiled and running follow these simple steps.


### Prerequisites

* A working Go install with a functional GOPATH
* A Q/Nickserv account registered for the bot with auto-op/half-op set for its account in the desired IRC channel
* Lastly, a bot account must be created in Discord and added to your Discord server with the Manage Channels permission (to be able to edit the topic)


### Installation
 
Download the latest code into your go directory:
```sh
cd $GOPATH
go get github.com/CoordSpace/topicsync
```

Once that's done, you should have a compiled binary in your /bin directory!

Now open up a new yaml file called `config.yaml` and paste the contents of the `config-sample.yaml` file from this repo into it and edit to meet your needs.


#### Example Config.yaml

```
---
irc:
  server: irc.coolnetwork.chat
  port: 6667
  user: TopicBot
  nick: TopicBot
  name: A Topic Sync Bot
  # IRC channel to join (e.g. '#TopicBotTest')
  channel: '#CoolChatroom
  # Optional - Bot name and command to auth with IRC services
  auth:
    bot: NickServ@coolnetwork.chat
    cmd: IDENTIFY foo password
discord:
  token: 'a-big-string-from-the-discord-dev-portal-bot-settings'
  # The ID of the channel to track, found using discord dev-mode 
  # and right-clicking on the channel
  channelID: '38104041843693923425'
# Topic update message formatting string for Discord
# [Emoji] Topic Updated: [This is the topic message.]
updateFormat: "%s Topic Updated: %s"
# Optional - Random Emojis for added topic update flair in Discord
emojis:
  - ⚠️
  - 🔔
  - 🚨
...

```

This configuration file must be stored in one of three valid directories in order for this program to use it:

* `/etc/topicsync/`
* `$HOME/.topicsync`
* `./` (The same directory as the binary)

Lastly just run the binary using the process supervisor of your choice. It will automatically find your config file and start connecting to the servers and channels of your choice. By default, __the existing channel topics in the pair will not sync immediately upon joining.__ So issue a new topic on either side to get the process started. 


### Docker
A hyper-minimal, statically-linked Docker image is also available for fast integration into existing containerized stacks.

Just create your local `config.yaml` configuration from the repo's sample template and mount it to the latest Docker build with:

`docker run -ti -v /path/to/config.yaml:/etc/topicsync/config.yaml coordspace/topicsync`

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/CoordSpace/topicsync/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Chris Earley - [@CoordSpace](https://twitter.com/CoordSpace) - chris@coord.space

Project Link: [https://github.com/CoordSpace/topicsync](https://github.com/CoordSpace/topicsync)


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=flat-square&logo=linkedin&colorB=555
[linkedin-url]: https://coord.space/in
[product-screenshot]: https://giant.gfycat.com/ExcellentIdleBluetickcoonhound.gif
