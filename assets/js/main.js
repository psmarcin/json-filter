const inputElement = document.querySelector("#youtubeUrl")
const exampleElement = document.querySelector("#exampleUrl")
const generateElement = document.querySelector("#generate")
const feedUrlElement = document.querySelector("#feedUrl")
const feedResultElement = document.querySelector(".feedResult")
const formElement = document.querySelector("#form")

const exampleUrl = exampleElement.innerHTML
const BASE_URL = "https://podcast.psmarcin.me/feed?youtubeUrl="
const fillExample = ()=>{
  inputElement.value = exampleUrl
}

const generateFeedUrl = (e)=>{
  e.preventDefault()
  const ytUrl = inputElement.value
  const feedUrl = BASE_URL + encodeURIComponent(ytUrl)
  feedUrlElement.value = feedUrl
  feedResultElement.style.display = "block"
  window.scroll(0, feedResultElement.offsetTop + 100)
}


exampleElement.addEventListener("click", fillExample)
formElement.addEventListener("submit", generateFeedUrl)
