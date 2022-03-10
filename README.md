# VrcSpotifyIntegration
Abusing the Spotify API for cross-platform spotify control in VRChat

## THIS REQUIRES SPOTIFY PREMIUM TO FUNCTION

## Quickstart

On first launch, the mod will prompt you to follow this guide to get set up.<br>
Setup is quite simple: go to [Spotify Developers](https://developer.spotify.com/dashboard/login), log in, and create a new Application.<br>
On the top right, click on `Edit Settings`, and add a callback URL like this.

### Important:
the url **Has** to be `http`

![image](https://user-images.githubusercontent.com/29461788/157617637-004d4240-0952-40be-b0b7-8f9bf415e51e.png)
The port number is up to you, just remember it.

Click `Add`, and click `Save` at the bottom

Write down the Client ID, Client Secret, and the port you chose<br>
![image](https://user-images.githubusercontent.com/29461788/156894335-286ff528-b5ad-40dc-bd06-d6e499dfa2c4.png)

In the game, it will ask you for all three of these values. enter them, and your Console should tell you that you've logged in successfully.<br>
Remember that the mod **ONLY** wants the number you put for the port, **NOT** the entire callback uri!

## Notes & Credits

* This mod uses ReMod.Core by [RequiDev](https://github.com/RequiDev)
* Note that should you not have ReMod.Core, the mod will inject it automatically. Please consider [downloading it yourself](https://github.com/RequiDev/ReMod.Core/releases/latest)
* The Spotify API wrapper by [zmb3](https://github.com/zmb3/spotify)

## Building it Yourself
* Clone the Repo
* Navigate to the `GotifyNative` folder
* `go build -trimpath --buildmode=c-shared -ldflags="-s -w"  -o GotifyNative.dll gotify.go MelonLogger.go mono.go`
* Place it in the `VrcSpotifyMod` directory
* build the C# project

