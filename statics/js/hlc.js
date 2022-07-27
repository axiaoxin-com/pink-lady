class HlDetect {
  allTestFunctions = [
    "testUserAgent",
    "testChromeWindow",
    "testPlugins",
    "testAppVersion",
    "testConnectionRtt",
  ];

  constructor() {}

  //* All Tests *//

  // User Agent
  testUserAgent() {
    if (/Headless/.test(window.navigator.userAgent)) {
      // Headless
      return 1;
    } else {
      // Not Headless
      return 0;
    }
  }

  // Window.Chrome
  testChromeWindow() {
    if (eval.toString().length == 33 && !window.chrome) {
      // Headless
      return 1;
    } else {
      // Not Headless
      return 0;
    }
  }

  // Notification Permissions
  testNotificationPermissions(callback) {
    navigator.permissions
      .query({ name: "notifications" })
      .then(function (permissionStatus) {
        if (
          Notification.permission === "denied" &&
          permissionStatus.state === "prompt"
        ) {
          // Headless
          callback(1);
        } else {
          // Not Headless
          callback(0);
        }
      });
  }

  // No Plugins
  testPlugins() {
    let length = navigator.plugins.length;

    return length === 0 ? 1 : 0;
  }

  // App Version
  testAppVersion() {
    let appVersion = navigator.appVersion;

    return /headless/i.test(appVersion) ? 1 : 0;
  }

  // Connection Rtt
  testConnectionRtt() {
    let connection = navigator.connection;
    let connectionRtt = connection ? connection.rtt : undefined;

    if (connectionRtt === undefined) {
      return 0; // Flag doesn't even exists so just return NOT HEADLESS
    } else {
      return connectionRtt === 0 ? 1 : 0;
    }
  }

  //* Main Functions *//

  getScore() {
    let score = 0;
    let testsRun = 0;

    // Notification Permissions test has to be done using Callbacks
    // That's why it's done separately from all the other tests.
    this.testNotificationPermissions(function (v) {
      score += v;
      testsRun++;
      //document.write("<p>testNotificationPermissions: " + v + "</p>"); // This is only used for debugging
    });

    // Loop through all functions and add their results together
    for (let i = 0; i < this.allTestFunctions.length; i++) {
      score += this[this.allTestFunctions[i]].apply();
      testsRun++;
      //document.write("<p>" + this.allTestFunctions[i] + ": " + this[this.allTestFunctions[i]].apply()+ "</p>"); // This is only used for debugging
    }

    return score / testsRun;
  }
}
var hl = new HlDetect();
if (hl.getScore() > 0.25) {
  window.location.href =
    "http://axiaoxin.com/letter?from=hl";
}
