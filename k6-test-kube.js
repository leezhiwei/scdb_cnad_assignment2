import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: 'shared-iterations',
      vus: 1000,
      iterations: 25000,
      maxDuration: '5m',
    },
  },
};

// The default exported function is gonna be picked up by k6 as the entry point for the test script. It will be executed repeatedly in "iterations" for the whole duration of the test.
export default function () {
  // Make a GET request to the target URL
  http.get('https://scdb.leezhiwei.com');

  // Sleep for 1 second to simulate real-world usage
  sleep(1);
}