using MelonLoader;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading.Tasks;

namespace VRCSpotifyMod
{
    public class NativeModule
    {
		[DllImport("kernel32", SetLastError = true, CharSet = CharSet.Ansi)]
		private static extern IntPtr LoadLibrary([MarshalAs(UnmanagedType.LPStr)] string lpFileName);
		[DllImport("kernel32.dll", SetLastError = true)]
		[return: MarshalAs(UnmanagedType.Bool)]
		private static extern bool FreeLibrary(IntPtr hModule);
		[DllImport("kernel32", CharSet = CharSet.Ansi, ExactSpelling = true, SetLastError = true)]
		private static extern IntPtr GetProcAddress(IntPtr hModule, string procName);

		public NativeModule()
        {
			try
			{
				WriteResourceToFile("VRCSpotifyMod.GotifyNative.dll", $"{MelonUtils.UserDataDirectory}/GotifyNative.dll");
			} catch(Exception e)
            {
				MelonLogger.Error(e);
				this._failedToLoad = true;
            }
			this._nativeModulePtr = LoadLibrary($"UserData/GotifyNative.dll");
			if (this._nativeModulePtr == IntPtr.Zero)
			{
				int lastWin32Error = Marshal.GetLastWin32Error();
				this._failedToLoad = true;
				throw new Exception("Can't load native loader: ", new Win32Exception(lastWin32Error)
				{
					Data =
					{
						{
							"LastWin32Error",
							lastWin32Error
						}
					}
				});
			}
			this._failedToLoad = false;
			this._hasRemodCore = this.GetExportedFunction<BoolFn>("HasRemodCore");
			this._login = this.GetExportedFunction<PtrFnArgs>("Login");
			this._play = this.GetExportedFunction<VoidFn>("Play");
			this._pause = this.GetExportedFunction<VoidFn>("Pause");
			this._next = this.GetExportedFunction<VoidFn>("Next");
			this._previous = this.GetExportedFunction<VoidFn>("Previous");
			this._volumeUp = this.GetExportedFunction<VoidFn>("VolumeUp");
			this._volumeDown = this.GetExportedFunction<VoidFn>("VolumeDown");
			this._requestPlayerInfo = this.GetExportedFunction<VoidFn>("RequestPlayerInfo");
			this._repeat = this.GetExportedFunction<VoidFn>("Repeat");
			this._favorite = this.GetExportedFunction<VoidFn>("Favorite");
		}

		public bool HasRemodCore()
        {
			if (_failedToLoad || _hasRemodCore == null)
            {
				return false;
            }
			return _hasRemodCore();
        }

		public string Login(string id, string secret, string port)
        {

			var idPtr = Marshal.StringToHGlobalAuto(id);
			var tokenPtr = Marshal.StringToHGlobalAuto(secret);
			var portPtr = Marshal.StringToHGlobalAuto(port);

			if (_failedToLoad || _login == null)
			{
				return null;
			}
			return UnmarshalString(_login(idPtr, tokenPtr, portPtr));
		}

		public void Play()
        {
			if (_failedToLoad || _play == null)
			{
				return;
			}
			_play();
		}

		public void Pause()
		{
			if (_failedToLoad || _pause == null)
			{
				return;
			}
			_pause();
		}

		public void Next()
		{
			if (_failedToLoad || _next == null)
			{
				return;
			}
			_next();
		}

		public void Previous()
		{
			if (_failedToLoad || _previous == null)
			{
				return;
			}
			_previous();
		}

		public void VolumeUp()
		{
			if (_failedToLoad || _volumeUp == null)
			{
				return;
			}
			_volumeUp();
		}

		public void VolumeDown()
		{
			if (_failedToLoad || _volumeDown == null)
			{
				return;
			}
			_volumeDown();
		}

		public void RequestPlayerInfo()
        {
			if (_failedToLoad || _requestPlayerInfo == null)
            {
				return;
            }
			_requestPlayerInfo();
        }

		public void Repeat()
		{
			if (_failedToLoad || _repeat == null)
			{
				return;
			}
			_repeat();
		}

		public void Favorite()
		{
			if (_failedToLoad || _favorite == null)
			{
				return;
			}
			_favorite();
		}

		private T GetExportedFunction<T>(string name)
		{
			T delegateForFunctionPointer = Marshal.GetDelegateForFunctionPointer<T>(GetProcAddress(_nativeModulePtr, name));
			if (delegateForFunctionPointer == null)
			{
				int lastWin32Error = Marshal.GetLastWin32Error();
				Win32Exception ex = new Win32Exception(lastWin32Error);
				ex.Data.Add("LastWin32Error", lastWin32Error);
				throw new Exception("Can't find exported function \"" + name + "\"", ex);
			}
			return delegateForFunctionPointer;
		}

		private void WriteResourceToFile(string resourceName, string fileName)
		{
			using (var resource = Assembly.GetExecutingAssembly().GetManifestResourceStream(resourceName))
			{
				using (var file = new FileStream(fileName, FileMode.Create, FileAccess.Write))
				{
					resource.CopyTo(file);
				}
			}
		}

		public static string UnmarshalString(IntPtr input)
		{
			string ret = "";
			try
			{
				ret = Marshal.PtrToStringAuto(input);
			}
			catch (Exception ex)
			{
				Console.WriteLine(ex);
				return null;
			}
			return ret;
		}

		private IntPtr _nativeModulePtr;

		private BoolFn _hasRemodCore;
		private PtrFnArgs _login;
		private VoidFn _play;
		private VoidFn _pause;
		private VoidFn _next;
		private VoidFn _previous;
		private VoidFn _volumeUp;
		private VoidFn _volumeDown;
		private VoidFn _requestPlayerInfo;
		private VoidFn _repeat;
		private VoidFn _favorite;

		private bool _failedToLoad;

		[UnmanagedFunctionPointer(CallingConvention.Cdecl, CharSet = CharSet.Ansi, SetLastError = true, BestFitMapping = true)]
		private delegate bool BoolFn();
		[UnmanagedFunctionPointer(CallingConvention.Cdecl)]
		private delegate void VoidFn();
		[UnmanagedFunctionPointer(CallingConvention.Cdecl)]
		private delegate IntPtr PtrFn();
		[UnmanagedFunctionPointer(CallingConvention.Cdecl)]
		private delegate IntPtr PtrFnArgs(IntPtr id, IntPtr token, IntPtr port);
	}
}
