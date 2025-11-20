//Blank

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
        return
    }

    const data = await resp.json()

    if (data.status == "success") {
        var items = data.items
        for (let e of items) { //Shud always be 8 max
            console.log("Item name: " + e.Name)
            console.log("Item ID: " + e.ID)
        }
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    await fetchItems()
})