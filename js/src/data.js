export function fetchQuiz() {
    const url = window.location.hostname === 'localhost'
    ? 'http://localhost:8081/api/v1/quiz'       // Localdev
    : 'https://api.tomedome.io/api/v1/quiz';    // Production API URL
    return fetch(url).then(response => response.json());
}

export function markupAbilityAnswer(question) {
    return `<strong>${question.abilityName}: </strong>${question.answer}`
}