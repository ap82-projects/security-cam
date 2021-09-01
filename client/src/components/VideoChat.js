import React, { useState, useCallback, useEffect } from "react";
import Video from "twilio-video";
import Lobby from "./Lobby";
import Room from "./Room";

const VideoChat = (props) => {
  const { guestName, guestRoom, axios } = props;
  const [username, setUsername] = useState(guestName ? guestName : "");
  const [roomName, setRoomName] = useState(guestRoom ? guestRoom : "");
  const [room, setRoom] = useState(null);
  const [connecting, setConnecting] = useState(false);

  const handleUsernameChange = useCallback((event) => {
    setUsername(event.target.value);
  }, []);

  const handleRoomNameChange = useCallback((event) => {
    setRoomName(event.target.value);
  }, []);

  const handleSubmit = useCallback(
    async (event) => {
      event.preventDefault();
      setConnecting(true);
      ///////////////////////////
      // For getting token from server
      //
      // const response = await axios.get(`/api/videotoken/${roomName}/${username}`);
      // const { data } = response;
      ///////////////////////////

      const AccessToken = require('twilio').jwt.AccessToken;
      const VideoGrant = AccessToken.VideoGrant;

      const identity = username;

      // Create Video Grant
      const videoGrant = new VideoGrant({
        room: roomName,
      });

      // Create an access token
      const twilioData = (await axios.get("/api/twiliodata")).data
      const token = new AccessToken(
        twilioData.accountSid,
        twilioData.apiKeySid,
        twilioData.apiSecretKey,
        { identity: identity }
      );
      token.addGrant(videoGrant);

      Video.connect(token.toJwt(), {
        name: roomName,
        audio: true,
        video: true
      })
        .then((room) => {
          setConnecting(false);
          setRoom(room);
        })
        .catch((err) => {
          console.error(err);
          setConnecting(false);
        });
    },
    [roomName, username]
  );

  const handleLogout = useCallback(() => {
    setRoom((prevRoom) => {
      if (prevRoom) {
        prevRoom.localParticipant.tracks.forEach((trackPub) => {
          trackPub.track.stop();
        });
        prevRoom.disconnect();
      }
      return null;
    });
  }, []);

  useEffect(() => {
    if (room) {
      const tidyUp = (event) => {
        if (event.persisted) {
          return;
        }
        if (room) {
          handleLogout();
        }
      };
      window.addEventListener("pagehide", tidyUp);
      window.addEventListener("beforeunload", tidyUp);
      return () => {
        window.removeEventListener("pagehide", tidyUp);
        window.removeEventListener("beforeunload", tidyUp);
      };
    }
  }, [room, handleLogout]);

  let render;
  if (room) {
    render = (
      <Room roomName={roomName} room={room} handleLogout={handleLogout} guestRoom={guestRoom} />
    );
  } else {
    render = (
      <Lobby
        username={username}
        roomName={roomName}
        handleUsernameChange={handleUsernameChange}
        handleRoomNameChange={handleRoomNameChange}
        handleSubmit={handleSubmit}
        connecting={connecting}
      />
    );
  }
  return render;
};

export default VideoChat;
