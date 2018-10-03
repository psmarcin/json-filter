const inputElement = document.querySelector("#youtubeUrl")
const exampleElement = document.querySelector("#exampleUrl")
const generateElement = document.querySelector("#generate")
const feedUrlElement = document.querySelector("#feedUrl")
const feedResultElement = document.querySelector(".feedResult")
const exampleUrl = exampleElement.innerHTML
const BASE_URL = "http://podcast.psmarcin.me/feed?youtubeUrl="
const fillExample = ()=>{
  inputElement.value = exampleUrl
}

const generateFeedUrl = (e)=>{
  e.preventDefault()
  const ytUrl = inputElement.value
  const feedUrl = BASE_URL + encodeURIComponent(ytUrl)
  feedUrlElement.value = feedUrl
  feedResultElement.style.display = "block"
}


exampleElement.addEventListener("click", fillExample)
generateElement.addEventListener("click", generateFeedUrl)
