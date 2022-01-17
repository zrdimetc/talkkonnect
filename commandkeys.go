/*
 * talkkonnect headless mumble client/gateway with lcd screen and channel control
 * Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * talkkonnect is the based on talkiepi and barnard by Daniel Chote and Tim Cooper
 *
 * The Initial Developer of the Original Code is
 * Suvir Kumar <suvir@talkkonnect.com>
 * Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
 *
 * Contributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 * Zoran Dimitrijevic
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 *
 */

package talkkonnect

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/talkkonnect/volume-go"
)

func (b *Talkkonnect) cmdDisplayMenu() {
	log.Println("debug: Delete Key Pressed Menu and Session Information Requested")

	TTSEvent("displaymenu")
	b.talkkonnectMenu("\u001b[44;1m") // add blue background to banner reference https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#background-colors
	b.ParticipantLEDUpdate(true)
}

func (b *Talkkonnect) cmdChannelUp(source string) {
	log.Printf("debug: F1 pressed Channel Up (+) Requested from %v\n", source)
	b.ChannelUp()
}

func (b *Talkkonnect) cmdChannelDown(source string) {
	log.Printf("debug: F2 pressed Channel Down (-) Requested from %v\n", source)
	b.ChannelDown()
}

func (b *Talkkonnect) cmdMuteUnmute(subCommand, source string) {

	log.Printf("debug: F3 pressed %v Speaker Requested from %v\n", subCommand, source)
	OrigMuted, err := volume.GetMuted(Config.Global.Software.Settings.OutputMuteControlDevice)

	if err != nil {
		log.Println("error: Unable to get current Muted/Unmuted State ", err)
	} else {
		if OrigMuted {
			log.Println("debug: Originally Device is Muted")
		} else {
			log.Println("debug: Originally Device is Unmuted")
		}
	}

	if subCommand == "toggle" {
		if OrigMuted {
			err := volume.Unmute(Config.Global.Software.Settings.OutputMuteControlDevice)
			if err != nil {
				log.Println("error: Unmuting Failed", err)
				return
			}
			TTSEvent("unmutespeaker")
			log.Println("info: Output Device Unmuted")
			if Config.Global.Hardware.TargetBoard == "rpi" {
				if LCDEnabled {
					LcdText = [4]string{"nil", "nil", "nil", "UnMuted"}
					LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
				}
				if OLEDEnabled {
					oledDisplay(false, 6, 1, "Unmuted")
				}
			}
			return
		} else {
			TTSEvent("mutespeaker")
			err = volume.Mute(Config.Global.Software.Settings.OutputMuteControlDevice)
			if err != nil {
				log.Println("error: Muting Failed", err)
			}
			log.Println("info: Output Device Muted")
			if Config.Global.Hardware.TargetBoard == "rpi" {
				if LCDEnabled {
					LcdText = [4]string{"nil", "nil", "nil", "Muted"}
					LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
				}
				if OLEDEnabled {
					oledDisplay(false, 6, 1, "Muted")
				}
			}
			return
		}
	}

	//force mute
	if subCommand == "mute" {
		TTSEvent("mutespeaker")
		err = volume.Mute(Config.Global.Software.Settings.OutputMuteControlDevice)
		if err != nil {
			log.Println("error: Muting Failed ", err)
			return
		}
		log.Println("info: Output Device Muted")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "Muted"}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)

				if OLEDEnabled {
					oledDisplay(false, 6, 1, "Muted")
				}
			}
			return
		}
	}
	//force unmute
	if subCommand == "unmute" {
		err := volume.Unmute(Config.Global.Software.Settings.OutputMuteControlDevice)
		TTSEvent("unmutespeaker")
		if err != nil {
			log.Println("error: Unmute Failed ", err)
			return
		}
		log.Println("info: Output Device Unmuted")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "UnMuted"}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(false, 6, 1, "Unmuted")
			}
		}
		return
	}
}
func (b *Talkkonnect) cmdCurrentVolume(source string) {
	OrigVolume, err := volume.GetVolume(Config.Global.Software.Settings.OutputVolControlDevice)
	if err != nil {
		log.Printf("error: Unable to get current volume: %+v\n", err)
	}

	log.Printf("debug: F4 pressed Volume Level Requested From =%v\n", source)
	log.Println("info: Volume Level is at", OrigVolume, "%")

	TTSEvent("currentvolumelevel")
	if Config.Global.Hardware.TargetBoard == "rpi" {
		if LCDEnabled {
			LcdText = [4]string{"nil", "nil", "nil", "Volume " + strconv.Itoa(OrigVolume)}
			LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
		}
		if OLEDEnabled {
			oledDisplay(false, 6, 1, "Volume "+strconv.Itoa(OrigVolume))
		}

	}
}

func (b *Talkkonnect) cmdVolumeUp(source string) {
	log.Printf("debug: F5 pressed Volume UP (+) from %v\n", source)
	origVolume, err := volume.GetVolume(Config.Global.Software.Settings.OutputVolControlDevice)
	if err != nil {
		log.Printf("warn: unable to get original volume: %+v volume control will not work!\n", err)
		return
	}

	if origVolume < 100 {
		err := volume.IncreaseVolume(Config.Global.Hardware.IO.VolumeButtonStep.VolUpStep, Config.Global.Software.Settings.OutputVolControlDevice)
		if err != nil {
			log.Println("warn: F5 Increase Volume Failed! ", err)
		}
		origVolume, _ := volume.GetVolume(Config.Global.Software.Settings.OutputVolControlDevice)
		log.Println("info: Volume UP (+) Now At ", origVolume, "%")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "Volume + " + strconv.Itoa(origVolume)}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(false, 6, 1, "Volume "+strconv.Itoa(origVolume))
			}
		}
	} else {
		log.Println("debug: F5 Increase Volume")
		log.Println("info: Already at Maximum Possible Volume")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "Max Vol"}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(false, 6, 1, "Max Vol")
			}
		}
	}
	TTSEvent("digitalvolumeup")
}

func (b *Talkkonnect) cmdVolumeDown(source string) {
	log.Printf("info: F6 pressed Volume Down (-) from %v\n", source)
	origVolume, err := volume.GetVolume(Config.Global.Software.Settings.OutputVolControlDevice)
	if err != nil {
		log.Printf("warn: unable to get original volume: %+v volume control will not work!\n", err)
		return
	}

	if origVolume > 0 {
		err := volume.IncreaseVolume(Config.Global.Hardware.IO.VolumeButtonStep.VolDownStep, Config.Global.Software.Settings.OutputVolControlDevice)
		if err != nil {
			log.Println("error: F6 Decrease Volume Failed! ", err)
		}
		origVolume, _ := volume.GetVolume(Config.Global.Software.Settings.OutputVolControlDevice)
		log.Println("info: Volume Down (-) Now At ", origVolume, "%")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "Volume - " + strconv.Itoa(origVolume)}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(false, 6, 1, "Volume "+strconv.Itoa(origVolume))
			}

		}
	} else {
		log.Println("debug: F6 Increase Volume Already")
		log.Println("info: Already at Minimum Possible Volume")
		if Config.Global.Hardware.TargetBoard == "rpi" {
			if LCDEnabled {
				LcdText = [4]string{"nil", "nil", "nil", "Min Vol"}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(false, 6, 1, "Min Vol")
			}
		}
	}
	TTSEvent("digitalvolumedown")
}

func (b *Talkkonnect) cmdListServerChannels(source string) {
	log.Printf("debug: F7 pressed Channel List Requested %v\n", source)

	TTSEvent("listserverchannels")
	b.ListChannels(true)
	b.ParticipantLEDUpdate(true)
}

func (b *Talkkonnect) cmdStartTransmitting(source string) {
	log.Printf("debug: F8 pressed TX Mode Requested (Start Transmitting) from %v\n", source)
	log.Println("info: Start Transmitting")

	TTSEvent("starttransmitting")

	if IsPlayStream {
		IsPlayStream = false
		NowStreaming = false

		var eventSound EventSoundStruct = findEventSound("stream")
		if eventSound.Enabled {
			if s, err := strconv.ParseFloat(eventSound.Volume, 32); err == nil {
				b.playIntoStream(eventSound.FileName, float32(s))
			}
		}
	}

	if !b.IsTransmitting {
		time.Sleep(100 * time.Millisecond)
		b.TransmitStart()
	} else {
		log.Println("error: Already in Transmitting Mode")
	}
}

func (b *Talkkonnect) cmdStopTransmitting(source string) {
	log.Printf("debug: F9 pressed RX Mode Request (Stop Transmitting) from %v\n", source)
	log.Println("info: Stop Transmitting")

	TTSEvent("stoptransmitting")

	if IsPlayStream {
		IsPlayStream = false
		NowStreaming = false

		var eventSound EventSoundStruct = findEventSound("stream")
		if eventSound.Enabled {
			if s, err := strconv.ParseFloat(eventSound.Volume, 32); err == nil {
				b.playIntoStream(eventSound.FileName, float32(s))
			}
		}
	}

	if b.IsTransmitting {
		time.Sleep(100 * time.Millisecond)
		b.TransmitStop(true)
	} else {
		log.Println("info: Not Already Transmitting")
	}
}

func (b *Talkkonnect) cmdListOnlineUsers(source string) {
	log.Printf("debug: F10 pressed Online User(s) in Current Channel Requested from %v\n", source)
	log.Println("info: F10 Online User(s) in Current Channel")

	TTSEvent("listonlineusers")

	log.Println(fmt.Sprintf("info: Channel %#v Has %d Online User(s)", b.Client.Self.Channel.Name, len(b.Client.Self.Channel.Users)))
	b.ListUsers()
	b.ParticipantLEDUpdate(true)
}

func (b *Talkkonnect) cmdPlayback(source string) {
	log.Printf("debug: F11 pressed Start/Stop Stream Stream into Current Channel Requested from %v", source)
	log.Println("info: Stream into Current Channel")

	b.BackLightTimer()

	TTSEvent("playstream")

	if b.IsTransmitting {
		log.Println("alert: talkkonnect was already transmitting will now stop transmitting and start the stream")
		b.TransmitStop(false)
	}

	IsPlayStream = !IsPlayStream
	NowStreaming = IsPlayStream

	if IsPlayStream && Config.Global.Software.Settings.StreamSendMessage {
		b.SendMessage(fmt.Sprintf("%s Streaming", b.Username), false)
	}

	var eventSound EventSoundStruct = findEventSound("stream")
	if eventSound.Enabled {
		if s, err := strconv.ParseFloat(eventSound.Volume, 32); err == nil {
			go b.playIntoStream(eventSound.FileName, float32(s))
		}
	}
}

func (b *Talkkonnect) cmdGPSPosition(source string) {
	log.Printf("debug: F12 pressed from %v\n", source)
	log.Println("info: GPS details requested")

	TTSEvent("requestgpsposition")

	var i int = 0
	var tries int = 10

	for i = 0; i < tries; i++ {
		goodGPSRead, err := getGpsPosition(3)

		if err != nil {
			log.Println("error: GPS Function Returned Error Message", err)
			break
		}

		if goodGPSRead {
			break
		}

	}

	if i == tries {
		log.Println("warn: Could Not Get a Good GPS Read")
	}

}

func (b *Talkkonnect) cmdQuitTalkkonnect(source string) {
	log.Printf("debug: Ctrl-C Terminate Program Requested from %v\n", source)
	duration := time.Since(StartTime)
	log.Printf("info: Talkkonnect Now Running For %v \n", secondsToHuman(int(duration.Seconds())))
	TTSEvent("quittalkkonnect")
	CleanUp()
}

func (b *Talkkonnect) cmdDebugStacktrace() {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	var debug bytes.Buffer
	debug.WriteString(string(buf[0:stackSize]))
	scanner := bufio.NewScanner(&debug)
	var line int = 1
	log.Println("debug: Pressed Ctrl-D")
	log.Println("info: Stack Dump Requested")
	for scanner.Scan() {
		log.Printf("debug: line: %d %s\n", line, scanner.Text())
		line++
	}
	goStreamStats()
}

func (b *Talkkonnect) cmdSendEmail(source string) {
	log.Printf("debug: Ctrl-E Pressed from %v\n", source)
	log.Println("info: Send Email Requested")

	var i int = 0
	var tries int = 10

	for i = 0; i < tries; i++ {
		goodGPSRead, err := getGpsPosition(3)

		if err != nil {
			log.Println("error: GPS Function Returned Error Message", err)
			break
		}

		if goodGPSRead {
			break
		}

	}

	if i == tries {
		log.Println("warn: Could Not Get a Good GPS Read")
		return
	}

	TTSEvent("sendemail")

	if Config.Global.Software.SMTP.Enabled {

		emailMessage := fmt.Sprintf(Config.Global.Software.SMTP.Message + "\n")
		emailMessage = emailMessage + fmt.Sprintf("Ident: %s \n", b.Ident)
		emailMessage = emailMessage + fmt.Sprintf("Mumble Username: %s \n", b.Username)

		if Config.Global.Software.SMTP.GpsDateTime {
			emailMessage = emailMessage + fmt.Sprintf("Date "+GNSSData.Date+" UTC Time "+GNSSData.Time+"\n")
		}

		if Config.Global.Software.SMTP.GpsLatLong {
			emailMessage = emailMessage + fmt.Sprintf("Latitude "+strconv.FormatFloat(GNSSData.Lattitude, 'f', 6, 64)+" Longitude "+strconv.FormatFloat(GNSSData.Longitude, 'f', 6, 64)+"\n")
		}

		if Config.Global.Software.SMTP.GoogleMapsURL {
			emailMessage = emailMessage + "http://www.google.com/maps/place/" + strconv.FormatFloat(GNSSData.Lattitude, 'f', 6, 64) + "," + strconv.FormatFloat(GNSSData.Longitude, 'f', 6, 64)
		}

		err := sendviagmail(Config.Global.Software.SMTP.Username, Config.Global.Software.SMTP.Password, Config.Global.Software.SMTP.Receiver, Config.Global.Software.SMTP.Subject, emailMessage)
		if err != nil {
			log.Println("error: Error from Email Module: ", err)
		}
	} else {
		log.Println("warning: Sending Email Disabled in Config")
	}
}

func (b *Talkkonnect) cmdConnPreviousServer(source string) {
	log.Printf("debug: Ctrl-F Pressed from %v", source)
	log.Println("info: Previous Server Requested")

	TTSEvent("previousserver")

	if AccountCount > 1 {
		if AccountIndex > 0 {
			AccountIndex--
		} else {
			AccountIndex = AccountCount - 1
		}
		modifyXMLTagServerHopping(ConfigXMLFile, AccountIndex)
	}
}

func (b *Talkkonnect) cmdClearScreen(source string) {
	reset()
	log.Printf("debug: Ctrl-L Pressed Cleared Screen from %v\n", source)
	if Config.Global.Hardware.TargetBoard == "rpi" {
		if LCDEnabled {
			LcdText = [4]string{"nil", "nil", "nil", "nil"}
			LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
		}

		if OLEDEnabled {
			Oled.DisplayOn()
			LCDIsDark = false
			oledDisplay(true, 0, 0, "")
		}
	}
}

func (b *Talkkonnect) cmdPingServers(source string) {
	log.Printf("debug: Ctrl-O Pressed %v\n", source)
	log.Println("info: Ping Servers")
	TTSEvent("pingservers")
	b.pingServers()
}

func (b *Talkkonnect) cmdConnNextServer(source string) {
	log.Printf("debug: Ctrl-N Pressed from %v\n", source)
	log.Println("info: Next Server Requested Killing This Session, talkkonnect should be restarted by systemd")

	TTSEvent("nextserver")

	if AccountCount > 1 {
		if AccountIndex < AccountCount-1 {
			AccountIndex++
		} else {
			AccountIndex = 0
		}
		modifyXMLTagServerHopping(ConfigXMLFile, AccountIndex)
	}

}

func (b *Talkkonnect) cmdAudioTrafficRecord(source string) {
	log.Printf("debug: Ctrl-I Pressed from %v\n", source)
	log.Println("info: Traffic Recording Requested")
	if !Config.Global.Hardware.AudioRecordFunction.Enabled {
		log.Println("warn: Audio Recording Function Not Enabled")
	}
	if Config.Global.Hardware.AudioRecordFunction.RecordMode != "traffic" {
		log.Println("warn: Traffic Recording Not Enabled")
	}

	if Config.Global.Hardware.AudioRecordFunction.Enabled {
		if Config.Global.Hardware.AudioRecordFunction.RecordMode == "traffic" {
			if Config.Global.Hardware.AudioRecordFunction.RecordFromOutput != "" {
				if Config.Global.Hardware.AudioRecordFunction.RecordSoft == "sox" {
					go AudioRecordTraffic()
					if Config.Global.Hardware.TargetBoard == "rpi" {
						if LCDEnabled {
							LcdText = [4]string{"nil", "nil", "Traffic Audio Rec ->", "nil"} // 4 or 3
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled {
							oledDisplay(false, 5, 1, "Traffic Audio Rec ->") // 6 or 5
						}
					}
				} else {
					log.Println("info: Traffic Recording is not Enabled or sox Encountered Problems")
				}
			}
		}
	}
}

func (b *Talkkonnect) cmdAudioMicRecord(source string) {
	log.Printf("debug: Ctrl-J Pressed from %v\n", source)
	log.Println("info: Ambient (Mic) Recording Requested")
	if !Config.Global.Hardware.AudioRecordFunction.Enabled {
		log.Println("warn: Audio Recording Function Not Enabled")
	}
	if Config.Global.Hardware.AudioRecordFunction.RecordMode != "ambient" {
		log.Println("warn: Ambient (Mic) Recording Not Enabled")
	}

	if Config.Global.Hardware.AudioRecordFunction.Enabled {
		if Config.Global.Hardware.AudioRecordFunction.RecordMode == "ambient" {
			if Config.Global.Hardware.AudioRecordFunction.RecordFromInput != "" {
				if Config.Global.Hardware.AudioRecordFunction.RecordSoft == "sox" {
					go AudioRecordAmbient()
					if Config.Global.Hardware.TargetBoard == "rpi" {
						if LCDEnabled {
							LcdText = [4]string{"nil", "nil", "Mic Audio Rec ->", "nil"} // 4 or 3
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled {
							oledDisplay(false, 5, 1, "Mic Audio Rec ->") // 6 or 5
						}
					}
				} else {
					log.Println("error: Ambient (Mic) Recording is not Enabled or sox Encountered Problems")
				}
			}
		}
	}
}

func (b *Talkkonnect) cmdAudioMicTrafficRecord(source string) {
	log.Printf("debug: Ctrl-K Pressed from %v\n", source)
	log.Println("info: Recording (Traffic and Mic) Requested")
	if !Config.Global.Hardware.AudioRecordFunction.Enabled {
		log.Println("warn: Audio Recording Function Not Enabled")
	}
	if Config.Global.Hardware.AudioRecordFunction.RecordMode != "combo" {
		log.Println("warn: Combo Recording (Traffic and Mic) Not Enabled")
	}

	if Config.Global.Hardware.AudioRecordFunction.Enabled {
		if Config.Global.Hardware.AudioRecordFunction.RecordMode == "combo" {
			if Config.Global.Hardware.AudioRecordFunction.RecordFromInput != "" {
				if Config.Global.Hardware.AudioRecordFunction.RecordSoft == "sox" {
					go AudioRecordCombo()
					if Config.Global.Hardware.TargetBoard == "rpi" {
						if LCDEnabled {
							LcdText = [4]string{"nil", "nil", "Combo Audio Rec ->", "nil"} // 4 or 3
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled {
							oledDisplay(false, 5, 1, "Combo Audio Rec ->") // 6 or 5
						}
					}
				} else {
					log.Println("error: Combo Recording (Traffic and Mic) is not Enabled or sox Encountered Problems")
				}
			}
		}
	}
}

func (b *Talkkonnect) cmdPanicSimulation(source string) {
	if !(IsConnected) {
		return
	}
	b.BackLightTimer()
	log.Printf("debug: Ctrl-P Pressed from %v\n", source)
	log.Println("info: Panic Button Start/Stop Simulation Requested")

	TTSEvent("panicsimulation")

	if Config.Global.Hardware.PanicFunction.Enabled {

		if b.IsTransmitting {
			b.TransmitStop(false)
		} else {
			b.IsTransmitting = true
			b.SendMessage(Config.Global.Hardware.PanicFunction.Message, Config.Global.Hardware.PanicFunction.RecursiveSendMessage)

		}

		if Config.Global.Hardware.PanicFunction.SendIdent {
			b.SendMessage(fmt.Sprintf("My Username is %s and Ident is %s", b.Username, b.Ident), Config.Global.Hardware.PanicFunction.RecursiveSendMessage)
		}

		if Config.Global.Hardware.PanicFunction.SendGpsLocation {

			var i int = 0
			var tries int = 10

			for i = 0; i < tries; i++ {
				goodGPSRead, err := getGpsPosition(3)

				if err != nil {
					log.Println("error: GPS Function Returned Error Message", err)
					break
				}

				if goodGPSRead {
					break
				}
			}

			if i == tries {
				log.Println("warn: Could Not Get a Good GPS Read")
			}

			if goodGPSRead && i != tries {
				log.Println("info: Sending GPS Info My Message")
				gpsMessage := "My GPS Coordinates are " + fmt.Sprintf(" Latitude "+strconv.FormatFloat(GNSSData.Lattitude, 'f', 6, 64)) + fmt.Sprintf(" Longitude "+strconv.FormatFloat(GNSSData.Longitude, 'f', 6, 64))
				b.SendMessage(gpsMessage, Config.Global.Hardware.PanicFunction.RecursiveSendMessage)
			}

			IsPlayStream = true
			b.playIntoStream(Config.Global.Hardware.PanicFunction.FilenameAndPath, Config.Global.Hardware.PanicFunction.Volume)
			if Config.Global.Hardware.TargetBoard == "rpi" {
				if LCDEnabled {
					LcdText = [4]string{"nil", "nil", "nil", "Panic Message Sent!"}
					LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
				}
				if OLEDEnabled {
					oledDisplay(false, 6, 1, "Panic Message Sent!")
				}
			}
			if Config.Global.Hardware.PanicFunction.TxLockEnabled && Config.Global.Hardware.PanicFunction.TxLockTimeOutSecs > 0 {
				b.TxLockTimer()
			}

			// New. Send email after Panic Event //
			if Config.Global.Hardware.PanicFunction.PMailEnabled {
				b.cmdSendEmail("from panic")
				log.Println("info: Sending Panic Alert Email To Predefined Email Address")
			}
			//

			// New. Record ambient audio on Panic Event if recording is enabled
			if Config.Global.Hardware.AudioRecordFunction.Enabled {
				log.Println("info: Running sox for Audio Recording...")
				AudioRecordAmbient()
			}
			//

		} else {
			log.Println("warn: Panic Function Disabled in Config")
		}
		IsPlayStream = false
		b.IsTransmitting = false

		if Config.Global.Hardware.PanicFunction.PLowProfile {
			GPIOOutAll("led/relay", "off")
			log.Println("info: Low Profile Lights Option is Enabled. Turning All Leds Off During Panic Event")
			if LCDEnabled {
				log.Println("info: Low Profile Lights is Enabled. Turning Off Display During Panic Event")
				LcdText = [4]string{"", "", "", ""}
				LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
			}
			if OLEDEnabled {
				oledDisplay(true, 0, 0, "")
			}
		}
	}
}

func (b *Talkkonnect) cmdRepeatTxLoop(source string) {
	log.Printf("debug: Ctrl-R Pressed from %v\n", source)
	log.Println("info: Repeat TX Test Requested")
	isrepeattx = !isrepeattx
	go b.repeatTx()
}

func (b *Talkkonnect) cmdScanChannels(source string) {
	log.Printf("debug: Ctrl-S Pressed fgrom %v\n", source)
	log.Println("info: Scanning Channels")

	TTSEvent("startscanning")
	b.Scan()
}

func cmdThanks(source string) {
	log.Printf("debug: Ctrl-T Pressed from %v\n", source)
	log.Println("info: Thanks and Acknowledgements Screen Request ")
	talkkonnectAcknowledgements("\u001b[44;1m") // add blue background to banner reference https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#background-colors
}

func (b *Talkkonnect) cmdShowUptime(source string) {
	log.Printf("debug: Ctrl-U Pressed from %v\n", source)
	log.Println("info: Talkkonnect Uptime Request ")
	duration := time.Since(StartTime)
	log.Printf("info: Talkkonnect Now Running For %v \n", secondsToHuman(int(duration.Seconds())))
}

func (b *Talkkonnect) cmdDisplayVersion(source string) {
	log.Printf("debug: Ctrl-V Pressed from %v\n", source)
	log.Println("info: Talkkonnect Version Request ")
	releasedVersion := checkGitHubVersion()
	if talkkonnectVersion != releasedVersion {
		log.Printf("warn: Ver %v Rel %v (Different Ver %v Available!)\n", talkkonnectVersion, talkkonnectReleased, releasedVersion)
	} else {
		log.Printf("info: Ver %v Rel %v (Latest Release)\n", talkkonnectVersion, talkkonnectReleased)
	}
}

func (b *Talkkonnect) cmdDumpXMLConfig(source string) {
	log.Printf("debug: Ctrl-X Pressed from %v\n", source)
	log.Println("info: Print XML Config " + ConfigXMLFile)
	TTSEvent("printxmlconfig")
	printxmlconfig()
}

func (b *Talkkonnect) cmdPlayRepeaterTone(source string) {
	log.Printf("debug: Ctrl-G Pressed from %v\n", source)
	log.Println("info: Play Repeater Tone on Speaker and Simulate RX Signal")

	b.BackLightTimer()

	if Config.Global.Software.Sounds.RepeaterTone.Enabled {
		b.PlayTone(Config.Global.Software.Sounds.RepeaterTone.ToneFrequencyHz, Config.Global.Software.Sounds.RepeaterTone.ToneDurationSec, "local", true)
	} else {
		log.Println("warn: Repeater Tone Disabled by Config")
	}

}

func (b *Talkkonnect) cmdLiveReload(source string) {
	log.Printf("debug: Ctrl-B Pressed from %v\n", source)
	log.Println("info: XML Config Live Reload")
	err := readxmlconfig(ConfigXMLFile, true)
	if err != nil {
		message := err.Error()
		FatalCleanUp(message)
	}
}

func cmdSanityCheck(source string) {
	log.Printf("debug: Ctrl-H Pressed from %v\n", source)
	log.Println("info: XML Sanity Checker")
	CheckConfigSanity(false)
}
