async function signup(username, password) {
    const response = fetch("/api/signup", {
        method: "POST",
        body: JSON.stringify({ Username: username, Password: password })
    })

    const data = await response

    if (!data.ok) {
        console.log("Error")
        return
    }

    const realData = await data.json()
    console.log(realData)
}

document.addEventListener("DOMContentLoaded", async () => {
    //signup("itsbailey444", "TestPassword123")
})