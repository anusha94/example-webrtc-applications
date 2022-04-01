package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/pion/webrtc/v3"

	gst "github.com/pion/example-webrtc-applications/v3/internal/gstreamer-sink"
	"github.com/pion/example-webrtc-applications/v3/internal/signal"
)

// gstreamerReceiveMain is launched in a goroutine because the main thread is needed
// for Glib's main loop (Gstreamer uses Glib)
func gstreamerReceiveMain() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		fmt.Println("got ice candidate")

		if i != nil {
			peerConnection.AddICECandidate(i.ToJSON())
		}
	})

	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		// fmt.Println(dc.Label())
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Println(msg)
		})
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	offer := webrtc.SessionDescription{}

	err = json.Unmarshal([]byte(signal.MustReadStdin()), &offer)
	if err != nil {
		panic(err)
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	a, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(a))

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Block forever
	select {}
}

func init() {
	// This example uses Gstreamer's autovideosink element to display the received video
	// This element, along with some others, sometimes require that the process' main thread is used
	runtime.LockOSThread()
}

func main() {
	// Start a new thread to do the actual work for this application
	go gstreamerReceiveMain()
	// Use this goroutine (which has been runtime.LockOSThread'd to he the main thread) to run the Glib loop that Gstreamer requires
	gst.StartMainLoop()
}
