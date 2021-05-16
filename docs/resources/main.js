function getGraph() {
  var endpoint = 'http://localhost:3000/api/contributions?';
  var markdown = '![gh-contributions](endpoint)';

  var usernames = document.getElementById("usernames").value;
  var year = document.getElementById("year").value;
  var theme = document.getElementById("theme").value;

  if (usernames == "") {
    alert("Username is required");
    return false;
  }

  if (year == "") {
    var date = new Date();
    var currentYear = date.getFullYear();
    year = currentYear;
  }

  usernames = getUsernames();

  // Add usernames
  for (var i = 0; i < usernames.length; i++) {
    if (i == 0) {
      endpoint = endpoint + 'username=' + usernames[i]
    } else {
      endpoint = endpoint + '&username=' + usernames[i]
    }
  }

  // Add year
  endpoint += '&year=' + year

  // Add theme
  endpoint += '&theme=' + theme
  console.log(endpoint);

  markdownText = markdown.replace('endpoint', endpoint);

  document.getElementById("markdown-text").innerHTML = markdownText;
  document.getElementById("markdown-text").style.display = 'block';

  document.getElementById("graph").src = endpoint;
  document.getElementById("graph").style.display = 'block';
}

function getUsernames() {
  var username = document.getElementById("usernames").value;
  var usernamesList = username.replace(/\s/g, '').split(",");

  return usernamesList
}
