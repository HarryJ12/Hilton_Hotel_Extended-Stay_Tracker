// Load all guests from the server and show them on the page
async function loadGuests() {

  // Ask the backend for the current list of guests (GET)
  // const res = await fetch("http://localhost:8080/guests"); // for local testing
  const res = await fetch("/guests");

  // Convert response into a JavaScript array of guest objects
  const guests = await res.json();

  // Find the guest list element and clear it
  const list = document.getElementById("guestList");
  list.innerHTML = "";

  // For each guest, create a list item with guest details and add it to the list
  guests.forEach(g => {
    const li = document.createElement("li");
  li.innerHTML = `
    <strong>${g.name}</strong> <br>
    <span>Room Number:</span> ${g.room_number}<br>
    <span>Check-In Date:</span> ${g.check_in_date.split("T")[0]}<br>
    <span>Contact Information:</span> ${g.contact}<br> 
    `;

  // Creates a Delete button for the guest
  const delBtn = document.createElement("button");
  delBtn.textContent = "Delete";
  delBtn.style.marginTop = "3px";  

    // When Delete is clicked:
    // 1. Tell the backend to remove this guest by ID (DELETE)
    // 2. Reload the list so the UI updates
    delBtn.onclick = async () => {
      // const res = await fetch(`http://localhost:8080/guests/${g.id}`, { // for local testing
      const res = await fetch(`/guests/${g.id}`, {
        method: "DELETE"
      });

      // If delete failed, show an alert and do not reload the list
      if (!res.ok) {
        alert("Delete failed");
        return;
      }

      // Reload the guest list to reflect the deletion
      loadGuests();
    };

    // Append the Delete button to the list item and the list item to the list
    li.appendChild(delBtn);
    list.appendChild(li);
  });
}

// Handle the form submission to add a new guest to stop the browser from reloading the page
document.getElementById("guestForm").addEventListener("submit", async e => {
  e.preventDefault();

  // Sends the form data to the backend to create a new guest (POST)
  // await fetch("http://localhost:8080/guests", { // for local testing
  await fetch("/guests", {
    method: "POST", 
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      // Read values from the input fields
      name: document.getElementById("name").value,
      contact: document.getElementById("contact").value,
      room_number: document.getElementById("room").value,
      daily_rate: parseInt(document.getElementById("rate").value, 10), // Convert rate to a number
      check_in_date: document.getElementById("checkin").value
    })
  });

  // Clear the form after submission
  e.target.reset();

  // Reload guests so the new one appears
  loadGuests();
});

// Initial load of guests when the page is opened
loadGuests();
