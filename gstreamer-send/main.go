package main

import (
	"encoding/json"
	"fmt"

	"github.com/pion/webrtc/v3"

	"github.com/pion/example-webrtc-applications/v3/internal/signal"
)

func main() {
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

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	a, err := json.Marshal(offer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(a))

	answer := webrtc.SessionDescription{}

	err = json.Unmarshal([]byte(signal.MustReadStdin()), &answer)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	<-gatherComplete

	// // Start pushing buffers on these tracks
	// gst.CreatePipeline("opus", []*webrtc.TrackLocalStaticSample{audioTrack}, *audioSrc).Start()
	// gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{firstVideoTrack, secondVideoTrack}, *videoSrc).Start()
	dataChannel.Send([]byte("jMeAss"))
	// Block forever
	select {}
}
