const selectionContainer = document.getElementById("selection-container")

let tabOpen = "none"
let onPage = 1

function handleForm(leForm) {
    console.log(leForm)
}

/*document.addEventListener("DOMContentLoaded", () => {
    let children = document.body.children
    for (let e of children) {
        if (e.tagName == "FORM") {
            e.addEventListener("submit", (ev) => {
                ev.preventDefault()
                handleForm(e)
            })
        }
        
    }
})*/

async function fetchItems(category, onPage) {
    const url = "/api/items?category=" + category + "&page=" + onPage
    const resp = await fetch(url)

    if (!resp.ok) {
        console.log("fetch failed")
        return ["fail"]
    }

    const data = await resp.json()

    if (data.status == "success") {
        var items = data.items
        if (typeof items != "object" /*Will always be object but yk*/ || items == null) {
            console.log("No items.")
            return ["fail"]
        }

        /*for (let e of items) { //Shud always be 8 max
            console.log("Item name: " + e.Name)
            console.log("Item ID: " + e.ID)
        }*/

        return items
    }
}

async function submitRating(uuid, score) {
    const resp = await fetch('/api/submitRating', {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "uuid" : uuid,
            "rating" : score
        })
    })
}

async function openModal(uuid) {
    console.log(uuid);
    await submitRating(uuid, 5)
}

async function switchTab(tabName) {
    selectionContainer.innerHTML = ""
    switch (tabName) {
        case "development":
            console.log("Switching to development tab")
            tabOpen = tabName
            const items = await fetchItems(tabOpen, onPage)
            console.log(items)
            for (let e of items) {
                const div = document.createElement("div")
                const p = document.createElement("p")
                div.id = "frame-card"
                p.id = "center"
                div.setAttribute("uuid", e.ID)
                p.innerText = e.Name
                div.appendChild(p)

                div.addEventListener("click", () => {
                    openModal(div.getAttribute("uuid"))
                })
                selectionContainer.appendChild(div)
            }
            break;
    
        default:
            break;
    }
}

function registerDivs() {
    var divs = selectionContainer.children
    for (let e of divs) {
        if (e.tagName == "DIV") {
            e.addEventListener("click", (ev) => {
                ev.preventDefault() //Just in case
                var href = e.getAttribute("href")
                switchTab(href)
            })
        }
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    //await fetchItems("development", 1)
    registerDivs()
    //switchTab("development")
})