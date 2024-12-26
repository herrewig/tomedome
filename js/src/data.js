export function fetchQuiz(url = "https://api.tomedome.io/api/v1/quiz") {
    return fetch(url).then(response => response.json());
}

export function markupAbilityAnswer(question) {
    return `<strong>${question.abilityName}: </strong>${question.answer}`
}