const inputElement = document.querySelector("#youtubeUrl")
const exampleElement = document.querySelector("#exampleUrl")
const exampleUrl = exampleElement.innerHTML

const fillExample = ()=>{
  inputElement.value = exampleUrl
}

exampleElement.addEventListener("click", fillExample)
