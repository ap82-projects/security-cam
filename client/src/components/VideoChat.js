import React, { useState, useCallback, useEffect } from "react";
import Video from "twilio-video";
import Lobby from "./Lobby";
import Room from "./Room";
import axios from "axios";

const VideoChat = (props) => {
  // const { guestName, guestRoom, setRoomID } = props;
  const { guestName, guestRoom, serverURL } = props;
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

      // Used when generating any kind of tokens
      // To set up environmental variables, see http://twil.io/secure
      const twilioAccountSid = process.env.REACT_APP_TWILIO_ACCOUNT_SID;
      const twilioApiKey = process.env.REACT_APP_TWILIO_API_KEY_SID;
      const twilioApiSecret = process.env.REACT_APP_TWILIO_API_KEY_SECRET;

      const identity = username;

      // Create Video Grant
      const videoGrant = new VideoGrant({
        room: roomName,
      });

      // Create an access token which we will sign and return to the client,
      // containing the grant we just created
      const token = new AccessToken(
        twilioAccountSid,
        twilioApiKey,
        twilioApiSecret,
        { identity: identity }
      );
      token.addGrant(videoGrant);

      // Serialize the token to a JWT string
      // console.log(token.toJwt());

      // Video.connect(data.token, {
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
      // setRoomID();
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
