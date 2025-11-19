function getCookie(name) {
	const cookies = document.cookie.split(';');
	for (let i = 0; i < cookies.length; i++) {
		const pair = cookies[i].trim().split('=');
		if (pair[0] === name) {
			return decodeURIComponent(pair[1] || '');
		}
	}
	return null;
}

async function ValidateSession() {
    const resp = await fetch("/api/logout", {
        method: "POST",
        body: JSON.stringify({
            "sessionID" : getCookie("session_id")
        })
    })

    if (!resp.ok) {
        console.log("Error fetching data")
        return
    }

    const data = await resp.json()

    console.log(data)
}

document.addEventListener("DOMContentLoaded", async () => {
    await ValidateSession()
})