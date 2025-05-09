/* NEC protocol documentation is here: https://assets.sharpnecdisplays.us/documents/usermanuals/external_controls.pdf */

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	/* "encoding/json"
	"time" */
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

// Protocol for our models is documented at
// https://assets.sharpnecdisplays.us/documents/usermanuals/external_control_e705-805.pdf
func setPower(socketKey string, value string) (string, error) {
	function := "setPower"
	var command string
	var statusStr string

	framework.Log(function + " - setting: " + socketKey + " to: " + value)
	if value == `"on"` {
		command = HexToBin("01304130413043024332303344363030303103730d")

	} else if value == `"off"` {
		command = HexToBin("01304130413043024332303344363030303403760d")
	} else {
		errMsg := fmt.Sprintf(function + " - aknke8e error illegal value for on or off: " + value)
		framework.Log(errMsg)
		return command, errors.New(errMsg)
	}

	framework.Log(fmt.Sprintf(function+" - asdasnw34 command is: %s", command))

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - r3fad5 error writing to " + socketKey + "closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - vcxxr343 still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return string(command), errors.New(errMsg)
		}
	}

	// get server response
	res := framework.ReadLineFromSocket(socketKey)

	raw := ExtractBody(BinToHex(res))

	framework.Log(function + "raw hex response is: " + raw)

	if !strings.Contains(raw, "303043323033443630303031") &&
		!strings.Contains(raw, "303043323033443630303034") {
		errMsg := function + " - gre654gf got invalid response: " + raw
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid device response")
	}

	statusStr = `"ok"` // once we get here, we know we succeeded

	return statusStr, nil
}

func getPower(socketKey string) (string, error) {
	function := "getPower"
	command := HexToBin("01304130433036023030443603770d")

	framework.Log(function + " - getting power status for: " + socketKey)

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - grgrdf error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - ndgrt still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return "", errors.New(errMsg)
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	powerStateReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + powerStateReply)

	var result string

	// Check for power on responses
	if strings.Contains(powerStateReply, "30303030443630313030303430303031") ||
		strings.Contains(powerStateReply, "30303030443630303030303430303031") {
		result = `"on"`
		// Check for power off responses
	} else if strings.Contains(powerStateReply, "30303030443630313030303430303034") ||
		strings.Contains(powerStateReply, "30303030443630303030303430303034") {
		result = `"off"`

		// Unknown response handling
	} else {
		errMsg := function + " - btrsdb unknown getpower response: " + powerStateReply
		framework.AddToErrors(socketKey, errMsg)
		result = `"unknown power response"`
		return result, errors.New(errMsg)
	}

	return result, nil
}

func setVideoRoute(socketKey string, outputNumber string, inputNumber string) (string, error) {
	function := "setInput"
	var command string
	var result string

	// outputNumber is ignored in this device since there is only one output (the screen)
	if outputNumber != "1" {
		errMsg := function + " - avegwava invalid output number value: " + outputNumber + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	framework.Log(function + " - setting input: " + inputNumber + " for: " + socketKey)

	inputNumber = strings.Trim(inputNumber, `"`) // trim quotes off JSON body input string
	switch inputNumber {
	case "1":
		command = HexToBin("0130413045304102303036303030313103720d")
	case "2":
		command = HexToBin("0130413045304102303036303030313203710d")
	default:
		errMsg := function + " - trying to set unknown input: " + inputNumber
		framework.Log(errMsg)
		return errMsg, errors.New("bad input provided error")
	}

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - efwrfc error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - gjkiyu still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("write error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + rawReply)

	// The differences in the codes below are from NEC's "maximum value that monitor can accept".
	// One could check more thoroughly by masking bytes 10 through 13 in the reply.
	// For now, we're simply adding strings as we encounter new responses.
	// This effect will happen in other set operations.
	if rawReply == "30303030363030303030383030303131" ||
		rawReply == "30303030363030303030383230303131" ||
		rawReply == "30303030363030303030383830303131" ||
		rawReply == "30303030363030313030383030303131" {
		//	rawReply == "303030303630303030303832303031" ||
		//	rawReply == "303030303630303030303830303031" {
		result = `"1"`
	} else if rawReply == "30303030363030303030383030303132" ||
		rawReply == "30303030363030303030383830303132" ||
		rawReply == "30303030363030303030383230303132" {
		result = `"2"`
	} else {
		errMsg := function + " - vv4ds unknown response: " + rawReply
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("response error")
	}

	return result, nil
}

func getVideoRoute(socketKey string, outputNumber string) (string, error) {
	function := "getInput"
	//      01304130433036023030363003030D (from PHP)
	//           01304130433036023032324503700D  (from NEC Tool)
	command := HexToBin("01304130433036023030363003030d")

	framework.Log(function + " - getting input status for: " + socketKey)

	// outputNumber is ignored in this device since there is only one output (the screen)
	if outputNumber != "1" {
		errMsg := function + " - btwrff invalid output number value: " + outputNumber + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - btsrfg error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - drbtfdbs still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("write error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	//rawReply := strings.ToLower(BinToHex(response))
	hexReply := BinToHex(response)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - hexReply is: " + hexReply + " rawReply is: " + rawReply)

	var result string

	switch rawReply {
	/* case "303030303630303030303830303031310309":
		result = `"1"`
	case "30303030363030303030383230303131030b":
		result = `"1"`
	case "303030303630303030303838303031310303":
		result = `"1"`
	case "30303030363030303030383030303132030a":
		result = `"2"`
	case "303030303630303030303832303031320308":
		//01 30 30 41 44 31 32 02 30 30 30 30 36 30 30 30 30 30 38 32 30 30 31 32 03 08
		result = `"2"`
	case "303030303630303030303838303031320300":
		result = `"2"`
	case "3030303036303030303038303030313103":
		result = `"2"` */
	/* 303030303630303030303830303031
	303030303630303030303830303031
	303030303630303030303830303031
	303030303630303030303830303031
	303030303630303030303830303031
	30303032324530303030304130303034 from NEC tool */
	case "30303030363030303030383030303131":
		result = `"1"`
	case "30303030363030303030383830303131":
		result = `"1"`
	case "30303030363030303030383030303132":
		result = `"2"`
	case "30303030363030303030383230303132":
		result = `"2"`
	case "30303030363030303030383830303132":
		result = `"2"`

	default:
		errMsg := function + " - gw45gd unknown input response: " + rawReply
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("unknown response error")
	}

	return result, nil
}

// HexToBin is from https://www.php2golang.com/method/function.hex2bin.html
func HexToBin(hex string) string {
	var retStr string = ""
	tempStr := hex

	for len(tempStr) > 0 {
		chunk := tempStr[:2]
		b, _ := strconv.ParseUint(chunk, 16, 8)
		retStr = retStr + string(rune(b))
		tempStr = tempStr[2:]
	}

	tempStr = retStr
	for len(tempStr) > 0 {
		chunk := tempStr[:1]
		i := int(chunk[0])
		fmt.Printf("%02x", i)
		tempStr = tempStr[1:]
	}
	return retStr
}

func ExtractBody(message string) string {
	// Remove the header, checksum, etc. from the message
	// (reduces variability in the responses to make them easier to match)

	if len(message) < 20 { // message is too short, processing would cause a panic
		return message
	}
	// Header is 14 bytes and then there are two bytes start of message delimiter
	body := message[16:]

	// Remove last bytes until we remove end of message delimiter (some devices don't emit a checksum)
	// Note: the framework sometimes removes the NEC device checksum bytes (e.g. 0x09 which is tab)
	// This function accounts for that by simply removing until it reaches 0x03.  We don't verify checksums
	// which would be redundant on TCP/IP.
	for {
		bytesPair := body[len(body)-2:]
		body = body[:len(body)-2]
		if bytesPair == "03" {
			break
		}
		if len(body) < 2 {
			break
		}
	}

	return body
}

func BinToHex(str string) string {
	retStr := ""
	tempStr := str
	var chunk string
	var chunkStr string

	// fmt.Println("Initial tempStr length:", len(tempStr))
	for len(tempStr) > 0 {
		// fmt.Println("tempStr before: " + tempStr + " ")
		// fmt.Println("tempStr length:", len(tempStr))
		chunk = tempStr[:1]
		i := int(chunk[0])
		chunkStr = fmt.Sprintf("%02x", i)
		// fmt.Println("chunk: " + chunkStr)
		retStr = retStr + chunkStr
		tempStr = tempStr[1:]
		// fmt.Println("tempStr after: " + tempStr + " ")
	}

	return retStr
}

// volumeSetValueMapper converts an integer volume value to the corresponding NEC command
func volumeSetValueMapper(volumeInt int) (string, error) {
	var hexCommand string

	switch volumeInt {
	case 0:
		hexCommand = "0130413045304102303036323030303003700d"
	case 1:
		hexCommand = "0130413045304102303036323030303103710d"
	case 2:
		hexCommand = "0130413045304102303036323030303203720d"
	case 3:
		hexCommand = "0130413045304102303036323030303303730d"
	case 4:
		hexCommand = "0130413045304102303036323030303403740d"
	case 5:
		hexCommand = "0130413045304102303036323030303503750d"
	case 6:
		hexCommand = "0130413045304102303036323030303603760d"
	case 7:
		hexCommand = "0130413045304102303036323030303703770d"
	case 8:
		hexCommand = "0130413045304102303036323030303803780d"
	case 9:
		hexCommand = "0130413045304102303036323030303903790d"
	case 10:
		hexCommand = "0130413045304102303036323030304103010d"
	case 11:
		hexCommand = "0130413045304102303036323030304203020d"
	case 12:
		hexCommand = "0130413045304102303036323030304303030d"
	case 13:
		hexCommand = "0130413045304102303036323030304403040d"
	case 14:
		hexCommand = "0130413045304102303036323030304503050d"
	case 15:
		hexCommand = "0130413045304102303036323030304603060d"
	case 16:
		hexCommand = "0130413045304102303036323030313003710d"
	case 17:
		hexCommand = "0130413045304102303036323030313103700d"
	case 18:
		hexCommand = "0130413045304102303036323030313203730d"
	case 19:
		hexCommand = "0130413045304102303036323030313303720d"
	case 20:
		hexCommand = "0130413045304102303036323030313403750d"
	case 21:
		hexCommand = "0130413045304102303036323030313503740d"
	case 22:
		hexCommand = "0130413045304102303036323030313603770d"
	case 23:
		hexCommand = "0130413045304102303036323030313703760d"
	case 24:
		hexCommand = "0130413045304102303036323030313803790d"
	case 25:
		hexCommand = "0130413045304102303036323030313903780d"
	case 26:
		hexCommand = "0130413045304102303036323030314103000d"
	case 27:
		hexCommand = "0130413045304102303036323030314203030d"
	case 28:
		hexCommand = "0130413045304102303036323030314303020d"
	case 29:
		hexCommand = "0130413045304102303036323030314403050d"
	case 30:
		hexCommand = "0130413045304102303036323030314503040d"
	case 31:
		hexCommand = "0130413045304102303036323030314603070d"
	case 32:
		hexCommand = "0130413045304102303036323030323003720d"
	case 33:
		hexCommand = "0130413045304102303036323030323103730d"
	case 34:
		hexCommand = "0130413045304102303036323030323203700d"
	case 35:
		hexCommand = "0130413045304102303036323030323303710d"
	case 36:
		hexCommand = "0130413045304102303036323030323403760d"
	case 37:
		hexCommand = "0130413045304102303036323030323503770d"
	case 38:
		hexCommand = "0130413045304102303036323030323603740d"
	case 39:
		hexCommand = "0130413045304102303036323030323703750d"
	case 40:
		hexCommand = "01304130453041023030363230303238037a0d"
	case 41:
		hexCommand = "01304130453041023030363230303239037b0d"
	case 42:
		hexCommand = "0130413045304102303036323030324103030d"
	case 43:
		hexCommand = "0130413045304102303036323030324203000d"
	case 44:
		hexCommand = "0130413045304102303036323030324303010d"
	case 45:
		hexCommand = "0130413045304102303036323030324403060d"
	case 46:
		hexCommand = "0130413045304102303036323030324503070d"
	case 47:
		hexCommand = "0130413045304102303036323030324603040d"
	case 48:
		hexCommand = "0130413045304102303036323030333003730d"
	case 49:
		hexCommand = "0130413045304102303036323030333103720d"
	case 50:
		hexCommand = "0130413045304102303036323030333203710d"
	case 51:
		hexCommand = "0130413045304102303036323030333303700d"
	case 52:
		hexCommand = "0130413045304102303036323030333403770d"
	case 53:
		hexCommand = "0130413045304102303036323030333503760d"
	case 54:
		hexCommand = "0130413045304102303036323030333603750d"
	case 55:
		hexCommand = "0130413045304102303036323030333703740d"
	case 56:
		hexCommand = "01304130453041023030363230303338037b0d"
	case 57:
		hexCommand = "01304130453041023030363230303339037a0d"
	case 58:
		hexCommand = "0130413045304102303036323030334103020d"
	case 59:
		hexCommand = "0130413045304102303036323030334203010d"
	case 60:
		hexCommand = "0130413045304102303036323030334303000d"
	case 61:
		hexCommand = "0130413045304102303036323030334403070d"
	case 62:
		hexCommand = "0130413045304102303036323030334503060d"
	case 63:
		hexCommand = "0130413045304102303036323030334603050d"
	case 64:
		hexCommand = "0130413045304102303036323030343003740d"
	case 65:
		hexCommand = "0130413045304102303036323030343103750d"
	case 66:
		hexCommand = "0130413045304102303036323030343203760d"
	case 67:
		hexCommand = "0130413045304102303036323030343303770d"
	case 68:
		hexCommand = "0130413045304102303036323030343403700d"
	case 69:
		hexCommand = "0130413045304102303036323030343503710d"
	case 70:
		hexCommand = "0130413045304102303036323030343603720d"
	case 71:
		hexCommand = "0130413045304102303036323030343703730d"
	case 72:
		hexCommand = "01304130453041023030363230303438037c0d"
	case 73:
		hexCommand = "01304130453041023030363230303439037d0d"
	case 74:
		hexCommand = "0130413045304102303036323030344103050d"
	case 75:
		hexCommand = "0130413045304102303036323030344203060d"
	case 76:
		hexCommand = "0130413045304102303036323030344303070d"
	case 77:
		hexCommand = "0130413045304102303036323030344403000d"
	case 78:
		hexCommand = "0130413045304102303036323030344503010d"
	case 79:
		hexCommand = "0130413045304102303036323030344603020d"
	case 80:
		hexCommand = "0130413045304102303036323030353003750d"
	case 81:
		hexCommand = "0130413045304102303036323030353103740d"
	case 82:
		hexCommand = "0130413045304102303036323030353203770d"
	case 83:
		hexCommand = "0130413045304102303036323030353303760d"
	case 84:
		hexCommand = "0130413045304102303036323030353403710d"
	case 85:
		hexCommand = "0130413045304102303036323030353503700d"
	case 86:
		hexCommand = "0130413045304102303036323030353603730d"
	case 87:
		hexCommand = "0130413045304102303036323030353703720d"
	case 88:
		hexCommand = "01304130453041023030363230303538037d0d"
	case 89:
		hexCommand = "01304130453041023030363230303539037c0d"
	case 90:
		hexCommand = "0130413045304102303036323030354103040d"
	case 91:
		hexCommand = "0130413045304102303036323030354203070d"
	case 92:
		hexCommand = "0130413045304102303036323030354303060d"
	case 93:
		hexCommand = "0130413045304102303036323030354403010d"
	case 94:
		hexCommand = "0130413045304102303036323030354503000d"
	case 95:
		hexCommand = "0130413045304102303036323030354603030d"
	case 96:
		hexCommand = "0130413045304102303036323030363003760d"
	case 97:
		hexCommand = "0130413045304102303036323030363103770d"
	case 98:
		hexCommand = "0130413045304102303036323030363203740d"
	case 99:
		hexCommand = "0130413045304102303036323030363303750d"
	case 100:
		hexCommand = "0130413045304102303036323030363403720d"
	default:
		return "", errors.New("invalid volume level: must be between 0 and 100")
	}

	// Convert hex string to binary command
	command := HexToBin(hexCommand)
	return command, nil
}

func setVolume(socketKey string, output string, volume string) (string, error) {
	function := "setVolume"

	// The output parameter is ignored because there is only one output on the device
	framework.Log(function + " - setting volume for: " + socketKey + " to: " + volume)
	if output != "1" {
		errMsg := function + " - aseegdr invalid output number value: " + volume + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	// Parse the volume value from the string (removing any quotes)
	volumeStr := strings.Trim(volume, `"`)
	volumeInt, err := strconv.Atoi(volumeStr)
	if err != nil {
		errMsg := function + " - berabbd invalid volume value: " + volume
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid volume error")
	}

	// Get the appropriate command for this volume level
	command, err := volumeSetValueMapper(volumeInt)
	if err != nil {
		errMsg := function + " - asdfewasdb " + err.Error()
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid volume value error")
	}

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - brtnd error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - drbtbsfd still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("write error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + rawReply)

	// The volume response format would need to be verified against actual device responses
	// For now, we'll just assume success if we got any response
	if len(rawReply) > 0 {
		return volumeStr, nil
	} else {
		errMsg := function + " - ewafdgs no response from device"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("read error")
	}
}

// volumeResponseMapper maps a raw volume response to an integer volume value
func volumeResponseMapper(reply string) (int, error) {
	var volume int

	switch reply {
	case "30303030363230303030363430303030":
		volume = 0
	case "30303030363230303030363430303031":
		volume = 1
	case "30303030363230303030363430303032":
		volume = 2
	case "30303030363230303030363430303033":
		volume = 3
	case "30303030363230303030363430303034":
		volume = 4
	case "30303030363230303030363430303035":
		volume = 5
	case "30303030363230303030363430303036":
		volume = 6
	case "30303030363230303030363430303037":
		volume = 7
	case "30303030363230303030363430303038":
		volume = 8
	case "30303030363230303030363430303039":
		volume = 9
	case "30303030363230303030363430303041":
		volume = 10
	case "30303030363230303030363430303042":
		volume = 11
	case "30303030363230303030363430303043":
		volume = 12
	case "30303030363230303030363430303044":
		volume = 13
	case "30303030363230303030363430303045":
		volume = 14
	case "30303030363230303030363430303046":
		volume = 15
	case "30303030363230303030363430303130":
		volume = 16
	case "30303030363230303030363430303131":
		volume = 17
	case "30303030363230303030363430303132":
		volume = 18
	case "30303030363230303030363430303133":
		volume = 19
	case "30303030363230303030363430303134":
		volume = 20
	case "30303030363230303030363430303135":
		volume = 21
	case "30303030363230303030363430303136":
		volume = 22
	case "30303030363230303030363430303137":
		volume = 23
	case "30303030363230303030363430303138":
		volume = 24
	case "30303030363230303030363430303139":
		volume = 25
	case "30303030363230303030363430303141":
		volume = 26
	case "30303030363230303030363430303142":
		volume = 27
	case "30303030363230303030363430303143":
		volume = 28
	case "30303030363230303030363430303144":
		volume = 29
	case "30303030363230303030363430303145":
		volume = 30
	case "30303030363230303030363430303146":
		volume = 31
	case "30303030363230303030363430303230":
		volume = 32
	case "30303030363230303030363430303231":
		volume = 33
	case "30303030363230303030363430303232":
		volume = 34
	case "30303030363230303030363430303233":
		volume = 35
	case "30303030363230303030363430303234":
		volume = 36
	case "30303030363230303030363430303235":
		volume = 37
	case "30303030363230303030363430303236":
		volume = 38
	case "30303030363230303030363430303237":
		volume = 39
	case "30303030363230303030363430303238":
		volume = 40
	case "30303030363230303030363430303239":
		volume = 41
	case "30303030363230303030363430303241":
		volume = 42
	case "30303030363230303030363430303242":
		volume = 43
	case "30303030363230303030363430303243":
		volume = 44
	case "30303030363230303030363430303244":
		volume = 45
	case "30303030363230303030363430303245":
		volume = 46
	case "30303030363230303030363430303246":
		volume = 47
	case "30303030363230303030363430303330":
		volume = 48
	case "30303030363230303030363430303331":
		volume = 49
	case "30303030363230303030363430303332":
		volume = 50
	case "30303030363230303030363430303333":
		volume = 51
	case "30303030363230303030363430303334":
		volume = 52
	case "30303030363230303030363430303335":
		volume = 53
	case "30303030363230303030363430303336":
		volume = 54
	case "30303030363230303030363430303337":
		volume = 55
	case "30303030363230303030363430303338":
		volume = 56
	case "30303030363230303030363430303339":
		volume = 57
	case "30303030363230303030363430303341":
		volume = 58
	case "30303030363230303030363430303342":
		volume = 59
	case "30303030363230303030363430303343":
		volume = 60
	case "30303030363230303030363430303344":
		volume = 61
	case "30303030363230303030363430303345":
		volume = 62
	case "30303030363230303030363430303346":
		volume = 63
	case "30303030363230303030363430303430":
		volume = 64
	case "30303030363230303030363430303431":
		volume = 65
	case "30303030363230303030363430303432":
		volume = 66
	case "30303030363230303030363430303433":
		volume = 67
	case "30303030363230303030363430303434":
		volume = 68
	case "30303030363230303030363430303435":
		volume = 69
	case "30303030363230303030363430303436":
		volume = 70
	case "30303030363230303030363430303437":
		volume = 71
	case "30303030363230303030363430303438":
		volume = 72
	case "30303030363230303030363430303439":
		volume = 73
	case "30303030363230303030363430303441":
		volume = 74
	case "30303030363230303030363430303442":
		volume = 75
	case "30303030363230303030363430303443":
		volume = 76
	case "30303030363230303030363430303444":
		volume = 77
	case "30303030363230303030363430303445":
		volume = 78
	case "30303030363230303030363430303446":
		volume = 79
	case "30303030363230303030363430303530":
		volume = 80
	case "30303030363230303030363430303531":
		volume = 81
	case "30303030363230303030363430303532":
		volume = 82
	case "30303030363230303030363430303533":
		volume = 83
	case "30303030363230303030363430303534":
		volume = 84
	case "30303030363230303030363430303535":
		volume = 85
	case "30303030363230303030363430303536":
		volume = 86
	case "30303030363230303030363430303537":
		volume = 87
	case "30303030363230303030363430303538":
		volume = 88
	case "30303030363230303030363430303539":
		volume = 89
	case "30303030363230303030363430303541":
		volume = 90
	case "30303030363230303030363430303542":
		volume = 91
	case "30303030363230303030363430303543":
		volume = 92
	case "30303030363230303030363430303544":
		volume = 93
	case "30303030363230303030363430303545":
		volume = 94
	case "30303030363230303030363430303546":
		volume = 95
	case "30303030363230303030363430303630":
		volume = 96
	case "30303030363230303030363430303631":
		volume = 97
	case "30303030363230303030363430303632":
		volume = 98
	case "30303030363230303030363430303633":
		volume = 99
	case "30303030363230303030363430303634":
		volume = 100
	default:
		return 0, errors.New("unknown volume response")
	}

	return volume, nil
}

func getVolume(socketKey string, outputNumber string) (string, error) {
	function := "getVolume"
	command := HexToBin("01304130433036023030363203010d")

	// output parameter is ignored in this device
	if outputNumber != "1" {
		errMsg := function + " - brearj6 invalid output number value: " + outputNumber + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	framework.Log(function + " - getting volume for: " + socketKey)

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - aswfefawdf error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - vacdca still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("write error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + rawReply)

	// convert the raw response to a volume integer
	volume, err := volumeResponseMapper(rawReply)
	if err != nil {
		errMsg := function + " - waewafdv unknown volume response: " + rawReply
		framework.AddToErrors(socketKey, errMsg)
		return "", errors.New("unknown response error")
	}

	// Return the volume as a string
	volumeStr := strconv.Itoa(volume)
	return `"` + volumeStr + `"`, nil
}

func getAudioMute(socketKey string, outputNumber string) (string, error) {
	function := "getAudioMute"

	// output parameter is ignored in this device
	if outputNumber != "1" {
		errMsg := function + " - vaewaeg invalid output number value: " + outputNumber + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	framework.Log(function + " - getting audio mute status for: " + socketKey)

	command := HexToBin("01304130433036023030384403790d")

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - awefavd error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - vfvawe still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("get mute error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + rawReply)

	var result string

	// Check for unmuted responses
	if rawReply == "30303030384430303030303230303032" ||
		rawReply == "30303030384430303030303230303030" {
		result = `"false"`
		// Check for muted response
	} else if rawReply == "30303030384430303030303230303031" {
		result = `"true"`
	} else {
		errMsg := function + " - asfafee unknown mute response: " + rawReply
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("unknown response error")
	}

	return result, nil
}

func setAudioMute(socketKey string, outputNumber string, state string) (string, error) {
	function := "setAudioMute"
	var command string
	var expectedState string

	framework.Log(function + " - setting audio mute for: " + socketKey + " state: " + state)

	// output parameter is ignored in this device
	if outputNumber != "1" {
		errMsg := function + " - vafewaf invalid output number value: " + outputNumber + " (must be 1)"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("invalid output error")
	}

	// Remove quotes from state if present
	state = strings.Trim(state, `"`)

	// Determine command based on state
	switch state {
	case "true":
		command = HexToBin("0130413045304102303038443030303103090d")
		expectedState = "true"
	case "false":
		command = HexToBin("01304130453041023030384430303032030a0d")
		expectedState = "false"
	case "toggle":
		// Get current mute state
		currentState, err := getAudioMute(socketKey, "1")
		if err != nil {
			return "", err
		}

		// Toggle the state
		if currentState == `"false"` {
			command = HexToBin("0130413045304102303038443030303103090d")
			expectedState = "true"
		} else if currentState == `"true"` {
			command = HexToBin("01304130453041023030384430303032030a0d")
			expectedState = "false"
		} else {
			errMsg := function + " - bterwbfs unknown current mute state: " + currentState
			framework.AddToErrors(socketKey, errMsg)
			return errMsg, errors.New("unknown state error")
		}
	default:
		errMsg := function + " - kmvasdv unknown mute command: " + state
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("unknown command error")
	}

	framework.Log(function + " - setting audio mute for: " + socketKey + " to: " + expectedState)

	if !framework.WriteLineToSocket(socketKey, command) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - vafewaf error writing to " + socketKey + " closing and trying again"
		framework.AddToErrors(socketKey, errMsg)
		framework.CloseSocketConnection(socketKey)

		if !framework.WriteLineToSocket(socketKey, command) {
			errMsg := function + " - vafewaf still getting an error writing to " + socketKey + " giving up"
			framework.AddToErrors(socketKey, errMsg)
			framework.CloseSocketConnection(socketKey)
			return errMsg, errors.New("write error")
		}
	}

	// get server response
	response := framework.ReadLineFromSocket(socketKey)
	rawReply := ExtractBody(BinToHex(response))

	framework.Log(function + " - raw hex response is: " + rawReply)

	var result string

	// Check for unmuted response
	if rawReply == "30303030384430303030303230303032" ||
		rawReply == "30303030384430313030303230303032" {
		result = `"false"`
		// Check for muted response
	} else if rawReply == "30303030384430303030303230303031" ||
		rawReply == "30303030384430313030303230303031" {
		result = `"true"`
	} else {
		errMsg := function + " - vafewaf unknown mute response: " + rawReply
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New("unknown response error")
	}

	return result, nil
}
