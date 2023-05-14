// Establish a WebSocket connection
const socket = new WebSocket("ws://localhost:3000/ws");

// Event handler for WebSocket connection open
socket.onopen = () => {
  console.log("WebSocket connection established");
};

// Event handler for WebSocket connection close
socket.onclose = () => {
  console.log("WebSocket connection closed");
};

// Event handler for receiving WebSocket messages
socket.onmessage = (event) => {
    console.log("returned message")
  const message = event.data;
  console.log("Received message from server:", message);
  
  // Handle the received message as needed
  // For example, parse the message as JSON
  // const parsedMessage = JSON.parse(message);
  // console.log("Parsed message:", parsedMessage);

  // Perform specific actions based on the received message
  if (parsedMessage.type === "notification") {
    displayNotification(parsedMessage.content);
  } else if (parsedMessage.type === "data") {
    processData(parsedMessage.data);
  }
};

// Function to send a message via WebSocket
function sendMessage(message) {
  socket.send(message);
}

// Example usage
sendMessage("Hello, server!");
