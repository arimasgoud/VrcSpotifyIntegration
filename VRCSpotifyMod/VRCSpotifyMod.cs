using MelonLoader;
using ReMod.Core.Managers;
using ReMod.Core.VRChat;
using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using TMPro;
using UnityEngine;
using UnityEngine.UI;
using VRC.UI.Core;

namespace VRCSpotifyMod
{
    public class VRCSpotifyMod : MelonMod
    {
        public static MelonLogger.Instance Log = new MelonLogger.Instance("VRCSpotifyIntegration");

        public static MelonPreferences_Category MyPreferenceCategory;
        public static MelonPreferences_Entry<string> Spotify_Id;
        public static MelonPreferences_Entry<string> Spotify_Secret;
        public static MelonPreferences_Entry<string> Spotify_Port;

        public static NativeModule GotifyNative;

        private static bool _isSetup;

        public VRCSpotifyMod()
        {
            GotifyNative = new NativeModule();
            if (!GotifyNative.HasRemodCore())
            {
                Log.Warning("You do not have Remod.Core installed. Consider downloading it from https://github.com/RequiDev/Remod.Core/releases/latest and putting it in your VRChat root directory.\n");
                Log.Msg("Downloading & Injecting Remod.Core");

                LoadReModCore();
            }
        }

        public override void OnApplicationStart()
        {
            MyPreferenceCategory = MelonPreferences.CreateCategory("VrcSpotifyMod");
            Spotify_Id = MyPreferenceCategory.CreateEntry("Spotify_Id", "");
            Spotify_Secret = MyPreferenceCategory.CreateEntry("Spotify_Secret", "");
            Spotify_Port = MyPreferenceCategory.CreateEntry("Spotify_Port", "42069");

            MelonPreferences.Save();

            _isSetup = Spotify_Id.Value != "" && Spotify_Secret.Value != "";

            CacheIcons();

            if (_isSetup)
                MelonCoroutines.Start(LoginRoutine(Spotify_Id.Value, Spotify_Secret.Value, Spotify_Port.Value));
            
            OnUIManagerInitialized(() =>
            {
                RinMenu.PrepareMenu();
            });
        }

        private void OnUIManagerInitialized(Action code)
        {
            MelonCoroutines.Start(OnUiManagerInitCoroutine(code));
        }

        private IEnumerator OnUiManagerInitCoroutine(Action code)
        {
            while (VRCUiManager.prop_VRCUiManager_0 == null) yield return null;

            //early init

            while (UIManager.field_Private_Static_UIManager_0 == null)
                yield return null;
            while (GameObject.Find("UserInterface").GetComponentInChildren<VRC.UI.Elements.QuickMenu>(true) == null)
                yield return null;
            while (QuickMenuEx.Instance == null)
                yield return null;
            code();
            while (VRCPlayer.field_Internal_Static_VRCPlayer_0 == null)
                yield return null;
            if (!_isSetup)
            {
                yield return new WaitForSeconds(1f);
                Setup.FirstTime();
            }
        }

        public static IEnumerator LoginRoutine(string id, string token, string port)
        {
            var res = GotifyNative.Login(id, token, port);
            if (res == null)
            {
                Log.Error("Failed to get a response from the native Spotify module. Please try logging in again in-game!");
            }
            Log.Msg(res);
            yield break;
        }

        private void LoadReModCore()
        {
            try
            {
                var bytes = new WebClient().DownloadData("https://github.com/RequiDev/ReMod.Core/releases/latest/download/ReMod.Core.dll");
                Assembly.Load(bytes);
            }
            catch (Exception e)
            {
                MelonLogger.Error($"Unable to Load ReModCore Dependency: {e}");
            }
        }

        public static void UpdatePlayerInfo(string playing, string volume, string repeat, string favorite)
        {
            RinMenu.PlayingCategory.Header.Title = playing;
            RinMenu.VolumeCategory.Header.Title = volume;
            RinMenu.RepeatButton.Text = repeat;
            Sprite sprite = null;
            switch(repeat)
            {
                case "off":
                    sprite = ResourceManager.GetSprite("Spotify.repeat-off");
                    break;
                case "track":
                    sprite = ResourceManager.GetSprite("Spotify.repeat-once");
                    break;
                case "context":
                    sprite = ResourceManager.GetSprite("Spotify.repeat");
                    break;
            }

            var image = RinMenu.RepeatButton.RectTransform.Find("Icon").GetComponent<Image>();
            image.sprite = sprite;
            image.overrideSprite = sprite;

            RinMenu.FavoriteButton.Text = favorite;
            Sprite sprite2 = null;
            switch (favorite)
            {
                case "Favorited!":
                    sprite2 = ResourceManager.GetSprite("Spotify.heart");
                    break;
                default:
                    sprite2 = ResourceManager.GetSprite("Spotify.heart-outline");
                    break;
            }

            var image2 = RinMenu.FavoriteButton.RectTransform.Find("Icon").GetComponent<Image>();
            image2.sprite = sprite2;
            image2.overrideSprite = sprite2;
        }

        public static void CacheIcons()
        {
            //https://github.com/RequiDev/ReModCE/blob/master/ReModCE/ReMod.cs
            var ourAssembly = Assembly.GetExecutingAssembly();
            var resources = ourAssembly.GetManifestResourceNames();
            foreach (var resource in resources)
            {
                if (!resource.EndsWith(".png"))
                    continue;

                var stream = ourAssembly.GetManifestResourceStream(resource);

                var ms = new MemoryStream();
                stream.CopyTo(ms);
                var resourceName = Regex.Match(resource, @"([a-zA-Z\d\-_]+)\.png").Groups[1].ToString();
                ResourceManager.LoadSprite("Spotify", resourceName, ms.ToArray());
            }
        }
    }
}
