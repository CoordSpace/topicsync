[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <h1 align="center">Topic Sync</h1>
  <h2 align="center">Low-latency IRC<->Discord Channel Topic Syncing</h2>
  <h3 align="center">Written in Go</h3>
</p>


<!-- TABLE OF CONTENTS -->
## Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)


<!-- ABOUT THE PROJECT -->
## About The Project

[![Short Topic Sync Demo][product-screenshot]]

This is a simple, small, and easy to configure golang bot for connecting and bidirectionally syncing a __single__ IRC channel's topic to the topic of a __single__ Discord channel. 

Any topic updates in those channels will be synced with the other nearly instantly and a topic change alert will be posted in Discord to mirror the alert normally shown in IRC.

It was made to compliment a large stack of chat and site services for an independent game streaming community spread out over two chat networks. Since [matterbridge](https://github.com/42wim/matterbridge) doesn't support topic syncing for the time being, this was developed to fill in the gap.

If there's demand, this project could be extended to support multiple channels or other configuration needs so feel free to <a href="https://github.com/CoordSpace/topicsync/issues">request a feature!</a>


### Built With

* [discordgo](https://github.com/bwmarrin/discordgo)
* [girc](https://github.com/lrstanley/girc)
* [viper](https://github.com/spf13/viper)



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* npm
```sh
npm install npm@latest -g
```

### Installation
 
1. Clone the topicsync
```sh
git clone https:://github.com/CoordSpace/topicsync.git
```
2. Install NPM packages
```sh
npm install
```

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
