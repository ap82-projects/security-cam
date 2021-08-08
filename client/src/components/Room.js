import React, { useEffect, useState } from "react";
// import { React } from "../react";
import Participant from "./Participant";

const Room = ({ roomName, room, handleLogout, guestRoom }) => {
  const [participants, setParticipants] = useState([]);
  // const guestURL = window.location.protocol 
  //                  + "//" + window.location.host
  //                  + "/?guestRoom=" + guestRoom;

  useEffect(() => {
    const participantConnected = (participant) => {
      setParticipants((prevParticipants) => [...prevParticipants, participant]);
    };

    const participantDisconnected = (participant) => {
      setParticipants((prevParticipants) =>
        prevParticipants.filter((p) => p !== participant)
      );
    };

    room.on("participantConnected", participantConnected);
    room.on("participantDisconnected", participantDisconnected);
    room.participants.forEach(participantConnected);
    return () => {
      room.off("participantConnected", participantConnected);
      room.off("participantDisconnected", participantDisconnected);
    };
  }, [room]);

  const remoteParticipants = participants.map((participant) => (
    <Participant key={participant.sid} participant={participant} />
  ));

  return (
    <div className="room">
      {/* <h4>Room: {roomName}</h4> */}
      
      {/* <h5>Remote Participants</h5> */}
      <div className="remote-participants">{remoteParticipants}</div>
      <div className="local-participant">
        {room ? (
          <Participant
            key={room.localParticipant.sid}
            participant={room.localParticipant}
          />
        ) : (
          ""
        )}
        <button onClick={handleLogout}>Leave Room</button>
      </div>
    </div>
  );
};

export default Room;
