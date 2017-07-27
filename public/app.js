new Vue({
  el: '#app',

  data: {
    ws: null, // Our websocket
    newMsg: '', // Holds new messages to be sent to the server
    chatContent: '', // A running list of chat messages displayed on the screen
    email: null, // Email address used for grabbing an avatar
    username: null, // Our username
    destination: '', // Actual user who received messages
    joined: false // True if email and username have been filled in
  },
  created: function() {
    var self = this;
    this.ws = new WebSocket('wss://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      if (msg.type == "send") {
        self.chatContent += '<div class="chip">'
        + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
        + msg.username
        + '</div>'
        + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
      }
      else if (msg.type == "error") {
        Materialize.toast(msg.message, 2000);
      }

      var element = document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
    });
  },
  methods: {
    send: function () {
      if (this.newMsg != '') {
        this.ws.send(
          JSON.stringify({
            type: "send",
            email: this.email,
            username: this.username,
            destination: this.destination,
            message: $('<p>').html(this.newMsg).text() // Strip out html
          }
        ));
        this.newMsg = ''; // Reset newMsg
      }
      else {
        Materialize.toast('You message is empty', 2000);
        return
      }
    },
    join: function () {
      if (!this.email) {
        Materialize.toast('You must enter an email', 2000);
        return
      }
      if (!this.username) {
        Materialize.toast('You must choose a username', 2000);
        return
      }
      this.email = $('<p>').html(this.email).text();
      this.username = $('<p>').html(this.username).text();
      this.ws.send(
        JSON.stringify({
          type: "connect",
          email: this.email,
          username: this.username,
          destination: this.destination,
          message: ""
        }
      ));
      this.joined = true;
    },
    gravatarURL: function(email) {
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
  }
});
