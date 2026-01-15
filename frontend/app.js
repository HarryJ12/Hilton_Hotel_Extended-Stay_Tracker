async function loadGuests() {
  // const res = await fetch("http://localhost:8080/guests"); // for local testing
  const res = await fetch("/guests");
  const guests = await res.json();

  const list = document.getElementById("guestList");
  list.innerHTML = "";

  guests.forEach(g => {
    const li = document.createElement("li");

  li.innerHTML = `
    <strong>${g.name}</strong> <br>
    <span>Room Number:</span> ${g.room_number}<br>
    <span>Check-In Date:</span> ${g.check_in_date.split("T")[0]}<br>
    <span>Contact Information:</span> ${g.contact}<br>
  `;

  const delBtn = document.createElement("button");
  delBtn.textContent = "Delete";
  delBtn.style.marginTop = "3px";  

    delBtn.onclick = async () => {
      // const res = await fetch(`http://localhost:8080/guests/${g.id}`, { // for local testing
      const res = await fetch(`/guests/${g.id}`, {
        method: "DELETE"
      });

      if (!res.ok) {
        alert("Delete failed");
        return;
      }

      loadGuests();
    };

    li.appendChild(delBtn);
    list.appendChild(li);
  });
}

document.getElementById("guestForm").addEventListener("submit", async e => {
  e.preventDefault();

  // await fetch("http://localhost:8080/guests", { // for local testing
  await fetch("/guests", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: document.getElementById("name").value,
      contact: document.getElementById("contact").value,
      room_number: document.getElementById("room").value,
      daily_rate: parseInt(document.getElementById("rate").value, 10),
      check_in_date: document.getElementById("checkin").value
    })
  });

  e.target.reset();
  loadGuests();
});

loadGuests();
