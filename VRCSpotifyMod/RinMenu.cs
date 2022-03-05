using ReMod.Core.Managers;
using ReMod.Core.UI.QuickMenu;
using System.Threading.Tasks;
using TMPro;
using UnityEngine.UI;
using static VRCSpotifyMod.VRCSpotifyMod;

namespace VRCSpotifyMod
{
    internal static class RinMenu
    {
        public static ReMenuCategory PlayingCategory;
        public static ReMenuCategory VolumeCategory;
        public static ReMenuButton RepeatButton;
        public static ReMenuButton FavoriteButton;
        public static void PrepareMenu()
        {
            var page = new ReCategoryPage("VrcSpotifyIntegration", true);
            ReTabButton.Create("SpotifyTab", "VRCSpotifyIntegration tab", "VrcSpotifyIntegration", ResourceManager.GetSprite("Spotify.logo"));
            PlayingCategory = page.AddCategory("<color=#03dffc>Vrc</color><color=#03fc3d>Spotify</color><color=#03dffc>Integration</color>", false);
            VolumeCategory = page.AddCategory("Volume: NaN");
            var text = PlayingCategory.Header.GameObject.GetComponentInChildren<TextMeshProUGUI>(true);
            text.enableAutoSizing = true;
            text.enableWordWrapping = false;


            page.OnOpen += () =>
            {
                Task.Run(() =>
                {
                    GotifyNative.RequestPlayerInfo();
                });
            };

            PlayingCategory.AddButton("Previous", "Go one song back", () => Task.Run(GotifyNative.Previous), ResourceManager.GetSprite("Spotify.last"));
            PlayingCategory.AddButton("Pause", "Pause Playback", () => Task.Run(GotifyNative.Pause), ResourceManager.GetSprite("Spotify.pause"));
            PlayingCategory.AddButton("Play", "Resume Playback", () => Task.Run(GotifyNative.Play), ResourceManager.GetSprite("Spotify.play"));
            PlayingCategory.AddButton("Next", "Go one song forward", () => Task.Run(GotifyNative.Next), ResourceManager.GetSprite("Spotify.next"));

            RepeatButton = VolumeCategory.AddButton("Repeat", "Change the Repeat State", () => Task.Run(GotifyNative.Repeat), ResourceManager.GetSprite("Spotify.repeat"));
            VolumeCategory.AddButton("Vol. Up", "Increase the Volume by 10%", () => Task.Run(GotifyNative.VolumeUp), ResourceManager.GetSprite("Spotify.plus"));
            VolumeCategory.AddButton("Vol. Down", "Decrease the Volume by 10%", () => Task.Run(GotifyNative.VolumeDown), ResourceManager.GetSprite("Spotify.minus"));
            FavoriteButton = VolumeCategory.AddButton("Favorite", "Favorite/Unfavorite the current song", () => Task.Run(GotifyNative.Favorite), ResourceManager.GetSprite("Spotify.heart"));

        }
    }
}
