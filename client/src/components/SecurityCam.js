import React, { useState, useEffect, useRef } from 'react';
import Webcam from 'react-webcam';

function SecurityCam(props) {

  const { addIncident } = props;
  const [movement, setMovement] = useState('OK')
  const [videoConstraints, setVideoConstraints] = useState('user') // user-facing/selfie
  let pre, post;
  const threshold = 15;
  const interval = 500;
  let diffImg = '';

  const webcamRef = useRef(null);
  const capture = async () => {
    if (webcamRef.current) {
      const pic = await webcamRef.current.getScreenshot();
      if (typeof pic === 'string') {

        pre = post ? post : pic;
        post = pic;

        compare(pre, post, function (result) {
          if (result > threshold) {
            console.log("MOTION DETECTED");
            addIncident(post);
          }
        });
      }
    }
  }

  useEffect(() => {
    const captureInterval = setInterval(capture, interval)
    return () => clearInterval(captureInterval)
  }, [])

  return (
    <div className='SecurityCam'>
      <h3>Monitoring</h3>
      <h4>{movement}</h4>
      {/* <button
        onClick={() => setVideoConstraints(videoConstraints === 'user' ? { exact: 'environment' } : 'user')}
      >
        {videoConstraints === 'user' ? 'Use Forward Camera' : 'Use Selfie Camera'}
      </button> */}
      <Webcam
        audio={false}
        screenshotFormat='image/jpeg'
        videoConstraints={videoConstraints}
        ref={webcamRef}
      />
    </div>
  )


  // The following function taken from
  // https://rosettacode.org/wiki/Percentage_difference_between_images
  function getImageData(url, callback) {
    const img = document.createElement('img');
    const canvas = document.createElement('canvas');
    
    img.onload = function () {
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext('2d');
      ctx.drawImage(img, 0, 0);
      callback(ctx.getImageData(0, 0, img.width, img.height));
    };
    
    img.src = url;
  }
  
  // The following function taken from
  // https://rosettacode.org/wiki/Percentage_difference_between_images
  function compare(firstImage, secondImage, callback) {
    getImageData(firstImage, function (img1) {
      getImageData(secondImage, function (img2) {
        if (img1.width !== img2.width || img1.height != img2.height) {
          callback(NaN);
          return;
        }

        let diff = 0;

        for (var i = 0; i < img1.data.length / 4; i++) {
          diff += Math.abs(img1.data[4 * i + 0] - img2.data[4 * i + 0]) / 255;
          diff += Math.abs(img1.data[4 * i + 1] - img2.data[4 * i + 1]) / 255;
          diff += Math.abs(img1.data[4 * i + 2] - img2.data[4 * i + 2]) / 255;
        }
        diffImg = diff
        callback(100 * diff / (img1.width * img1.height * 3));
      });
    });
  }

}

export default SecurityCam;
