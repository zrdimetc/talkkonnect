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
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 * gps.go -> talkkonnect function to interface to USB GPS Neo6M
 */

package talkkonnect

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
	"github.com/talkkonnect/go-nmea"
)

//GPS Related Global Variables
var (
	GPSTime      string
	GPSDate      string
	GPSLatitude  float64
	GPSLongitude float64
	GPSSpeed     float64
	GPSCourse    float64
	GPSVariation float64
)

var goodGPSRead bool = false

func getGpsPosition(verbose bool) (bool, error) {
	if Config.Global.Hardware.GPS.Enabled {

		if Config.Global.Hardware.GPS.Port == "" {
			return false, errors.New("You Must Specify Port")
		}

		if Config.Global.Hardware.GPS.Even && Config.Global.Hardware.GPS.Odd {
			return false, errors.New("can't specify both even and odd parity")
		}

		parity := serial.PARITY_NONE

		if Config.Global.Hardware.GPS.Even {
			parity = serial.PARITY_EVEN
		} else if Config.Global.Hardware.GPS.Odd {
			parity = serial.PARITY_ODD
		}

		options := serial.OpenOptions{
			PortName:               Config.Global.Hardware.GPS.Port,
			BaudRate:               Config.Global.Hardware.GPS.Baud,
			DataBits:               Config.Global.Hardware.GPS.DataBits,
			StopBits:               Config.Global.Hardware.GPS.StopBits,
			MinimumReadSize:        Config.Global.Hardware.GPS.MinRead,
			InterCharacterTimeout:  Config.Global.Hardware.GPS.CharTimeOut,
			ParityMode:             parity,
			Rs485Enable:            Config.Global.Hardware.GPS.Rs485,
			Rs485RtsHighDuringSend: Config.Global.Hardware.GPS.Rs485HighDuringSend,
			Rs485RtsHighAfterSend:  Config.Global.Hardware.GPS.Rs485HighAfterSend,
		}

		f, err := serial.Open(options)

		if err != nil {
			Config.Global.Hardware.GPS.Enabled = false
			return false, errors.New("Cannot Open Serial Port")
		} else {
			defer f.Close()
		}

		if Config.Global.Hardware.GPS.TxData != "" {
			txData_, err := hex.DecodeString(Config.Global.Hardware.GPS.TxData)

			if err != nil {
				Config.Global.Hardware.GPS.Enabled = false
				return false, errors.New("Cannot Decode Hex Data")
			}

			log.Println("Sending: ", hex.EncodeToString(txData_))

			count, err := f.Write(txData_)

			if err != nil {
				return false, errors.New("Error writing to serial port")
			} else {
				log.Println("Wrote %v bytes\n", count)
			}

		}

		if Config.Global.Hardware.GPS.Rx {
			serialPort, err := serial.Open(options)
			if err != nil {
				log.Println("warn: Unable to Open Serial Port Error ", err)
			}

			defer serialPort.Close()

			reader := bufio.NewReader(serialPort)
			scanner := bufio.NewScanner(reader)

			goodGPSRead = false
			for scanner.Scan() {
				s, err := nmea.Parse(scanner.Text())

				if err == nil {
					if s.DataType() == nmea.TypeRMC {
						m := s.(nmea.RMC)
						if m.Latitude != 0 && m.Longitude != 0 {
							goodGPSRead = true
							GPSTime = fmt.Sprintf("%v", m.Time)
							GPSDate = fmt.Sprintf("%v", m.Date)
							GPSLatitude = m.Latitude
							GPSLongitude = m.Longitude
							if verbose {
								log.Println("info: Raw Sentence ", m)
								log.Println("info: Time: ", m.Time)
								log.Println("info: Validity: ", m.Validity)
								log.Println("info: Latitude GPS: ", nmea.FormatGPS(m.Latitude))
								log.Println("info: Latitude DMS: ", nmea.FormatDMS(m.Latitude))
								log.Println("info: Longitude GPS: ", nmea.FormatGPS(m.Longitude))
								log.Println("info: Longitude DMS: ", nmea.FormatDMS(m.Longitude))
								log.Println("info: Speed: ", m.Speed)
								log.Println("info: Course: ", m.Course)
								log.Println("info: Date: ", m.Date)
								log.Println("info: Variation: ", m.Variation)
							}
							break
						} else {
							log.Println("warn: Got Latitude 0 and Longtitude 0 from GPS")
						}
					} else {
						log.Println("warn: GPS Sentence Format Was not nmea.RMC")
					}
				}
			}
		} else {
			return false, errors.New("Rx Not Set")
		}

		return goodGPSRead, nil
	}
	return false, errors.New("GPS Not Enabled")
}
