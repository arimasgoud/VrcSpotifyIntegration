# VrcSpotifyIntegration
Abusing the Spotify API for cross-platform spotify control in VRChat

## THIS REQUIRES SPOTIFY PREMIUM TO FUNCTION

## Quickstart

On first launch, the mod will prompt you to follow this guide to get set up.<br>
Setup is quite simple: go to [Spotify Developers](https://developer.spotify.com/dashboard/login), log in, and create a new Application.<br>
On the top right, click on `Edit Settings`, and add a callback URL like this.
![image](https://user-images.githubusercontent.com/29461788/156894291-fd429bca-6e20-4972-a8c8-04e0b370fcd2.png)<br>
The port number is up to you, just remember it.

Write down the Client ID, Client Secret, and the port you chose<br>
![image](https://user-images.githubusercontent.com/29461788/156894335-286ff528-b5ad-40dc-bd06-d6e499dfa2c4.png)

In the game, it will ask you for all three of these values. enter them, and your Console should tell you that you've logged in successfully.

## Notes & Credits

* This mod uses ReMod.Core by [RequiDev](https://github.com/RequiDev)
* Note that should you not have ReMod.Core, the mod will inject it automatically. Please consider [downloading it yourself](https://github.com/RequiDev/ReMod.Core/releases/latest)
* The Spotify API wrapper by [zmb3](https://github.com/zmb3/spotify)

