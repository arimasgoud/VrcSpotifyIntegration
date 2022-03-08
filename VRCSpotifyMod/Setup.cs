using MelonLoader;
using ReMod.Core.VRChat;
using System.Collections;
using System.Diagnostics;
using System.Threading.Tasks;
using UnityEngine;
using VRC;
using VRC.SDKBase;
using VRC.UI;

namespace VRCSpotifyMod
{
    internal class Setup
    {
        public static VRCUiPopupManager Manager() => VRCUiPopupManager.prop_VRCUiPopupManager_0;
        public static void FirstTime()
        {
            Manager().ShowStandardPopupV2("Welcome to VRChatSpotifyIntegration!", "To get set up, please click Okay and follow the instructions in the ReadMe, then continue in VRChat.", "Okay", () =>
            {
                Process.Start("https://github.com/RinLovesYou/VrcSpotifyIntegration");
                MelonCoroutines.Start(SpotifyIdThing());
            });
        }

        public static IEnumerator SpotifyIdThing()
        {
            yield return new WaitForSeconds(1f);

            Manager().ShowInputPopup("VRCSpotifyIntegration", "Enter the Spotify Application ID", (s) =>
            {
                VRCSpotifyMod.Spotify_Id.Value = s;
                MelonCoroutines.Start(SpotifyTokenThing());
            });
        }

        public static IEnumerator SpotifyTokenThing()
        {
            yield return new WaitForSeconds(1f);

            Manager().ShowInputPopup("VRCSpotifyIntegration", "Enter the Spotify Application Secret", (s) =>
            {
                VRCSpotifyMod.Spotify_Secret.Value = s;
                MelonCoroutines.Start(SpotifyPortThing());
            });
        }

        public static IEnumerator SpotifyPortThing()
        {
            yield return new WaitForSeconds(1f);

            Manager().ShowInputPopup("VRCSpotifyIntegration", "Enter the Port", (s) =>
            {
                VRCSpotifyMod.Spotify_Port.Value = s;
                MelonPreferences.Save();
                var _isSetup = VRCSpotifyMod.Spotify_Id.Value != "" && VRCSpotifyMod.Spotify_Secret.Value != "";

                if (_isSetup)
                    Task.Run(() => VRCSpotifyMod.LogMeIn(VRCSpotifyMod.Spotify_Id.Value, VRCSpotifyMod.Spotify_Secret.Value, VRCSpotifyMod.Spotify_Port.Value));
            });
        }
    }
}
