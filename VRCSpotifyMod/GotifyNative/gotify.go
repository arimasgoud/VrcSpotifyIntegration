package main

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
import "C"
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

//export HasRemodCore
func HasRemodCore() bool {
	if _, err := os.Stat("ReMod.Core.dll"); err == nil {
		//ReMod.Core exists
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		//ReMod.Core does not exist
		return false
	}
	//Schrödinger's ReMod.Core
	return false
}

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
var redirectURI = "http://localhost:6942/callback"

var (
	auth          *spotifyauth.Authenticator
	ch            = make(chan *spotify.Client)
	SpotifyClient *spotify.Client
	state         = "abc123"
	Instance      *LoggerInstance
)

//export Login
func Login(id, token, port uintptr) uintptr {
	// first start an HTTP server
	spotify_id := getString(id)
	spotify_secret := getString(token)
	spotify_port := getString(port)
	redirectURI = "http://localhost:" + spotify_port + "/callback"
	srv := &http.Server{Addr: fmt.Sprintf(":%s", spotify_port)}
	var err error

	mono.ThreadAttach()

	Instance, err = NewLoggerInstance("GotifyNative")
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	os.Setenv("SPOTIFY_ID", spotify_id)
	os.Setenv("SPOTIFY_SECRET", spotify_secret)

	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryModify, spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState))

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			mono.ThreadAttach()
			Instance.ErrorString("Failed to start HTTP server: " + err.Error())
			return
		}
	}()

	url := auth.AuthURL(state)
	openbrowser(url)

	// wait for auth to complete
	SpotifyClient = <-ch
	srv.Close()

	user, err := SpotifyClient.CurrentUser(context.Background())
	if err != nil {
		Instance.ErrorString("Failed to get current user, Please make sure you have Spotify Premium!")
		return CSString("Failed to log in: " + err.Error())
	}

	Instance.MsgString("Native Module Initialized!")

	return CSString("You are logged in as: " + user.DisplayName)
}

func getString(ptr uintptr) string {
	str := syscall.UTF16ToString((*[1 << 30]uint16)(unsafe.Pointer(ptr))[:])
	return str
}

func UpdatePlayerInfo() {

	mono.ThreadAttach()

	mod, err := mono.GetAssemblyImage("VRCSpotifyIntegration")
	if err != nil {
		return
	}

	klass, err := mono.GetClass(mod, "VRCSpotifyMod", "VRCSpotifyMod")
	if err != nil {
		return
	}

	method, err := mono.GetMethod(klass, "UpdatePlayerInfo", 4)
	if err != nil {
		return
	}

	time.Sleep(time.Second / 4)
	args := make([]uintptr, 4)
	var handle1 uintptr
	var handle2 uintptr
	var handle3 uintptr
	var handle4 uintptr
	args[0], handle1, err = mono.NewString(Playing())
	args[1], handle2, err = mono.NewString(GetVolume())
	args[2], handle3, err = mono.NewString(GetRepeat())
	var favorited string
	if GetFavorited() {
		favorited = "Favorited!"
	} else {
		favorited = "Not Favorited"
	}
	args[3], handle4, err = mono.NewString(favorited)

	defer mono.mono_gchandle_free_(handle1)
	defer mono.mono_gchandle_free_(handle2)
	defer mono.mono_gchandle_free_(handle3)
	defer mono.mono_gchandle_free_(handle4)

	if err != nil {
		return
	}

	mono.RuntimeInvoke(method, 0, args)
}

func Playing() string {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return "<color=#03dffc>Vrc</color><color=#03fc3d>Spotify</color><color=#03dffc>Integration</color>"
	}
	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to Fetch Player State: " + err.Error())
		return "<color=#03dffc>Vrc</color><color=#03fc3d>Spotify</color><color=#03dffc>Integration</color>"
	}

	if !player.Playing {
		return "<color=#03dffc>Vrc</color><color=#03fc3d>Spotify</color><color=#03dffc>Integration</color>"
	}

	current, err := SpotifyClient.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		Instance.ErrorString("Failed to get currently playing Song: " + err.Error())
		return "<color=#03dffc>Vrc</color><color=#03fc3d>Spotify</color><color=#03dffc>Integration</color>"
	}
	artists := current.Item.Artists[0].Name
	if len(current.Item.Artists) > 1 {
		for i := 1; i < len(current.Item.Artists); i++ {
			concat := "&"
			if len(current.Item.Artists) > i {
				concat = ","
			}
			artists += concat + " " + current.Item.Artists[i].Name
		}
	}
	return artists + " - " + current.Item.Name
}

func GetVolume() string {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return "Volume: NaN"
	}
	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return "Volume: NaN"
	}

	return fmt.Sprintf("Volume: %v%%", player.Device.Volume)

}

func GetRepeat() string {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return "off"
	}
	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return "off"
	}

	return player.RepeatState
}

func GetFavorited() bool {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return false
	}
	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return false
	}

	if !player.Playing {
		return false
	}

	//check if current song is in the user's library
	current, err := SpotifyClient.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch currently playing: " + err.Error())
		return false
	}

	if current.Item.ID == "" {
		return false
	}

	//get user's library
	library, err := SpotifyClient.CurrentUsersTracks(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch User Library: " + err.Error())
		return false
	}

	for _, item := range library.Tracks {
		if item.ID == current.Item.ID {
			return true
		}
	}

	return false
}

//export RequestPlayerInfo
func RequestPlayerInfo() {
	UpdatePlayerInfo()
}

//export Favorite
func Favorite() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	current, err := SpotifyClient.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Currently Playing: " + err.Error())
		return
	}

	if current.Item.ID.String() == "" {
		return
	}

	//save the current item to the user's library, or remove it
	if GetFavorited() {
		SpotifyClient.RemoveTracksFromLibrary(ctx, current.Item.ID)
	} else {
		SpotifyClient.AddTracksToLibrary(ctx, current.Item.ID)
	}
	UpdatePlayerInfo()
}

//export Repeat
func Repeat() {
	//change the repeat state
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}
	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	//change the repeat state
	switch player.RepeatState {
	case "off":
		SpotifyClient.Repeat(ctx, "context")
	case "context":
		SpotifyClient.Repeat(ctx, "track")
	case "track":
		SpotifyClient.Repeat(ctx, "off")
	}

	UpdatePlayerInfo()
}

//export VolumeUp
func VolumeUp() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	//volume up
	err = SpotifyClient.Volume(ctx, player.Device.Volume+10)
	if err != nil {
		return
	}

	UpdatePlayerInfo()
}

//export VolumeDown
func VolumeDown() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	//volume down
	err = SpotifyClient.Volume(ctx, player.Device.Volume-10)
	if err != nil {
		Instance.ErrorString("Failed to set Volume: " + err.Error())
		return
	}

	UpdatePlayerInfo()
}

//export Play
func Play() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if player.Playing {
		return
	}

	//play
	err = SpotifyClient.Play(ctx)
	if err != nil {
		Instance.ErrorString("Failed to play: " + err.Error())
		return
	}

	UpdatePlayerInfo()
}

//export Pause
func Pause() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	//pause
	err = SpotifyClient.Pause(ctx)
	if err != nil {
		Instance.ErrorString("Failed to pause: " + err.Error())
		return
	}

	UpdatePlayerInfo()
}

//export Next
func Next() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	//next
	err = SpotifyClient.Next(ctx)
	if err != nil {
		Instance.ErrorString("Failed to skip: " + err.Error())
		return
	}

	UpdatePlayerInfo()
}

//export Previous
func Previous() {
	if SpotifyClient == nil {
		Instance.ErrorString("You are not logged in!")
		return
	}

	ctx := context.Background()

	//get player state
	player, err := SpotifyClient.PlayerState(ctx)
	if err != nil {
		Instance.ErrorString("Failed to fetch Player State: " + err.Error())
		return
	}

	if !player.Playing {
		return
	}

	//previous
	err = SpotifyClient.Previous(ctx)
	if err != nil {
		Instance.ErrorString("Failed to play previous: " + err.Error())
		return
	}

	UpdatePlayerInfo()
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		Instance.ErrorString("Couldn't get token: " + err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		Instance.ErrorString("State doesn't match: " + st)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func CSString(str string) uintptr {
	res, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		Instance.ErrorString("Failed to convert string to UTF16: " + err.Error())
		return 0
	}
	return uintptr(unsafe.Pointer(res))
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		Instance.ErrorString("Failed to open browser: " + err.Error())
	}

}

func main() {

}
