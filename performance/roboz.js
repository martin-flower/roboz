import http from 'k6/http';
import { check, sleep } from 'k6';

// default required - entry point - similar to main
// contains VU code - is run over and over again
export default function () {

  // health
  let response = http.get('http://localhost:5000/health');
  check(response, {
    'health status 200': r => r.status == 200,
    'health response time - good enough': r => r.timings.duration < 100, // 100 milliseconds
  });
  sleep(0.4); // seconds

  // documentation
  response = http.get('http://localhost:5000');
  check(response, {
    'documentation status 200': r => r.status == 200,
    'documentation response time - good enough': r => r.timings.duration < 200, // 200 milliseconds
  });
  sleep(0.6); // seconds

  // enter
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };  

  response = http.post(`http://localhost:5000/developer-test/enter-path`, randomEnterPayload(), params);
  check(response, {
    'enter status 200': r => r.status == 200,
    'enter response time - good enough': r => r.timings.duration < 800, // 800 milliseconds
  });
  sleep(0.5);
}

function randomEnterPayload() {
  return JSON.stringify(
    {
      start: {
        x: 10,
        y: 22
      },
      commands: [
        {
          direction: "west",
          steps: 4
        },
        {
          direction: "south",
          steps: 4
        },
        {
          direction: "east",
          steps: 3
        }
      ]
    }
  );
}

// instead of adding the parameters on the command line
export let options = {
  vus: 1,
  duration: '5s',
  //    stages: [
  //        { duration: '30s', target: 20 },
  //        { duration: '1m30s', target: 10 },
  //        { duration: '20s', target: 0 },
  //    ],
};
